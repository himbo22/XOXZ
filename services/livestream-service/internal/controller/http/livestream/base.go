package livestream

import (
	xoxz "github.com/himbo22/xoxz/common-service/xoxz/logger"
	"github.com/himbo22/xoxz/livestream-service/internal/model"
	"github.com/himbo22/xoxz/livestream-service/internal/service"
	"github.com/himbo22/xoxz/livestream-service/internal/util"
	"github.com/labstack/echo/v5"
)

type liveStreamController struct {
	livestreamService service.LiveStreamService
	logger            xoxz.XoxzLogger
}

func (l *liveStreamController) Create(ctx *echo.Context) error {
	req := model.CreateStreamRequest{}
	if err := ctx.Bind(&req); err != nil {
		return err
	}

	stream, err := l.livestreamService.Create(ctx.Request().Context(), req)
	if err != nil {
		return err
	}

	return util.SuccessResponse(ctx, 200, 123, "hoang ne", stream)
}

func (l *liveStreamController) Stop(ctx *echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func NewLiveStreamController(
	logger xoxz.XoxzLogger,
	livestreamService service.LiveStreamService,
) LiveStreamController {
	return &liveStreamController{
		logger:            logger,
		livestreamService: livestreamService,
	}
}
