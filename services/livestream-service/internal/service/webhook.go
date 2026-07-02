package service

import (
	"context"

	xoxz "github.com/himbo22/xoxz/common-service/xoxz/logger"
	"github.com/himbo22/xoxz/livestream-service/internal/domain/repository"
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
	livestreamRepo repository.LivestreamRepository
	livekitClient  LiveKitProvider
	logger         xoxz.XoxzLogger
}

func NewWebhookService(
	livestreamRepo repository.LivestreamRepository,
	livekitClient LiveKitProvider,
	logger xoxz.XoxzLogger,
) WebhookService {
	return &webhookService{
		livekitClient:  livekitClient,
		livestreamRepo: livestreamRepo,
		logger:         logger,
	}
}

func (s *webhookService) OnStreamStarted(ctx context.Context, roomName string) error {
	return nil
}

func (s *webhookService) OnStreamEnded(ctx context.Context, roomName string) error {
	return nil
}

func (s *webhookService) OnViewerJoined(ctx context.Context, roomName string, identity string) error {
	return nil
}

func (s *webhookService) OnViewerLeft(ctx context.Context, roomName string, identity string) error {
	return nil
}

func (s *webhookService) OnRecordStarted(ctx context.Context, roomName string, egressID string) error {
	return nil
}

func (s *webhookService) OnRecordEnded(ctx context.Context, roomName string, egressID string, videoURL string) error {
	return nil
}
