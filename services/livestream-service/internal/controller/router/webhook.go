package router

import (
	"github.com/himbo22/xoxz/livestream-service/internal/controller/http/webhook"
	"github.com/labstack/echo/v5"
)

func SetupWebhook(parentRouter *echo.Group, webhookController webhook.WebhookController) {
	group := parentRouter.Group("/webhook")
	{
		group.POST("/livekit", webhookController.ServeHTTP)
	}
}
