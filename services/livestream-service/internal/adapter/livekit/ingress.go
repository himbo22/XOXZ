package livekit

import (
	"context"
	"fmt"

	xoxz "github.com/himbo22/xoxz/common-service/xoxz/logger"
	"github.com/himbo22/xoxz/livestream-service/internal/model"
	"github.com/livekit/protocol/livekit"
)

func (c *LiveKitSDK) CreateIngress(ctx context.Context, roomName string) (model.IngressResult, error) {
	// 1. Enforce timeout
	ctx, cancel := context.WithTimeout(ctx, c.Config.Timeout)
	defer cancel()

	// 2. Build request
	req := &livekit.CreateIngressRequest{
		Name:                fmt.Sprintf("%s-ingress", roomName),
		RoomName:            roomName,
		ParticipantIdentity: fmt.Sprintf("obs-%s", roomName),
		ParticipantName:     "OBS Publisher",
		//Protocol:            parseProtocol(protocol),
		Video: &livekit.IngressVideoOptions{
			EncodingOptions: &livekit.IngressVideoOptions_Preset{
				Preset: livekit.IngressVideoEncodingPreset_H264_1080P_30FPS_3_LAYERS,
			},
		},
		Audio: &livekit.IngressAudioOptions{
			EncodingOptions: &livekit.IngressAudioOptions_Preset{
				Preset: livekit.IngressAudioEncodingPreset_OPUS_MONO_64KBS,
			},
		},
	}
	
	// 3. Call SDK with trace-friendly logging
	//c.Logger.DebugContext(ctx, "calling livekit create ingress",
	//	"room", roomName, "protocol", protocol)

	info, err := c.IngressClient.CreateIngress(ctx, req)
	if err != nil {
		c.Logger.Error("livekit create ingress failed", xoxz.String("room", roomName), xoxz.Error(err))
		return model.IngressResult{}, fmt.Errorf("ingress_client: create failed: %w", err)
	}

	// 4. Map -> domain result
	return model.IngressResult{
		ID:        info.IngressId,
		URL:       info.Url,
		StreamKey: info.StreamKey,
	}, nil
}
