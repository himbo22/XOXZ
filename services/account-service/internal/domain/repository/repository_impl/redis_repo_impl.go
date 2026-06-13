package repository_impl

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	_const "github.com/himbo22/xoxz/account-service/internal/const"
	"github.com/himbo22/xoxz/account-service/internal/domain/repository"
	"github.com/himbo22/xoxz/account-service/internal/model"
	xoxz "github.com/himbo22/xoxz/common-service/xoxz/logger"
	"github.com/redis/go-redis/v9"
)

type redisRepository struct {
	client *redis.Client
	logger xoxz.XoxzLogger
}

func NewRedisRepository(client *redis.Client, logger xoxz.XoxzLogger) repository.RedisRepository {
	return &redisRepository{
		client: client,
		logger: logger,
	}
}

// RevokeAllSessions implements repository.RedisRepository.
func (r *redisRepository) RevokeAllSessions(ctx context.Context, userID uuid.UUID) error {
	trackerKey := _const.UserDevicesKey(userID)

	devices, err := r.client.HGetAll(ctx, trackerKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	// Initialize a transaction pipeline.
	pipe := r.client.TxPipeline()

	// remove all refresh token from user
	for _, deviceID := range devices {
		sessionKey := _const.RefreshTokenKey(deviceID)
		pipe.Del(ctx, sessionKey)
	}

	// remove all devices
	pipe.Del(ctx, trackerKey)

	// add AT to blacklist
	blKey := _const.BlackAccessTokenKey(userID)
	pipe.Set(ctx, blKey, time.Now().Unix(), 15*time.Minute)

	_, err = pipe.Exec(ctx)
	if err != nil {
		r.logger.Error("Failed to save session atomically",
			xoxz.String("userID", userID.String()),
			xoxz.Error(err),
		)
		return err
	}

	return nil
}

// RemoveSession implements repository.RedisRepository.
func (r *redisRepository) RemoveSession(ctx context.Context, userID uuid.UUID, deviceID string, refreshToken string) error {
	sessionKey := _const.RefreshTokenKey(refreshToken)
	trackerKey := _const.UserDevicesKey(userID)

	pipe := r.client.Pipeline()

	pipe.Del(ctx, sessionKey)

	pipe.HDel(ctx, trackerKey, deviceID)

	_, err := pipe.Exec(ctx)
	if err != nil {
		r.logger.Error("Failed to delete Redis key",
			xoxz.String("key", sessionKey),
			xoxz.Error(err))
		return err
	}

	r.logger.Info("Successfully deleted Redis key",
		xoxz.String("key", sessionKey))

	return nil
}

func (r *redisRepository) CreateGoogleSession(ctx context.Context, req model.SessionData) error {
	sessionKey := _const.RefreshTokenKey(req.RefreshToken)
	trackerKey := _const.UserDevicesKey(req.UserID)
	googleTokenKey := _const.GoogleTokenKey(req.HashedToken)

	// Initialize a transaction pipeline.
	pipe := r.client.TxPipeline()

	// Find old token
	oldToken, err := r.client.HGet(ctx, trackerKey, req.DeviceID).Result()
	// delete old if exists
	if err == nil && oldToken != "" {
		stringKeyOld := _const.RefreshTokenKey(oldToken)
		pipe.Del(ctx, stringKeyOld)
	}

	// update new token in tracker
	pipe.HSet(ctx, trackerKey, req.DeviceID, req.RefreshToken)

	// replace by new one
	pipe.Set(ctx, sessionKey, req.DataSession, req.SessionExpiration)

	// set google token -> prevent reply attack
	pipe.Set(ctx, googleTokenKey, 1, req.HashedTokenExpiration)

	// 3. Execute all commands in one transaction (MULTI -> cmds -> EXEC).
	_, err = pipe.Exec(ctx)
	if err != nil {
		r.logger.Error("Failed to save session atomically",
			xoxz.String("userID", req.UserID.String()),
			xoxz.Error(err),
		)
		return err
	}

	r.logger.Info("Successfully set Redis key",
		xoxz.String("key", sessionKey),
		xoxz.String("refreshToken", req.RefreshToken),
		xoxz.String("expiration", req.SessionExpiration.String()))

	return nil
}

// redis pipeline: transaction
func (r *redisRepository) CreateSession(ctx context.Context, req model.SessionData) error {
	sessionKey := _const.RefreshTokenKey(req.RefreshToken)
	trackerKey := _const.UserDevicesKey(req.UserID)

	// Initialize a transaction pipeline.
	pipe := r.client.TxPipeline()

	// Find old token
	oldToken, err := r.client.HGet(ctx, trackerKey, req.DeviceID).Result()
	// delete old if exists
	if err == nil && oldToken != "" {
		stringKeyOld := _const.RefreshTokenKey(oldToken)
		pipe.Del(ctx, stringKeyOld)
	}

	// update new token in tracker
	pipe.HSet(ctx, trackerKey, req.DeviceID, req.RefreshToken)

	// replace by new one
	pipe.Set(ctx, sessionKey, req.DataSession, req.SessionExpiration)

	// 3. Execute all commands in one transaction (MULTI -> cmds -> EXEC).
	_, err = pipe.Exec(ctx)
	if err != nil {
		r.logger.Error("Failed to save session atomically",
			xoxz.String("userID", req.UserID.String()),
			xoxz.Error(err),
		)
		return err
	}

	r.logger.Info("Successfully set Redis key",
		xoxz.String("key", sessionKey),
		xoxz.String("refreshToken", req.RefreshToken),
		xoxz.String("expiration", req.SessionExpiration.String()))

	return nil
}

// Set sets a key-value pair with expiration
func (r *redisRepository) Set(ctx context.Context, key, value string, expiration time.Duration) error {
	err := r.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		r.logger.Error("Failed to set Redis key",
			xoxz.String("key", key),
			xoxz.Error(err))
		return err
	}

	r.logger.Info("Successfully set Redis key",
		xoxz.String("key", key),
		xoxz.String("expiration", expiration.String()))

	return nil
}

// Get gets a value by key
func (r *redisRepository) Get(ctx context.Context, key string) (string, error) {
	value, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		r.logger.Error("Failed to get Redis key",
			xoxz.String("key", key),
			xoxz.Error(err))
		return "", err
	}

	r.logger.Info("Successfully got Redis key",
		xoxz.String("key", key))

	return value, nil
}

// Del deletes a key
func (r *redisRepository) Del(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		r.logger.Error("Failed to delete Redis key",
			xoxz.String("key", key),
			xoxz.Error(err))
		return err
	}

	r.logger.Info("Successfully deleted Redis key",
		xoxz.String("key", key))

	return nil
}

// Exists checks if a key exists
func (r *redisRepository) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		r.logger.Error("Failed to check Redis key existence",
			xoxz.String("key", key),
			xoxz.Error(err))
		return false, err
	}

	return exists > 0, nil
}

// Expire sets expiration for a key
func (r *redisRepository) Expire(ctx context.Context, key string, expiration time.Duration) error {
	err := r.client.Expire(ctx, key, expiration).Err()
	if err != nil {
		r.logger.Error("Failed to set expiration for Redis key",
			xoxz.String("key", key),
			xoxz.Error(err))
		return err
	}

	r.logger.Info("Successfully set expiration for Redis key",
		xoxz.String("key", key),
		xoxz.String("expiration", expiration.String()))

	return nil
}
