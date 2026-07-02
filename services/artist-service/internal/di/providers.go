package di

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/wire"
	"github.com/himbo22/xoxz/artist-service/internal/bootstrap"
	"github.com/himbo22/xoxz/artist-service/internal/config"
	"github.com/himbo22/xoxz/artist-service/internal/controller/http/artist"
	"github.com/himbo22/xoxz/artist-service/internal/controller/router"
	"github.com/himbo22/xoxz/artist-service/internal/domain/repository/repo_impl"
	"github.com/himbo22/xoxz/artist-service/internal/middleware"
	"github.com/himbo22/xoxz/artist-service/internal/model"
	"github.com/himbo22/xoxz/artist-service/internal/service"
	"github.com/himbo22/xoxz/artist-service/internal/util"
	"github.com/himbo22/xoxz/common-service/monitoring/telemetry"
	xoxzLogger "github.com/himbo22/xoxz/common-service/xoxz/logger"
	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
	"gorm.io/gorm"
)

var InfrastructureSet = wire.NewSet(
	provideAppLogger,
	providePostgreSQL,
	provideOtelCollector,
	provideEchoApp,
)

var ControllerSet = wire.NewSet(
	artist.NewArtistController,
	provideControllers,
)

var RepositorySet = wire.NewSet(
	repo_impl.NewArtistRepository,
)

var ArtistSet = wire.NewSet(
	service.NewArtistService,
)

func provideAppLogger(cfg *config.Config) (xoxzLogger.XoxzLogger, func()) {
	logger, cleanup := bootstrap.InitLogger(cfg.Logger)
	return xoxzLogger.NewxoxzLogger(logger.Logger), cleanup
}

func provideControllers(artistController artist.ArtistController) router.Controllers {
	return router.Controllers{
		ArtistController: artistController,
	}
}

type OtelTracerToken struct{}

func provideOtelCollector(cfg *config.Config) (*OtelTracerToken, func(), error) {
	shutdown, err := bootstrap.InitOtelCollector(cfg.Otel)
	if err != nil {
		return nil, nil, err
	}
	return &OtelTracerToken{}, shutdown, nil
}

func provideEchoApp(
	controller router.Controllers,
	logger xoxzLogger.XoxzLogger,
) *echo.Echo {
	e := echo.New()

	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.RequestID())
	e.Use(echoMiddleware.BodyLimit(2_097_152))
	e.Use(middleware.ContextMiddleware())
	e.Use(telemetry.EchoTracingMiddleware())
	e.Use(telemetry.EchoMetricsMiddleware())
	e.IPExtractor = echo.ExtractIPFromXFFHeader()

	e.HTTPErrorHandler = func(c *echo.Context, err error) {
		httpCode := http.StatusInternalServerError
		res := model.ResponseModel{
			StatusCode: 9999, // General system error code
			Message:    "System error, please try again later",
		}

		// 1. Check if this is an AppError that we intentionally threw
		if appErr, ok := err.(*util.AppError); ok {
			httpCode = appErr.HTTPCode
			res.StatusCode = appErr.CustomCode
			res.Message = appErr.Message
			res.Error = appErr.Detail
		} else if errors.Is(err, echo.ErrNotFound) {
			// 2. Check for echo.ErrNotFound (default 404)
			httpCode = 404
			res.Message = "Not Found"
		} else if he, ok := err.(*echo.HTTPError); ok {
			// 2. Check if Echo internally threw this error (e.g., 404 wrong URL, 413 body too large)
			httpCode = he.Code
			res.Message = fmt.Sprintf("%v", he.Message)
		} else {
			// 3. If we land here, there's a code bug, panic, or raw DB error.
			// THIS IS WHERE YOU USE ZAP LOGGER TO WRITE LOG FILES AND FIND THE REAL ERROR
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				logger.Warn("Client disconnected or timeout", xoxzLogger.Error(err))
				httpCode = 499 // Nginx standard code for Client Closed Request
				res.Message = "Client disconnected"
			} else {
				// --- THIS IS A REAL BUG ---
				logger.Error("Unhandled system error",
					xoxzLogger.Error(err),
					xoxzLogger.String("path", c.Request().URL.Path), // Log the URL being called for easier debugging
				)
			}
		}

		// Send JSON response to Frontend
		if resp, uErr := echo.UnwrapResponse(c.Response()); uErr == nil {
			// Fixed logic: Only write JSON when HTTP headers have NOT been sent yet (!Committed)
			if !resp.Committed {
				err := c.JSON(httpCode, res)
				if err != nil {
					logger.Errorf("Error writing error response: %v", err)
				}
			} else {
				// Response was committed mid-flight (e.g., streaming file then DB crashed)
				logger.Errorf("Error %v occurred but HTTP was already committed. Skipping JSON write.", err)
			}
		}
	}

	router.SetupRouters(e, controller)
	return e
}

func providePostgreSQL(cfg *config.Config) (*gorm.DB, func(), error) {
	dbConfig := bootstrap.DatabaseConfig{
		Host:            cfg.Database.Host,
		Port:            cfg.Database.Port,
		User:            cfg.Database.User,
		Password:        cfg.Database.Password,
		DBName:          cfg.Database.DBName,
		SSLMode:         cfg.Database.SSLMode,
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: time.Duration(cfg.Database.ConnMaxLifetime) * time.Hour,
		ConnMaxIdleTime: time.Duration(cfg.Database.ConnMaxLifetime) * time.Hour,
		Timezone:        cfg.Database.Timezone,
	}

	db, err := bootstrap.InitPostgreSQL(dbConfig)
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() { _ = bootstrap.CloseDatabase(db) }
	return db, cleanup, nil
}

func NewApp(cfg *config.Config, echoApp *echo.Echo, _ *OtelTracerToken) *App {
	return &App{
		Config:  cfg,
		EchoApp: echoApp,
	}
}
