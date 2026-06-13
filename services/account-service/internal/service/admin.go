package service

import (
	"context"

	"github.com/himbo22/xoxz/account-service/internal/logic"
)

type AdminService interface {
	CreateArtistInviteRequest(ctx context.Context) error
	RevokeArtistAccount(ctx context.Context) error
}

type adminService struct {
	adminLogic *logic.AdminLogic
}

func (a *adminService) CreateArtistInviteRequest(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (a *adminService) RevokeArtistAccount(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func NewAdminService(adminLogic *logic.AdminLogic) AdminService {
	return &adminService{
		adminLogic: adminLogic,
	}
}
