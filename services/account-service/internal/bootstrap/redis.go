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
	IdleTimeout     string
	MaxConnLifetime string
	WaitTimeout     string
	DialTimeout     string
	ReadTimeout     string
	WriteTimeout    string
	MaxActive       int
}

func mustParseDuration(name, value string, def time.Duration) (time.Duration, error) {
	if value == "" {
		return def, nil
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", name, err)
	}
	return duration, nil
}

func InitRedis(cfg RedisConfig) (*redis.Client, error) {
	if cfg.Address == "" {
		return nil, fmt.Errorf("redis address must not be empty")
	}
	if cfg.MaxActive <= 0 {
		cfg.MaxActive = 20
	}

	dialTimeout, err := mustParseDuration("DialTimeout", cfg.DialTimeout, 5*time.Second)
	if err != nil {
		return nil, err
	}
	readTimeout, err := mustParseDuration("ReadTimeout", cfg.ReadTimeout, 3*time.Second)
	if err != nil {
		return nil, err
	}
	writeTimeout, err := mustParseDuration("WriteTimeout", cfg.WriteTimeout, 3*time.Second)
	if err != nil {
		return nil, err
	}
	idleTimeout, err := mustParseDuration("IdleTimeout", cfg.IdleTimeout, 5*time.Minute)
	if err != nil {
		return nil, err
	}
	maxConnLifetime, err := mustParseDuration("MaxConnLifetime", cfg.MaxConnLifetime, 30*time.Minute)
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:            cfg.Address,
		Password:        cfg.Password,
		DB:              cfg.DB,
		DialTimeout:     dialTimeout,
		ReadTimeout:     readTimeout,
		WriteTimeout:    writeTimeout,
		PoolSize:        cfg.MaxActive,
		ConnMaxIdleTime: idleTimeout,
		ConnMaxLifetime: maxConnLifetime,
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
