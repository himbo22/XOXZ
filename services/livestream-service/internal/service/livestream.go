package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/himbo22/xoxz/common-service/monitoring/telemetry"
	xoxz "github.com/himbo22/xoxz/common-service/xoxz/logger"
	"github.com/himbo22/xoxz/livestream-service/internal/domain/entity"
	"github.com/himbo22/xoxz/livestream-service/internal/domain/repository"
	"github.com/himbo22/xoxz/livestream-service/internal/model"
)

type LiveStreamService interface {
	Create(ctx context.Context, req model.CreateStreamRequest) (stream *model.StreamResponse, err error)
	Stop(ctx context.Context) error
}

type LiveKitProvider interface {
	GetPublisherToken(ctx context.Context, room, identity string) (string, error)
	CreateIngress(ctx context.Context, roomName string) (model.IngressResult, error)
	GetEndpoint() string
}

type liveStreamService struct {
	livestreamRepo repository.LivestreamRepository
	livekitClient  LiveKitProvider
	logger         xoxz.XoxzLogger
}

func NewLiveStreamService(
	livestreamRepo repository.LivestreamRepository,
	livekitClient LiveKitProvider,
	logger xoxz.XoxzLogger,
) LiveStreamService {
	return &liveStreamService{
		livestreamRepo: livestreamRepo,
		livekitClient:  livekitClient,
		logger:         logger,
	}
}

func (l *liveStreamService) Create(ctx context.Context, req model.CreateStreamRequest) (stream *model.StreamResponse, err error) {
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
		WSURL:    l.livekitClient.GetEndpoint(),
	}

	switch req.Source {
	case model.SourceDirect:
		id, _ := uuid.NewUUID()
		token, err := l.livekitClient.GetPublisherToken(ctx, room.RoomName, id.String())
		if err != nil {
			return nil, err
		}

		resp.PublisherToken = token
	case model.SourceIngress:
		ctx, end := telemetry.StartSpan(ctx, "test1", "usecase-create-stream")
		defer func() { end(err) }()

		if req.Protocol == "" {
			req.Protocol = "whip"
		}

		ingressResult, err := l.livekitClient.CreateIngress(ctx, req.RoomName)
		if err != nil {
			return nil, err
		}

		resp.IngressURL = ingressResult.URL
		resp.StreamKey = ingressResult.StreamKey
	}

	return resp, nil
}

func (l *liveStreamService) Stop(ctx context.Context) error {
	return nil
}
