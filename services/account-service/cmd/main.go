// @title           IAM Service
// @version         1.0
// @description     Identity and Access Management service providing authentication, authorization, and user management APIs
// @termsOfService  http://swagger.io/terms/

// @contact.name   IAM Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@iam.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"fmt"
	"log"

	_ "github.com/himbo22/xoxz/account-service/docs" // swagger docs
	"github.com/himbo22/xoxz/account-service/internal/config"
	"github.com/himbo22/xoxz/account-service/internal/di"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		if err := godotenv.Load("../.env"); err != nil {
			log.Println("Warning: .env file not found, using system environment variables")
		}
	}

	fmt.Println("Hello, Account Service!")
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

	// go func() {
	// 	app.EchoApp.Start(":" + config.Server.Port)
	// }()
}
