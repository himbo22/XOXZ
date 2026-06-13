package bootstrap

import (
	"context"
	"log"
	"time"

	"github.com/himbo22/xoxz/livestream-service/internal/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func InitMongoDB(cfg config.MongoDBConfig) (*mongo.Database, func(), error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	clientOptions := options.Client().
		ApplyURI(cfg.URI).
		SetMaxPoolSize(cfg.MaxPoolSize).
		SetMinPoolSize(cfg.MinPoolSize).
		SetMaxConnIdleTime(time.Duration(cfg.MaxConnsIdleTime) * time.Minute)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, nil, err
	}

	// Ping check
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, nil, err
	}

	cleanup := func() {
		ctxDisconnect, cancelDisconnect := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancelDisconnect()

		if err := client.Disconnect(ctxDisconnect); err != nil {
			log.Printf("mongo disconnect error: %v", err)
		}
	}

	db := client.Database(cfg.DBName)

	return db, cleanup, nil
}
