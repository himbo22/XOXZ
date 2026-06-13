package bootstrap

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisConfig
type RedisConfig struct {
	Address         string
	Password        string
	DB              int
	IdleTimeout     time.Duration
	MaxConnLifetime time.Duration
	WaitTimeout     time.Duration
	DialTimeout     time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	MaxActive       int
}

func InitRedis(cfg RedisConfig) (*redis.Client, error) {
	if cfg.Address == "" {
		return nil, fmt.Errorf("redis address must not be empty")
	}
	if cfg.MaxActive <= 0 {
		cfg.MaxActive = 20
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:            cfg.Address,
		Password:        cfg.Password,
		DB:              cfg.DB,
		DialTimeout:     cfg.DialTimeout,
		ReadTimeout:     cfg.ReadTimeout,
		WriteTimeout:    cfg.WriteTimeout,
		PoolSize:        cfg.MaxActive,
		ConnMaxIdleTime: cfg.IdleTimeout,
		ConnMaxLifetime: cfg.MaxConnLifetime,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Printf("connected to redis at %s", cfg.Address)
	return rdb, nil
}

func CloseRedis(rdb *redis.Client) {
	if rdb == nil {
		return
	}
	if err := rdb.Close(); err != nil {
		log.Printf("error closing redis connection: %v", err)
	}
}
