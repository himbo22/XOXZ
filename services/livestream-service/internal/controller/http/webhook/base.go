package webhook

import (
	"fmt"

	xoxz "github.com/himbo22/xoxz/common-service/xoxz/logger"
	"github.com/himbo22/xoxz/livestream-service/internal/config"
	"github.com/himbo22/xoxz/livestream-service/internal/util"
	"github.com/labstack/echo/v5"
	"github.com/livekit/protocol/auth"
	"github.com/livekit/protocol/webhook"
)

type webhookController struct {
	receiver auth.KeyProvider
	logger   xoxz.XoxzLogger
}

func NewWebhookController(cfg *config.Config, logger xoxz.XoxzLogger) WebhookController {
	keyProvider := auth.NewSimpleKeyProvider(cfg.LiveKit.APIKey, cfg.LiveKit.APISecret)
	return &webhookController{
		receiver: keyProvider,
		logger:   logger,
	}
}

func (w *webhookController) ServeHTTP(ctx *echo.Context) error {
	event, err := webhook.ReceiveWebhookEvent(ctx.Request(), w.receiver)
	if err != nil {
		w.logger.Errorf("Webhook verify failed: %v\n", err)
		return err
	}
	//
	//var roomName string
	//if event.Room != nil {
	//	roomName = event.Room.Name
	//} else if event.IngressInfo != nil {
	//	roomName = event.IngressInfo.RoomName
	//}
	//
	//if roomName == "" {
	//	w.logger.Warn("Webhook received without room info", xoxz.String("event", event.Event))
	//	return util.SuccessResponse(ctx, 200, 123, "webhook received", nil)
	//}
	//
	//switch event.Event {
	//case "ingress_started":
	//	// OBS connected successfully, Ingress started pushing stream
	//	fmt.Printf("STREAM STARTED: Room %s is now LIVE!\n", roomName)
	//
	//	// Call Repository Update MongoDB -> Status: "LIVE"
	//	// ctrl.useCase.UpdateStatus(ctx, roomName, "LIVE")
	//
	//case "ingress_ended":
	//	// OBS disconnected
	//	fmt.Printf("STREAM ENDED: Room %s has ENDED!\n", roomName)
	//
	//	// Call Repository Update MongoDB -> Status: "ENDED"
	//	// ctrl.useCase.UpdateStatus(ctx, roomName, "ENDED")
	//
	//case "participant_joined":
	//	// This event fires when a viewer or idol enters the room
	//	// Used for tracking viewer count
	//	fmt.Printf("STREAM ENDED: Room %s has ENDED!\n", "roomName")
	//
	//default:
	//	// Skip non-critical events (room_created, track_published...)
	//	fmt.Printf("Ignored event: %s\n", event.Event)
	//}
	fmt.Printf("Event: %s\n", event.Event)
	return util.SuccessResponse(ctx, 200, 123, "webhook received", nil)
}
