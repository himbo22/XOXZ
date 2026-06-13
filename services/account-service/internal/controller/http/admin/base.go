package admin

import (
	_const "github.com/himbo22/xoxz/account-service/internal/const"
	"github.com/himbo22/xoxz/account-service/internal/model"
	"github.com/himbo22/xoxz/account-service/internal/service"
	"github.com/himbo22/xoxz/account-service/internal/util"
	xoxz "github.com/himbo22/xoxz/common-service/xoxz/logger"
	"github.com/labstack/echo/v5"
)

type adminController struct {
	adminService service.AdminService
	logger       xoxz.XoxzLogger
}

func (a *adminController) CreateArtistInvite(c *echo.Context) error {
	//TODO implement me
	req := model.CreateArtistInviteRequest{}
	if err := c.Bind(&req); err != nil {
		return util.NewErrorByCode(_const.CodeInvalidRequest)
	}
	// validate

	err := a.adminService.CreateArtistInviteRequest(c.Request().Context())
	if err != nil {
		return util.NewErrorByCode(_const.CodeInvalidRequest)
	}

	return util.SuccessResponseByCode(c, _const.CodeSuccess, nil)
}

func (a *adminController) RevokeArtistAccount(c *echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func NewAdminController(
	adminService service.AdminService,
	logger xoxz.XoxzLogger,
) AdminController {
	return &adminController{
		adminService: adminService,
		logger:       logger,
	}
}
