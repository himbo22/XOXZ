package service

import (
	"context"

	"github.com/himbo22/xoxz/account-service/internal/bootstrap"
	xoxz "github.com/himbo22/xoxz/common-service/xoxz/logger"
)

type AdminService interface {
	CreateArtistInviteRequest(ctx context.Context) error
	RevokeArtistAccount(ctx context.Context) error
}

type adminService struct {
	permissions *bootstrap.AccessControl
	logger      xoxz.XoxzLogger
}

func NewAdminService(
	logger xoxz.XoxzLogger,
	permissions *bootstrap.AccessControl,
) AdminService {
	return &adminService{
		permissions: permissions,
		logger:      logger,
	}
}

func (a *adminService) CreateArtistInviteRequest(ctx context.Context) error {
	return nil
}

func (a *adminService) RevokeArtistAccount(ctx context.Context) error {
	return nil
}
