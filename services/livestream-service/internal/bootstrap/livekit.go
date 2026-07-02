package bootstrap

import (
	"fmt"
	"time"

	xoxz "github.com/himbo22/xoxz/common-service/xoxz/logger"
	"github.com/himbo22/xoxz/livestream-service/internal/adapter/livekit"
	"github.com/himbo22/xoxz/livestream-service/internal/config"
	twirp "github.com/livekit/protocol/utils/xtwirp"
	lksdk "github.com/livekit/server-sdk-go/v2"
)

func NewLiveKitSDK(cfg config.LiveKitConfig, logger xoxz.XoxzLogger) (*livekit.LiveKitSDK, error) {
	if cfg.APIKey == "" || cfg.APISecret == "" || cfg.Endpoint == "" {
		return nil, fmt.Errorf("livekit: missing required config (api_key, api_secret, endpoint)")
	}

	if cfg.Timeout == 0 {
		cfg.Timeout = 5 * time.Second // Default production timeout
	}

	// http.Client dùng chung cho cả 3 client, áp timeout thực sự
	//httpClient := &http.Client{
	//	Timeout: cfg.Timeout,
	//	Transport: &http.Transport{
	//		MaxIdleConns:        50,
	//		MaxIdleConnsPerHost: 10,
	//		IdleConnTimeout:     90 * time.Second,
	//	},
	//}

	clientOpt := twirp.DefaultClientOptions()

	ingressClient := lksdk.NewIngressClient(cfg.Endpoint, cfg.APIKey, cfg.APISecret, clientOpt...)
	egressClient := lksdk.NewEgressClient(cfg.Endpoint, cfg.APIKey, cfg.APISecret, clientOpt...)
	roomClient := lksdk.NewRoomServiceClient(cfg.Endpoint, cfg.APIKey, cfg.APISecret, clientOpt...)

	return &livekit.LiveKitSDK{
		RoomClient:    roomClient,
		IngressClient: ingressClient,
		EgressClient:  egressClient,
		Config:        cfg,
		Logger:        logger,
	}, nil
}

//func CloseLiveKit(sdk *livekit.LiveKitSDK) {
//	sdk.In
//}
