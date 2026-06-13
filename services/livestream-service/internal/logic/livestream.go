package logic

import (
	"context"

	"github.com/google/uuid"
	"github.com/himbo22/xoxz/common-service/monitoring/telemetry"
	xoxz "github.com/himbo22/xoxz/common-service/xoxz/logger"
	livekit2 "github.com/himbo22/xoxz/livestream-service/internal/adapter/livekit"
	"github.com/himbo22/xoxz/livestream-service/internal/domain/entity"
	"github.com/himbo22/xoxz/livestream-service/internal/domain/repository"
	"github.com/himbo22/xoxz/livestream-service/internal/model"
)

type LiveStreamLogic struct {
	livestreamRepo repository.LivestreamRepository
	livekitSDK     *livekit2.LiveKitSDK
	logger         xoxz.XoxzLogger
}

func NewLiveStreamLogic(
	livestreamRepo repository.LivestreamRepository,
	livekitSDK *livekit2.LiveKitSDK,
	logger xoxz.XoxzLogger,
) *LiveStreamLogic {
	return &LiveStreamLogic{
		livestreamRepo: livestreamRepo,
		livekitSDK:     livekitSDK,
		logger:         logger,
	}
}

func (l *LiveStreamLogic) CreateStream(ctx context.Context, req model.CreateStreamRequest) (stream *model.StreamResponse, err error) {
	room := &entity.LivestreamRoom{
		RoomName: req.RoomName,
		Status:   string(model.StatusPending),
	}
	if err := l.livestreamRepo.Create(ctx, room); err != nil {
		return nil, err
	}

	resp := &model.StreamResponse{
		StreamID: room.ID.Hex(),
		RoomName: room.RoomName,
		Status:   model.StatusPending,
		WSURL:    l.livekitSDK.Config.Endpoint,
	}

	switch req.Source {
	case model.SourceDirect:
		id, _ := uuid.NewUUID()
		token, err := l.livekitSDK.GetPublisherToken(room.RoomName, id.String())
		if err != nil {
			return nil, err
		}

		resp.PublisherToken = token
	// create ingress for obs/stream lab/ ...
	case model.SourceIngress:
		ctx, end := telemetry.StartSpan(ctx, "test1", "logic-create-stream")
		defer func() { end(err) }()

		if req.Protocol == "" {
			req.Protocol = "whip" // default
		}

		ingressResult, err := l.livekitSDK.CreateIngress(ctx, req.RoomName)
		if err != nil {
			return nil, err
		}

		resp.IngressURL = ingressResult.URL
		resp.StreamKey = ingressResult.StreamKey
	}

	return resp, nil
}
