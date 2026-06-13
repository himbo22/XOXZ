package di

import (
	"context"
	"log"
	"time"

	"github.com/google/wire"
	"github.com/himbo22/xoxz/common-service/monitoring/telemetry"
	xoxzEcho "github.com/himbo22/xoxz/common-service/xoxz/echo"
	xoxzLogger "github.com/himbo22/xoxz/common-service/xoxz/logger"
	livekit "github.com/himbo22/xoxz/livestream-service/internal/adapter/livekit"
	"github.com/himbo22/xoxz/livestream-service/internal/bootstrap"
	"github.com/himbo22/xoxz/livestream-service/internal/config"
	"github.com/himbo22/xoxz/livestream-service/internal/controller/http/livestream"
	"github.com/himbo22/xoxz/livestream-service/internal/controller/http/webhook"
	"github.com/himbo22/xoxz/livestream-service/internal/controller/router"
	"github.com/himbo22/xoxz/livestream-service/internal/domain/repository/repo_impl"
	"github.com/himbo22/xoxz/livestream-service/internal/logic"
	"github.com/himbo22/xoxz/livestream-service/internal/service"
	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var InfrastructureSet = wire.NewSet(
	provideAppLogger,
	provideMongoDatabase,
	provideRedisClient,
	provideOtelCollector,
	provideLiveKitSDK,
	provideEchoApp,
)

var ControllerSet = wire.NewSet(
	livestream.NewLiveStreamController,
	webhook.NewWebhookController,
	provideControllers,
)

var RepositorySet = wire.NewSet(
	repo_impl.NewLivestreamRepository,
)

var LiveStreamSet = wire.NewSet(
	logic.NewLiveStreamLogic,
	service.NewLiveStreamService,
)

// provide xoxz xoxz (interface)
func provideAppLogger(cfg *config.Config) (xoxzLogger.XoxzLogger, func()) {
	logger, cleanup := bootstrap.InitLogger(cfg.Logger)
	return xoxzLogger.NewxoxzLogger(logger.Logger), cleanup
}

func provideControllers(
	liveController livestream.LiveStreamController,
	webhookController webhook.WebhookController,
) router.Controllers {
	return router.Controllers{
		LiveStreamController: liveController,
		WebhookController:    webhookController,
	}
}

type OtelTracerToken struct{}

func provideOtelCollector(cfg *config.Config) (*OtelTracerToken, func(), error) {
	shutdown, err := bootstrap.InitOtelCollector(cfg.Otel)
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		// Allocate a dedicated Context with Timeout for the OTel shutdown process
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Call the actual shutdown function
		if err := shutdown(ctx); err != nil {
			// Wire cannot return error outside, so we must log here
			// (You can use a.logger if logger was initialized before,
			// or use Go's standard log to ensure it always prints during shutdown)
			log.Printf("[Graceful Shutdown] Otel Collector failed to clean up: %v", err)
		} else {
			log.Println("[Graceful Shutdown] Otel Collector closed gracefully")
		}
	}
	return &OtelTracerToken{}, cleanup, nil
}

func provideLiveKitSDK(cfg *config.Config, logger xoxzLogger.XoxzLogger) (*livekit.LiveKitSDK, error) {
	client, err := bootstrap.NewLiveKitSDK(cfg.LiveKit, logger)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func provideEchoApp(
	controller router.Controllers,
	logger xoxzLogger.XoxzLogger,
) *echo.Echo {
	e := echo.New()

	// Recover + request id + logging should be baseline
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.RequestID())
	e.Use(echoMiddleware.BodyLimit(2_097_152))
	//e.Use(echoMiddleware.CORS(""))
	e.Use(telemetry.EchoMiddleware())
	e.Use(telemetry.EchoMetricsMiddleware())
	e.IPExtractor = echo.ExtractIPFromXFFHeader()

	// Error handler
	e.HTTPErrorHandler = func(c *echo.Context, err error) {
		xoxzEcho.ErrorHandler(c, logger, err)
	}

	router.SetupRouters(e, controller)
	return e
}

func provideMongoDatabase(cfg *config.Config) (*mongo.Database, func(), error) {
	client, cleanup, err := bootstrap.InitMongoDB(cfg.Mongo)
	if err != nil {
		return nil, nil, err
	}

	return client, cleanup, nil
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

func NewApp(cfg *config.Config, echoApp *echo.Echo, _ *OtelTracerToken) *App {
	return &App{
		Config:  cfg,
		EchoApp: echoApp,
	}
}
