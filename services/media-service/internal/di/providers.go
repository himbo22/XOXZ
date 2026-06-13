package di

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/google/wire"
	"github.com/himbo22/xoxz/media-service/internal/adapter/storage"
	"github.com/himbo22/xoxz/media-service/internal/config"
	"github.com/himbo22/xoxz/media-service/internal/controller/http/media"
	"github.com/himbo22/xoxz/media-service/internal/controller/router"
	"github.com/himbo22/xoxz/media-service/internal/grpc"
	"github.com/himbo22/xoxz/media-service/internal/model"
	"github.com/himbo22/xoxz/media-service/internal/service"
	"github.com/himbo22/xoxz/media-service/internal/util"
	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
)

var InfrastructureSet = wire.NewSet(
	provideEchoApp,
	provideMinioClient,
	provideGrpcServer,
)

var ControllerSet = wire.NewSet(
	media.NewMediaController,
	provideControllers,
)

var MediaSet = wire.NewSet(
	service.NewMediaService,
)

func provideControllers(
	media media.MediaController,
) router.Controllers {
	return router.Controllers{
		MediaController: media,
	}
}

func provideMinioClient(cfg *config.Config) (*storage.MinioClient, error) {
	config := storage.Config{
		Endpoint:        cfg.Minio.Endpoint,
		AccessKeyID:     cfg.Minio.AccessKeyID,
		SecretAccessKey: cfg.Minio.SecretAccessKey,
		UseSSL:          cfg.Minio.UseSSL,
		BucketName:      cfg.Minio.BucketName,
	}

	minioClient, err := storage.NewClient(config)
	if err != nil {
		return nil, err
	}

	return minioClient, nil
}

func provideEchoApp(controller router.Controllers) *echo.Echo {
	e := echo.New()
	// Recover + request id + logging should be baseline
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.RequestID())
	e.Use(echoMiddleware.CORS("https://1266-2001-ee1-db02-fbb0-3456-37bf-a199-6e1c.ngrok-free.app"))
	e.Use(echoMiddleware.BodyLimit(2_097_152))
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
		} else if he, ok := err.(*echo.HTTPError); ok {
			// 2. Check if Echo internally threw this error (e.g., 404 wrong URL, 413 body too large)
			httpCode = he.Code
			res.Message = fmt.Sprintf("%v", he.Message)
		} else {
			// 3. If we land here, there's a code bug, panic, or raw DB error.
			// THIS IS WHERE YOU USE ZAP LOGGER TO WRITE LOG FILES AND FIND THE REAL ERROR
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				httpCode = 499 // Nginx standard code for Client Closed Request
				res.Message = "Client disconnected"
			} else {
				// --- THIS IS A REAL BUG ---
				log.Printf("Unhandled system error: %v", err)
			}
		}

		// Send JSON response to Frontend
		if resp, uErr := echo.UnwrapResponse(c.Response()); uErr == nil {
			// Fixed logic: Only write JSON when HTTP headers have NOT been sent yet (!Committed)
			if !resp.Committed {
				err := c.JSON(httpCode, res)
				if err != nil {
					log.Printf("Error writing error response: %v", err)
				}
			} else {
				// Response was committed mid-flight (e.g., streaming file then DB crashed)
				log.Printf("Error %v occurred but HTTP was already committed. Skipping JSON write.", err)
			}
		}
	}

	router.SetupRouters(e, controller)
	return e
}

func provideGrpcServer(storage *storage.MinioClient) *grpc.MediaServer {
	server := grpc.NewMediaServer(storage)
	return server
}

func NewApp(cfg *config.Config, echoApp *echo.Echo, grpcServer *grpc.MediaServer) *App {
	return &App{
		Config:     cfg,
		EchoApp:    echoApp,
		GrpcServer: grpcServer,
	}
}
