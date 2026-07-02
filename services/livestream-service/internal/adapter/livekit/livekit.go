package livekit

import (
	"context"
	"time"

	"github.com/livekit/protocol/auth"
)

func (c *LiveKitSDK) GetPublisherToken(ctx context.Context, room, identity string) (string, error) {
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
