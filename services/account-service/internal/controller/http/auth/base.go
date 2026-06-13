package auth

import (
	"net/http"

	"github.com/himbo22/xoxz/account-service/internal/config"
	_const "github.com/himbo22/xoxz/account-service/internal/const"
	"github.com/himbo22/xoxz/account-service/internal/middleware"
	"github.com/himbo22/xoxz/account-service/internal/model"
	"github.com/himbo22/xoxz/account-service/internal/service"
	"github.com/himbo22/xoxz/account-service/internal/util"
	"github.com/himbo22/xoxz/common-service/monitoring/telemetry"
	xoxz "github.com/himbo22/xoxz/common-service/xoxz/logger"
	"github.com/labstack/echo/v5"
)

type authController struct {
	authService service.AuthService
	logger      xoxz.XoxzLogger
}

// Google godoc
// @Summary      Sign in with Google
// @Description  Authenticates a user with a Google ID token. Web clients receive the refresh token via HttpOnly cookie; mobile clients receive both tokens in the response body and should persist them manually on the client side.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        Device-ID    header    string                    true   "Unique device identifier"
// @Param        Device-Type  header    string                    true   "Device type: 1=web, 2=mobile"
// @Param        payload        body      model.GoogleLoginRequest  true   "Google login payload"
// @Success      200            {object}  model.AuthTokenResponseDoc
// @Failure      400            {object}  model.MessageResponseDoc
// @Failure      401            {object}  model.MessageResponseDoc
// @Failure      500            {object}  model.MessageResponseDoc
// @Router       /api/v1/public/auth/google [post]
func (a *authController) Google(c *echo.Context) error {
	deviceID, ok := middleware.GetDeviceIDFromHeader(c)
	if !ok {
		return util.NewError(http.StatusBadRequest, _const.CodeInvalidRequest.Code, "missing header: device-id")
	}

	deviceType, ok := middleware.GetDeviceTypeFromHeader(c)
	if !ok {
		return util.NewError(http.StatusBadRequest, _const.CodeInvalidRequest.Code, "invalid device type header")
	}

	req := model.GoogleLoginRequest{}
	if err := c.Bind(&req); err != nil {
		return util.NewError(http.StatusBadRequest, _const.CodeInvalidRequest.Code, "invalid request body")
	}
	if req.Token == "" {
		return util.NewError(http.StatusBadRequest, _const.CodeInvalidRequest.Code, "invalid request body")
	}

	req.DeviceID = deviceID

	// start span for monitoring
	ctx := c.Request().Context()

	ctx, end := telemetry.StartSpan(ctx, "test1", "controller")

	payload, err := a.authService.AuthenticateWithGoogle(ctx, req)
	if err != nil {
		util.TraceFields(ctx, a.logger).Error("Not Ok")
		end(err)
		return err
	}

	end(nil)
	return a.respondByDeviceType(c, payload, deviceType)
	// 2
	//var err error
	//defer func() { end(err) }()
	//
	//payload, err := a.authService.AuthenticateWithGoogle(ctx, req)
	//if err != nil {
	//	return err
	//}
	//
	//return a.respondByDeviceType(c, payload, deviceType)
}

// Refresh godoc
// @Summary      Refresh access token
// @Description  Refreshes the current session token. Mobile clients send refresh_token in the request body. Web clients send device headers and use the HttpOnly refresh_token cookie.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        payload        body      model.RefreshTokenRequest  false  "Mobile refresh payload. Web clients use cookie-based refresh_token instead."
// @Param        Device-ID    header    string  true   "Unique device identifier"
// @Param        Device-Type  header    string  true   "Device type: 1=web, 2=mobile"
// @Success      200            {object}  model.AuthTokenResponseDoc
// @Failure      400            {object}  model.MessageResponseDoc
// @Failure      401            {object}  model.MessageResponseDoc
// @Failure      403            {object}  model.MessageResponseDoc
// @Failure      500            {object}  model.MessageResponseDoc
// @Router       /api/v1/public/auth/session/refresh [post]
func (a *authController) Refresh(ctx *echo.Context) error {
	// device ID
	deviceID, ok := middleware.GetDeviceIDFromHeader(ctx)
	if !ok {
		return util.NewError(http.StatusBadRequest, _const.CodeInvalidRequest.Code, "missing header: device-id")
	}

	deviceType, ok := middleware.GetDeviceTypeFromHeader(ctx)
	if !ok {
		return util.NewError(http.StatusBadRequest, _const.CodeInvalidRequest.Code, "invalid device type header")
	}

	switch deviceType {
	case _const.MobileType:
		req := model.RefreshTokenRequest{
			DeviceID: deviceID,
		}

		if err := ctx.Bind(&req); err != nil {
			return util.NewError(http.StatusBadRequest, _const.CodeInvalidRequest.Code, "invalid request body")
		}

		if req.RefreshToken == "" {
			return util.NewError(http.StatusBadRequest, _const.CodeInvalidRequest.Code, "missing refresh token")
		}

		payload, err := a.authService.RefreshToken(ctx.Request().Context(), req)
		if err != nil {
			return err
		}

		return util.SuccessResponse(ctx, http.StatusOK, _const.CodeSuccess.Code, "refresh ok", payload)
	case _const.WebType:
		// Web clients store the refresh token in an HttpOnly cookie here.
		// refresh token from cookies
		refreshToken, ok := middleware.GetRefreshToken(ctx)
		if !ok {
			return util.NewError(http.StatusBadRequest, _const.CodeInvalidRequest.Code, "missing refresh token in cookies")
		}

		req := model.RefreshTokenRequest{
			RefreshToken: refreshToken,
			DeviceID:     deviceID,
		}

		payload, err := a.authService.RefreshToken(ctx.Request().Context(), req)
		if err != nil {
			return err
		}

		c := &http.Cookie{
			Name:     string(_const.RefreshTokenCookieKey),
			Value:    payload.RefreshToken,
			Path:     "/api/v1/public/auth/session",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			MaxAge:   int(config.RefreshTokenExpiryTime.Seconds()), // max age: second
		}
		ctx.SetCookie(c)

		return util.SuccessResponse(ctx, http.StatusOK, _const.CodeSuccess.Code, "Ok refresh token", model.AuthTokenResponse{AccessToken: payload.AccessToken})
	default:
		// invalid (optional)
		return util.NewError(http.StatusBadRequest, _const.CodeInvalidRequest.Code, "invalid device type header")
	}
}

