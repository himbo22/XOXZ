package di

import (
	"github.com/himbo22/xoxz/media-service/internal/config"
	"github.com/himbo22/xoxz/media-service/internal/grpc"
	"github.com/labstack/echo/v5"
)

type App struct {
	Config     *config.Config
	EchoApp    *echo.Echo
	GrpcServer *grpc.MediaServer
	Cleanup    func()
}
