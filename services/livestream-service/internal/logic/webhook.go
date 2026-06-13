package logic

import (
	"context"

	xoxz "github.com/himbo22/xoxz/common-service/xoxz/logger"
	livekit2 "github.com/himbo22/xoxz/livestream-service/internal/adapter/livekit"
	"github.com/himbo22/xoxz/livestream-service/internal/domain/repository"
)

type WebhookStreamLogic struct {
	livestreamRepo repository.LivestreamRepository
	livekitSDK     *livekit2.LiveKitSDK
	logger         xoxz.XoxzLogger
}

func NewWebhookStreamLogic(
	livestreamRepo repository.LivestreamRepository,
	livekitSDK *livekit2.LiveKitSDK,
	logger xoxz.XoxzLogger,
) *WebhookStreamLogic {
	return &WebhookStreamLogic{
		livekitSDK:     livekitSDK,
		livestreamRepo: livestreamRepo,
		logger:         logger,
	}
}

func (s *WebhookStreamLogic) StartRecording(ctx context.Context, roomName string) error {
	// req := &livekit.TrackCompositeEgressRequest{
	// 	RoomName:     roomName,
	// 	VideoTrackId: "",
	// 	AudioTrackId: "",
	// 	Options: &livekit.TrackCompositeEgressRequest_Preset{
	// 		Preset: livekit.EncodingOptionsPreset_H264_720P_30,
	// 	},
	// 	SegmentOutputs: []*livekit.SegmentedFileOutput{{
	// 		FilenamePrefix:  "my-output",
	// 		PlaylistName:    "my-output.m3u8",
	// 		SegmentDuration: 2,
	// 		Output: &livekit.SegmentedFileOutput_S3{
	// 			S3: &livekit.S3Upload{
	// 				AccessKey:      "accesskey",
	// 				Bucket:         "my-bucket",
	// 				Secret:         "secret",
	// 				Endpoint:       "endpoint",
	// 				ForcePathStyle: true,
	// 			},
	// 		},
	// 	}},
	// }
	// info, err := s.livekitSDK.EgressClient.StartTrackCompositeEgress(ctx, req)
	return nil
}
