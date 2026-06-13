//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/himbo22/xoxz/livestream-service/internal/config"
)

func InitializeApp(cfg *config.Config) (*App, func(), error) {
	wire.Build(
		InfrastructureSet,
		RepositorySet,
		LiveStreamSet,
		ControllerSet,
		NewApp,
	)
	return nil, nil, nil
}