// Logout godoc
// @Summary      Logout current session
// @Description  Invalidates the current session. Requires authenticated device headers. Mobile clients send refresh_token in the request body. Web clients use the HttpOnly refresh_token cookie.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization  header    string               true   "Bearer access token"
// @Param        Device-ID    header    string               true   "Unique device identifier"
// @Param        Device-Type  header    string               true   "Device type: 1=web, 2=mobile"
// @Param        payload        body      model.LogoutRequest  false  "Mobile logout payload. Web clients use cookie-based refresh_token instead."
// @Success      200            {object}  model.MessageResponseDoc
// @Failure      400            {object}  model.MessageResponseDoc
// @Failure      401            {object}  model.MessageResponseDoc
// @Failure      403            {object}  model.MessageResponseDoc
// @Failure      500            {object}  model.MessageResponseDoc
// @Router       /api/v1/public/auth/session/logout [post]
func (a *authController) Logout(ctx *echo.Context) error {
	deviceType := middleware.GetDeviceType(ctx)
	deviceID := middleware.GetDeviceID(ctx)
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return util.NewErrorByCode(_const.CodeInvalidRequest)
	}

	req := model.LogoutRequest{
		DeviceID: deviceID,
		UserID:   userID,
	}
	switch deviceType {
	case _const.MobileType:
		if err := ctx.Bind(&req); err != nil {
			return util.NewError(http.StatusBadRequest, _const.CodeInvalidRequest.Code, "invalid request body")
		}

		if err := a.authService.Logout(ctx.Request().Context(), req); err != nil {
			return err
		}

		return util.SuccessResponse(ctx, http.StatusOK, _const.CodeSuccess.Code, "logout successful", nil)
	case _const.WebType:
		refreshToken, ok := middleware.GetRefreshToken(ctx)
		if !ok {
			return util.NewErrorByCode(_const.CodeInvalidRequest, "invalid request token")
		}

		req.RefreshToken = refreshToken

		if err := a.authService.Logout(ctx.Request().Context(), req); err != nil {
			return err
		}

		return util.SuccessResponse(ctx, http.StatusOK, _const.CodeSuccess.Code, "logout successful", nil)
	default:
		return util.NewError(http.StatusBadRequest, _const.CodeInvalidRequest.Code, "invalid device type header")
	}
}

// RevokeAllSessions godoc
// @Summary      Revoke all sessions
// @Description  Invalidates all active sessions for the authenticated user using the authenticated access token and device headers.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization  header    string  true  "Bearer access token"
// @Param        Device-ID    header    string  true  "Unique device identifier"
// @Param        Device-Type  header    string  true  "Device type: 1=web, 2=mobile"
// @Success      200            {object}  model.MessageResponseDoc
// @Failure      400            {object}  model.MessageResponseDoc
// @Failure      401            {object}  model.MessageResponseDoc
// @Failure      403            {object}  model.MessageResponseDoc
// @Failure      500            {object}  model.MessageResponseDoc
// @Router       /api/v1/public/auth/session/revoke-all [post]
func (a *authController) RevokeAllSessions(ctx *echo.Context) error {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return util.NewErrorByCode(_const.CodeInvalidRequest)
	}
	accessToken := middleware.GetAccessToken(ctx)
	deviceID := middleware.GetDeviceID(ctx)

	req := model.RevokeAllSessionsRequest{
		DeviceID:    deviceID,
		UserID:      userID,
		AccessToken: accessToken,
	}

	if err := a.authService.RevokeAllSessions(ctx.Request().Context(), req); err != nil {
		return err
	}

	return util.SuccessResponse(ctx, http.StatusOK, _const.CodeSuccess.Code, "Logout all devices successful", nil)
}

func (a *authController) respondByDeviceType(ctx *echo.Context, payload model.AuthTokenResponse, deviceType _const.DeviceType) error {
	switch deviceType {
	case _const.MobileType:
		// Mobile clients receive both tokens in the body and can persist them client-side manually.
		return util.SuccessResponse(ctx, http.StatusOK, _const.CodeSuccess.Code, "login successful", payload)
	case _const.WebType:
		// Web clients store the refresh token in an HttpOnly cookie here.
		// If you want to expand token persistence behavior, add that manually in this branch.
		c := &http.Cookie{
			Name:     string(_const.RefreshTokenCookieKey),
			Value:    payload.RefreshToken,
			Path:     "/api/v1/public/auth/session",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			MaxAge:   int(config.RefreshTokenExpiryTime.Seconds()), // max age: second
		}
		ctx.SetCookie(c)

		return util.SuccessResponse(ctx, http.StatusOK, _const.CodeSuccess.Code, "login successful", model.AuthTokenResponse{AccessToken: payload.AccessToken})
	default:
		// invalid (optional)
		return util.NewError(http.StatusBadRequest, _const.CodeInvalidRequest.Code, "invalid device type header")
	}
}

func NewAuthController(
	authService service.AuthService,
	logger xoxz.XoxzLogger,
) AuthController {
	return &authController{
		authService: authService,
		logger:      logger,
	}
}
