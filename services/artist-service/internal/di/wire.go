//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/himbo22/xoxz/artist-service/internal/config"
)

func InitializeApp(cfg *config.Config) (*App, func(), error) {
	wire.Build(
		InfrastructureSet,
		ArtistSet,
		RepositorySet,
		ControllerSet,
		NewApp,
	)
	return nil, nil, nil
}
