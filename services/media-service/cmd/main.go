package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/himbo22/xoxz/media-service/internal/config"
	"github.com/himbo22/xoxz/media-service/internal/di"
	"github.com/labstack/echo/v5"
)

func main() {
	config := config.LoadConfig()
	app, cleanup, err := di.InitializeApp(config)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer cleanup()

	// 1. Run HTTP Server (Echo) in a separate Goroutine
	// 1. [V5 WAY] Create Context to catch OS signals (Ctrl+C, Docker stop)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 2. Run gRPC Server in a separate Goroutine (runs in parallel)
	go func() {
		log.Println("Starting gRPC Server on port :50051")
		if err := app.GrpcServer.Start("50051"); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	// 3. Configure HTTP Server using Echo v5 StartConfig
	sc := echo.StartConfig{
		Address:         ":" + config.Server.Port,
		GracefulTimeout: 10 * time.Second, // Auto-wait up to 10s on shutdown
	}

	log.Printf("Starting Echo HTTP Server on port :%s", config.Server.Port)

	// 4. Start HTTP Server.
	// Start() will BLOCK here to serve requests.
	// The magic: When you press Ctrl+C, `ctx` is cancelled. Start() catches the signal,
	// performs Graceful Shutdown (waits 10s), then returns nil to proceed.
	if err := sc.Start(ctx, app.EchoApp); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP server error: %v", err)
	}

	log.Println("HTTP Server shut down safely. Proceeding to stop gRPC...")

	// 5. After HTTP has shut down safely (past sc.Start), proceed to stop gRPC
	app.GrpcServer.Stop()

	log.Println("Entire system shut down successfully. Goodbye!")
}
