package service

import (
	"context"

	"github.com/himbo22/xoxz/livestream-service/internal/logic"
	"github.com/himbo22/xoxz/livestream-service/internal/model"
)

type LiveStreamService interface {
	Create(ctx context.Context, req model.CreateStreamRequest) (stream *model.StreamResponse, err error)
	Stop(ctx context.Context) error
}

type liveStreamService struct {
	liveStreamLogic *logic.LiveStreamLogic
}

func (l *liveStreamService) Create(ctx context.Context, req model.CreateStreamRequest) (stream *model.StreamResponse, err error) {
	return l.liveStreamLogic.CreateStream(ctx, req)
}

func (l *liveStreamService) Stop(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func NewLiveStreamService(liveStreamLogic *logic.LiveStreamLogic) LiveStreamService {
	return &liveStreamService{
		liveStreamLogic: liveStreamLogic,
	}
}
