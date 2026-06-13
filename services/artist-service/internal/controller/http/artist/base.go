package artist

import (
	"net/http"

	"github.com/himbo22/xoxz/artist-service/internal/model"
	"github.com/himbo22/xoxz/artist-service/internal/service"
	"github.com/himbo22/xoxz/artist-service/internal/util"
	xoxz "github.com/himbo22/xoxz/common-service/xoxz/logger"
	"github.com/labstack/echo/v5"
)

type artistController struct {
	artistService service.ArtistService
	logger        xoxz.XoxzLogger
}

func (a *artistController) Create(ctx *echo.Context) error {
	req := model.CreateArtistRequest{}
	if err := ctx.Bind(&req); err != nil {
		return err
	}

	artist, err := a.artistService.Create(ctx.Request().Context(), req)
	if err != nil {
		return err
	}

	return util.SuccessResponse(ctx, http.StatusCreated, 0, "artist created", artist)
}

func (a *artistController) GetByID(ctx *echo.Context) error {
	artist, err := a.artistService.GetByID(ctx.Request().Context(), ctx.Param("id"))
	if err != nil {
		return err
	}

	return util.SuccessResponse(ctx, http.StatusOK, 0, "success", artist)
}

func (a *artistController) List(ctx *echo.Context) error {
	req := model.ListArtistRequest{
		Page:   util.ParseIntQuery(ctx, "page", 1),
		Limit:  util.ParseIntQuery(ctx, "limit", 20),
		Search: ctx.QueryParam("search"),
	}

	artists, err := a.artistService.List(ctx.Request().Context(), req)
	if err != nil {
		return err
	}

	return util.SuccessResponse(ctx, http.StatusOK, 0, "success", artists)
}

func (a *artistController) Update(ctx *echo.Context) error {
	req := model.UpdateArtistRequest{}
	if err := ctx.Bind(&req); err != nil {
		return err
	}

	artist, err := a.artistService.Update(ctx.Request().Context(), ctx.Param("id"), req)
	if err != nil {
		return err
	}

	return util.SuccessResponse(ctx, http.StatusOK, 0, "artist updated", artist)
}

func (a *artistController) Delete(ctx *echo.Context) error {
	if err := a.artistService.Delete(ctx.Request().Context(), ctx.Param("id")); err != nil {
		return err
	}

	return util.SuccessResponse(ctx, http.StatusNoContent, 0, "artist deleted", nil)
}

func NewArtistController(
	logger xoxz.XoxzLogger,
	artistService service.ArtistService,
) ArtistController {
	return &artistController{
		logger:        logger,
		artistService: artistService,
	}
}
