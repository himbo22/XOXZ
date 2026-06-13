package logic

import (
	"github.com/himbo22/xoxz/account-service/internal/bootstrap"
	xoxz "github.com/himbo22/xoxz/common-service/xoxz/logger"
)

type AdminLogic struct {
	permissions *bootstrap.AccessControl
	logger      xoxz.XoxzLogger
}

func NewAdminLogic(
	logger xoxz.XoxzLogger,
	permissions *bootstrap.AccessControl,
) *AdminLogic {
	return &AdminLogic{
		permissions: permissions,
		logger:      logger,
	}
}
