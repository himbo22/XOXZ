package di

import (
	"time"

	"github.com/google/wire"
	"github.com/himbo22/xoxz/artist-service/internal/bootstrap"
	"github.com/himbo22/xoxz/artist-service/internal/config"
	"github.com/himbo22/xoxz/artist-service/internal/controller/http/artist"
	"github.com/himbo22/xoxz/artist-service/internal/controller/router"
	"github.com/himbo22/xoxz/artist-service/internal/domain/repository/repo_impl"
	"github.com/himbo22/xoxz/artist-service/internal/logic"
	"github.com/himbo22/xoxz/artist-service/internal/middleware"
	"github.com/himbo22/xoxz/artist-service/internal/service"
	"github.com/himbo22/xoxz/common-service/monitoring/telemetry"
	xoxzEcho "github.com/himbo22/xoxz/common-service/xoxz/echo"
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
	logic.NewArtistLogic,
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
	e.Use(telemetry.EchoMiddleware())
	e.Use(telemetry.EchoMetricsMiddleware())
	e.IPExtractor = echo.ExtractIPFromXFFHeader()

	e.HTTPErrorHandler = func(c *echo.Context, err error) {
		xoxzEcho.ErrorHandler(c, logger, err)
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
