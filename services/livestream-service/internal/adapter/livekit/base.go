package livekit

import (
	xoxz "github.com/himbo22/xoxz/common-service/xoxz/logger"
	"github.com/himbo22/xoxz/livestream-service/internal/config"
	lksdk "github.com/livekit/server-sdk-go/v2"
)

type LiveKitSDK struct {
	RoomClient    *lksdk.RoomServiceClient
	IngressClient *lksdk.IngressClient
	EgressClient  *lksdk.EgressClient
	Config        config.LiveKitConfig
	Logger        xoxz.XoxzLogger
}

func (c *LiveKitSDK) GetEndpoint() string {
	return c.Config.Endpoint
}
