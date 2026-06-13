package profile

import (
	"net/http"
	"strings"

	_const "github.com/himbo22/xoxz/account-service/internal/const"
	"github.com/himbo22/xoxz/account-service/internal/middleware"
	"github.com/himbo22/xoxz/account-service/internal/model"
	"github.com/himbo22/xoxz/account-service/internal/service"
	"github.com/himbo22/xoxz/account-service/internal/util"
	"github.com/labstack/echo/v5"
)

type profileController struct {
	profileService service.ProfileService
}

// GetProfile godoc
// @Summary      Get current profile
// @Description  Returns the authenticated user's profile.
// @Tags         Profile
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization  header    string  true  "Bearer access token"
// @Param        Device-ID      header    string  true  "Unique device identifier"
// @Param        Device-Type    header    string  true  "Device type: 1=web, 2=mobile"
// @Success      200            {object}  model.ProfileResponseDoc
// @Failure      400            {object}  model.MessageResponseDoc
// @Failure      401            {object}  model.MessageResponseDoc
// @Failure      403            {object}  model.MessageResponseDoc
// @Failure      404            {object}  model.MessageResponseDoc
// @Failure      500            {object}  model.MessageResponseDoc
// @Router       /api/v1/public/profile/me [get]
// GetProfile implements [ProfileController].
func (p *profileController) GetProfile(ctx *echo.Context) error {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return util.NewErrorByCode(_const.CodeInvalidRequest)
	}

	payload, err := p.profileService.GetProfile(ctx.Request().Context(), userID)
	if err != nil {
		return err
	}

	return util.SuccessResponse(ctx, http.StatusOK, _const.CodeSuccess.Code, "get profile successful", payload)
}

// GetPublicProfile godoc
// @Summary      Get public profile
// @Description  Returns a public profile by username.
// @Tags         Profile
// @Produce      json
// @Param        username  path      string  true  "Username"
// @Success      200       {object}  model.PublicProfileResponseDoc
// @Failure      400       {object}  model.MessageResponseDoc
// @Failure      404       {object}  model.MessageResponseDoc
// @Failure      500       {object}  model.MessageResponseDoc
// @Router       /api/v1/public/profile/{username} [get]
// GetPublicProfile implements [ProfileController].
func (p *profileController) GetPublicProfile(c *echo.Context) error {
	username := strings.TrimSpace(c.Param("username"))
	if username == "" {
		return util.NewError(http.StatusBadRequest, _const.CodeInvalidRequest.Code, "invalid username")
	}

	payload, svcErr := p.profileService.GetPublicProfile(c.Request().Context(), username)
	if svcErr != nil {
		return svcErr
	}

	return util.SuccessResponse(c, http.StatusOK, _const.CodeSuccess.Code, "get public profile successful", payload)
}

// UpdateAvatar godoc
// @Summary      Update avatar
// @Description  Updates the authenticated user's avatar.
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization  header    string                     true  "Bearer access token"
// @Param        Device-ID      header    string                     true  "Unique device identifier"
// @Param        Device-Type    header    string                     true  "Device type: 1=web, 2=mobile"
// @Param        payload        body      model.UpdateAvatarRequest  true  "Avatar update payload"
// @Success      200            {object}  model.UpdateAvatarResponseDoc
// @Failure      400            {object}  model.MessageResponseDoc
// @Failure      401            {object}  model.MessageResponseDoc
// @Failure      403            {object}  model.MessageResponseDoc
// @Failure      404            {object}  model.MessageResponseDoc
// @Failure      500            {object}  model.MessageResponseDoc
// @Router       /api/v1/public/profile/me/avatar [put]
// UpdateAvatar implements [ProfileController].
func (p *profileController) UpdateAvatar(ctx *echo.Context) error {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return util.NewErrorByCode(_const.CodeInvalidRequest)
	}

	req := model.UpdateAvatarRequest{}
	if err := ctx.Bind(&req); err != nil {
		return util.NewError(http.StatusBadRequest, _const.CodeInvalidRequest.Code, "invalid request body")
	}

	req.TmpPath = strings.TrimSpace(req.TmpPath)
	req.PerPath = strings.TrimSpace(req.PerPath)
	if req.TmpPath == "" || req.PerPath == "" {
		return util.NewError(http.StatusBadRequest, _const.CodeInvalidRequest.Code, "tmp_path and per_path are required")
	}
	req.UserID = userID

	payload, err := p.profileService.UpdateAvatar(ctx.Request().Context(), req)
	if err != nil {
		return err
	}

	return util.SuccessResponse(ctx, http.StatusOK, _const.CodeSuccess.Code, "update avatar successful", payload)
}

// UpdateProfile godoc
// @Summary      Update current profile
// @Description  Updates editable fields on the authenticated user's profile.
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization  header    string                      true  "Bearer access token"
// @Param        Device-ID      header    string                      true  "Unique device identifier"
// @Param        Device-Type    header    string                      true  "Device type: 1=web, 2=mobile"
// @Param        payload        body      model.UpdateProfileRequest  true  "Profile update payload"
// @Success      200            {object}  model.ProfileResponseDoc
// @Failure      400            {object}  model.MessageResponseDoc
// @Failure      401            {object}  model.MessageResponseDoc
// @Failure      403            {object}  model.MessageResponseDoc
// @Failure      404            {object}  model.MessageResponseDoc
// @Failure      500            {object}  model.MessageResponseDoc
// @Router       /api/v1/public/profile/me [put]
// UpdateProfile implements [ProfileController].
func (p *profileController) UpdateProfile(ctx *echo.Context) error {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return util.NewErrorByCode(_const.CodeInvalidRequest)
	}

	req := model.UpdateProfileRequest{}
	if err := ctx.Bind(&req); err != nil {
		return util.NewError(http.StatusBadRequest, _const.CodeInvalidRequest.Code, "invalid request body")
	}
	if req.Username != nil {
		username := strings.TrimSpace(*req.Username)
		if username == "" {
			return util.NewError(http.StatusBadRequest, _const.CodeInvalidRequest.Code, "username cannot be empty")
		}
		req.Username = &username
	}

	payload, err := p.profileService.UpdateProfile(ctx.Request().Context(), userID, req)
	if err != nil {
		return err
	}

	return util.SuccessResponse(ctx, http.StatusOK, _const.CodeSuccess.Code, "update profile successful", payload)
}

func NewProfileController(profileService service.ProfileService) ProfileController {
	return &profileController{
		profileService: profileService,
	}
}
