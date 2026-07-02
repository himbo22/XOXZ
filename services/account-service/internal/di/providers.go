package di

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/wire"
	mediaGrpc "github.com/himbo22/xoxz/account-service/internal/adapter/grpc"
	"github.com/himbo22/xoxz/account-service/internal/adapter/oauth"
	"github.com/himbo22/xoxz/account-service/internal/bootstrap"
	"github.com/himbo22/xoxz/account-service/internal/config"
	"github.com/himbo22/xoxz/account-service/internal/controller/http/admin"
	"github.com/himbo22/xoxz/account-service/internal/controller/http/auth"
	profilecontroller "github.com/himbo22/xoxz/account-service/internal/controller/http/profile"
	"github.com/himbo22/xoxz/account-service/internal/controller/router"
	"github.com/himbo22/xoxz/account-service/internal/domain/repository"
	"github.com/himbo22/xoxz/account-service/internal/domain/repository/repository_impl"
	"github.com/himbo22/xoxz/account-service/internal/middleware"
	"github.com/himbo22/xoxz/account-service/internal/model"
	"github.com/himbo22/xoxz/account-service/internal/service"
	"github.com/himbo22/xoxz/account-service/internal/util"
	"github.com/himbo22/xoxz/common-service/monitoring/telemetry"
	xoxzLogger "github.com/himbo22/xoxz/common-service/xoxz/logger"
	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var InfrastructureSet = wire.NewSet(
	provideAppLogger,
	providePostgreSQL,
	provideRedisClient,
	provideOtelCollector,
)

var PartySet = wire.NewSet(
	provideEchoApp,
	provideGrpcMediaClient,
	provideGoogleAuth,
	wire.Bind(new(service.GoogleOAuthProvider), new(*oauth.GoogleOAuthAdapter)),
)

var ControllerSet = wire.NewSet(
	auth.NewAuthController,
	profilecontroller.NewProfileController,
	admin.NewAdminController,
	provideControllers,
)

var RepositorySet = wire.NewSet(
	repository_impl.NewRedisRepository,
	repository_impl.NewTxManager,
	repository_impl.NewUserRepository,
	repository_impl.NewIdentityRepository,
	repository_impl.NewRoleRepository,
	repository_impl.NewUserRoleRepository,
	repository_impl.NewPermissionRepository,
	repository_impl.NewRolePermissionRepository,
)

var AuthSet = wire.NewSet(
	service.NewAuthService,
)

var ProfileSet = wire.NewSet(
	service.NewProfileService,
)

var AdminSet = wire.NewSet(
	service.NewAdminService,
)

var MiddlewareSet = wire.NewSet(
	provideAuthMiddleware,
)

var CacheSet = wire.NewSet(
	provideAccessControl,
)

// provide xoxz logger (interface)
func provideAppLogger(cfg *config.Config) (xoxzLogger.XoxzLogger, func()) {
	logger, cleanup := bootstrap.InitLogger(cfg.Logger)
	return xoxzLogger.NewxoxzLogger(logger.Logger), cleanup
}

func provideControllers(
	auth auth.AuthController,
	profile profilecontroller.ProfileController,
	admin admin.AdminController,
) router.Controllers {
	return router.Controllers{
		AuthController:    auth,
		ProfileController: profile,
		AdminController:   admin,
	}
}

func provideAuthMiddleware(
	redisRepo repository.RedisRepository,
	permissions *bootstrap.AccessControl,
) middleware.AuthMiddleware {
	return *middleware.NewAuthMiddleware(redisRepo, permissions)
}

type OtelTracerToken struct{}

func provideOtelCollector(cfg *config.Config) (*OtelTracerToken, func(), error) {
	shutdown, err := bootstrap.InitOtelCollector(cfg.Otel)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		// Use a dedicated timeout context for OTel shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Run the actual shutdown function
		if err := shutdown(ctx); err != nil {
			// Wire cannot return this error, so log it here
			log.Printf("[Graceful Shutdown] Otel Collector failed to clean up: %v", err)
		} else {
			log.Println("[Graceful Shutdown] Otel Collector closed gracefully")
		}
	}

	return &OtelTracerToken{}, cleanup, nil
}

func provideEchoApp(
	controller router.Controllers,
	logger xoxzLogger.XoxzLogger,
	authMiddleware middleware.AuthMiddleware,
) *echo.Echo {
	e := echo.New()

	// Recover, request ID, and logging should be the baseline.
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.RequestID())
	e.Use(echoMiddleware.CORSWithConfig(
		echoMiddleware.CORSConfig{
			AllowOrigins: []string{"https://377a-2001-ee0-4b6e-e810-f896-768f-224d-9675.ngrok-free.app", "http://localhost:3000"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
			AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		},
	))

	e.Use(echoMiddleware.BodyLimit(2_097_152))
	e.Use(middleware.ContextMiddleware())
	e.Use(telemetry.EchoTracingMiddleware())
	e.Use(telemetry.EchoMetricsMiddleware())
	e.IPExtractor = echo.ExtractIPFromXFFHeader()

	// Error handler
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

	router.SetupRouters(e, controller, &authMiddleware)
	return e
}

func provideGoogleAuth(cfg *config.Config) *oauth.GoogleOAuthAdapter {
	authConfig := oauth.GoogleOAuthConfig{
		ClientID:     cfg.OAuth.Google.ClientID,
		ClientSecret: cfg.OAuth.Google.ClientSecret,
		RedirectURL:  cfg.OAuth.Google.RedirectURL,
	}
	return oauth.NewGoogleOAuthAdapter(authConfig)
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

func provideRedisClient(cfg *config.Config) (*redis.Client, func(), error) {
	redisConfig := bootstrap.RedisConfig{
		Address:      cfg.Redis.Default.Address,
		Password:     cfg.Redis.Default.Password,
		DB:           cfg.Redis.Default.DB,
		DialTimeout:  cfg.Redis.Default.DialTimeout,
		ReadTimeout:  cfg.Redis.Default.ReadTimeout,
		WriteTimeout: cfg.Redis.Default.WriteTimeout,
		MaxActive:    cfg.Redis.Default.MaxActive,
	}
	client, err := bootstrap.InitRedis(redisConfig)
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() { bootstrap.CloseRedis(client) }
	return client, cleanup, nil
}

func provideGrpcMediaClient() (mediaGrpc.MediaClient, func(), error) {
	// Address could be pulled from config. We'll use a hardcoded value or os.Getenv for simplicity
	// Alternatively, assuming es-svc is reachable at es-svc:50051 via docker-compose DNS
	address := "localhost:50051" // Use "media-service:50051" when running in Docker.
	client, err := mediaGrpc.NewMediaClient(address)
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() { _ = client.Close() }
	return client, cleanup, nil
}

func provideAccessControl(
	roleRepo repository.RolePermissionRepository,
	permRepo repository.PermissionRepository,
) (*bootstrap.AccessControl, func(), error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	roles, err := roleRepo.GetAll(ctx)
	if err != nil {
		return nil, nil, err
	}
	permissions, err := permRepo.GetAll(ctx)
	if err != nil {
		return nil, nil, err
	}

	ac, cleanup := bootstrap.InitAccessControl(roles, permissions)

	return ac, cleanup, nil
}

func NewApp(cfg *config.Config, echoApp *echo.Echo, _ *OtelTracerToken) *App {
	return &App{
		Config:  cfg,
		EchoApp: echoApp,
	}
}
