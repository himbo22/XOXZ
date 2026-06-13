package livekit

import (
	"time"

	xoxz "github.com/himbo22/xoxz/common-service/xoxz/logger"
	"github.com/himbo22/xoxz/livestream-service/internal/config"
	"github.com/livekit/protocol/auth"
	lksdk "github.com/livekit/server-sdk-go/v2"
)

type LiveKitSDK struct {
	RoomClient    *lksdk.RoomServiceClient
	IngressClient *lksdk.IngressClient
	EgressClient  *lksdk.EgressClient
	Config        config.LiveKitConfig
	Logger        xoxz.XoxzLogger
}

func (c *LiveKitSDK) GetPublisherToken(room, identity string) (string, error) {
	at := auth.NewAccessToken(c.Config.APIKey, c.Config.APISecret)
	grant := &auth.VideoGrant{
		RoomJoin: true,
		Room:     room,
	}
	at.SetVideoGrant(grant).
		SetIdentity(identity).
		SetValidFor(time.Hour)

	return at.ToJWT()
}
