package media

import (
	"net/http"

	_const "github.com/himbo22/xoxz/media-service/internal/const"
	"github.com/himbo22/xoxz/media-service/internal/model"
	"github.com/himbo22/xoxz/media-service/internal/service"
	"github.com/himbo22/xoxz/media-service/internal/util"
	"github.com/labstack/echo/v5"
)

type mediaController struct {
	mediaService service.MediaService
}

// ConfirmUpload implements [MediaController].
func (m *mediaController) ConfirmUpload(ctx *echo.Context) error {
	panic("unimplemented")
}

// DeleteMedia implements [MediaController].
func (m *mediaController) DeleteMedia(ctx *echo.Context) error {
	panic("unimplemented")
}

// GeneratePresignedURL implements [MediaController].
func (m *mediaController) GeneratePresignedURL(ctx *echo.Context) error {
	req := model.GenerateURLRequest{}
	if err := ctx.Bind(&req); err != nil {
		return util.NewError(http.StatusBadRequest, _const.CodeInvalidRequest.Code, "invalid request body")
	}

	// validation
	if req.ContentType == "" || req.FileName == "" {
		return util.NewError(http.StatusBadRequest, _const.CodeInvalidRequest.Code, "missing fields")
	}

	response, err := m.mediaService.GeneratePresignedURL(ctx.Request().Context(), req)
	if err != nil {
		return err
	}

	return util.SuccessResponse(ctx, http.StatusOK, _const.CodeSuccess.Code, "Ok presigned url", response)
}

func NewMediaController(mediaService service.MediaService) MediaController {
	return &mediaController{
		mediaService: mediaService,
	}
}
