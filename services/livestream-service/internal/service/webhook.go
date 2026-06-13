package service

import (
	"context"

	"github.com/himbo22/xoxz/livestream-service/internal/logic"
)

type WebhookService interface {
	OnStreamStarted(ctx context.Context, roomName string) error
	OnStreamEnded(ctx context.Context, roomName string) error
	OnViewerJoined(ctx context.Context, roomName string, identity string) error
	OnViewerLeft(ctx context.Context, roomName string, identity string) error
	OnRecordStarted(ctx context.Context, roomName string, egressID string) error
	OnRecordEnded(ctx context.Context, roomName string, egressID string, videoURL string) error
}

type webhookService struct {
	webhookLogic *logic.WebhookStreamLogic
}

func NewWebhookService(webhookLogic *logic.WebhookStreamLogic) WebhookService {
	return &webhookService{
		webhookLogic: webhookLogic,
	}
}

func (w *webhookService) OnStreamStarted(ctx context.Context, roomName string) error {
	//TODO implement me
	panic("implement me")
}

func (w *webhookService) OnStreamEnded(ctx context.Context, roomName string) error {
	//TODO implement me
	panic("implement me")
}

func (w *webhookService) OnViewerJoined(ctx context.Context, roomName string, identity string) error {
	//TODO implement me
	panic("implement me")
}

func (w *webhookService) OnViewerLeft(ctx context.Context, roomName string, identity string) error {
	//TODO implement me
	panic("implement me")
}

func (w *webhookService) OnRecordStarted(ctx context.Context, roomName string, egressID string) error {
	//TODO implement me
	panic("implement me")
}

func (w *webhookService) OnRecordEnded(ctx context.Context, roomName string, egressID string, videoURL string) error {
	//TODO implement me
	panic("implement me")
}
