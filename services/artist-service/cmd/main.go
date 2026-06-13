// @title           Artist Service
// @version         1.0
// @description     Artist service providing artist profile management APIs
// @termsOfService  http://swagger.io/terms/

// @contact.name   Artist Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@artist.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /
package main

import (
	"fmt"
	"log"

	_ "github.com/himbo22/xoxz/artist-service/docs"
	"github.com/himbo22/xoxz/artist-service/internal/config"
	"github.com/himbo22/xoxz/artist-service/internal/di"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		if err := godotenv.Load("../.env"); err != nil {
			log.Println("Warning: .env file not found, using system environment variables")
		}
	}

	fmt.Println("Hello, Artist Service!")
	loadConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	app, cleanup, err := di.InitializeApp(loadConfig)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer cleanup()

	if err := app.EchoApp.Start(":" + loadConfig.Server.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
