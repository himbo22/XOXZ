package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	MinioClient *minio.Client
	BucketName  string
}

type Config struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	BucketName      string
}

func NewClient(cfg Config) (*MinioClient, error) {
	// 1. Initialize Object Client
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot initialize MinIO config: %w", err)
	}

	// 2. [Fail-Fast Principle] Check actual network connectivity
	// Ping the server to check if the bucket exists. Set 5 second timeout to prevent hanging.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exists, err := client.BucketExists(ctx, cfg.BucketName)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to MinIO server (Endpoint: %s): %w", cfg.Endpoint, err)
	}
	if !exists {
		// Could auto-create bucket here, but minio-setup already handles that,
		// so if bucket is missing, we error out.
		return nil, fmt.Errorf("MinIO connected but bucket '%s' does not exist", cfg.BucketName)
	}

	// 3. Return fully initialized Client
	return &MinioClient{
		MinioClient: client,
		BucketName:  cfg.BucketName,
	}, nil
}

func (c *MinioClient) CommitFile(ctx context.Context, tmpPath string, perPath string) error {
	// 1. Configure Source
	srcOpts := minio.CopySrcOptions{
		Bucket: c.BucketName,
		Object: tmpPath,
	}

	// 2. Configure Destination
	dstOpts := minio.CopyDestOptions{
		Bucket: c.BucketName,
		Object: perPath,
	}

	// 3. Execute internal Copy command
	_, err := c.MinioClient.CopyObject(ctx, dstOpts, srcOpts)
	if err != nil {
		return fmt.Errorf("cannot copy file to permanent storage: %w", err)
	}

	// 4. Delete file from temp storage
	err = c.MinioClient.RemoveObject(ctx, c.BucketName, tmpPath, minio.RemoveObjectOptions{})
	if err != nil {
		// [KEY PRINCIPLE]
		// Don't return error here! The file is ALREADY SAFELY STORED in per/.
		// tmp deletion failure is just garbage cleanup, we only need to log a warning.
		// If deletion fails, MinIO ILM will auto-clean the next day.
		// c.logger.Warn("Failed to delete temp file after commit (ILM will auto-clean)",
		// 	zap.String("tmpPath", tmpPath),
		// 	zap.Error(err),
		// )
		log.Printf("Failed to delete temp file after commit (ILM will auto-clean): %v", err)
	}

	return nil
}
