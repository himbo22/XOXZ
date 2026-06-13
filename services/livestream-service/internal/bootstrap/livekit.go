package bootstrap

import (
	"errors"
	"time"

	xoxz "github.com/himbo22/xoxz/common-service/xoxz/logger"
	"github.com/himbo22/xoxz/livestream-service/internal/adapter/livekit"
	"github.com/himbo22/xoxz/livestream-service/internal/config"
	lksdk "github.com/livekit/server-sdk-go/v2"
)

func NewLiveKitSDK(cfg config.LiveKitConfig, logger xoxz.XoxzLogger) (*livekit.LiveKitSDK, error) {
	if cfg.APIKey == "" || cfg.APISecret == "" || cfg.Endpoint == "" {
		return nil, errors.New("ingress_client: missing required config (api_key, api_secret, endpoint)")
	}

	if cfg.Timeout == 0 {
		cfg.Timeout = 5 * time.Second // Default production timeout
	}

	ingressClient := lksdk.NewIngressClient(cfg.Endpoint, cfg.APIKey, cfg.APISecret)
	egressClient := lksdk.NewEgressClient(cfg.Endpoint, cfg.APIKey, cfg.APISecret)
	roomClient := lksdk.NewRoomServiceClient(cfg.Endpoint, cfg.APIKey, cfg.APISecret)

	return &livekit.LiveKitSDK{
		RoomClient:    roomClient,
		IngressClient: ingressClient,
		EgressClient:  egressClient,
		Config:        cfg,
		Logger:        logger,
	}, nil
}
