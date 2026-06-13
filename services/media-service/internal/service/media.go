package service

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/himbo22/xoxz/media-service/internal/adapter/storage"
	_const "github.com/himbo22/xoxz/media-service/internal/const"
	"github.com/himbo22/xoxz/media-service/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mediaService struct {
	MinioClient *storage.MinioClient
}

type MediaService interface {
	GeneratePresignedURL(ctx context.Context, req model.GenerateURLRequest) (*model.GenerateURLResponse, error)
	CommitFile(ctx context.Context, tmpPath string, perPath string) error
	DeleteFile(ctx context.Context, objectPath string) error
}

// GeneratePresignedURL implements [MediaService].
func (m *mediaService) GeneratePresignedURL(ctx context.Context, req model.GenerateURLRequest) (*model.GenerateURLResponse, error) {
	// check content type
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
	}

	if !allowedTypes[req.ContentType] {
		return nil, status.Error(codes.InvalidArgument, "File format not supported")
	}

	// 1. Auto-generate file name to prevent duplicates and enhance security
	// Example: uuid + original file extension
	ext := filepath.Ext(req.FileName)
	safeFileName := fmt.Sprintf("%s_%s%s", req.UserID, uuid.New().String(), ext)

	// 2. Create temporary (TmpPath) and permanent (PerPath) paths
	// Example: req.Module = "profile" -> tmp/profile/user1_uuid.jpg
	tmpPath := fmt.Sprintf("tmp/%s/%s", _const.MediaAvatar, safeFileName)
	perPath := fmt.Sprintf("per/%s/%s", _const.MediaAvatar, safeFileName)

	// 3. Set presigned URL expiry (e.g., 15 minutes)
	expiry := time.Minute * 15

	// 4. Enable Security: Enforce Content-Type
	// This forces the Frontend to send the correct Content-Type header on upload
	reqParams := make(url.Values)
	reqParams.Set("response-content-type", req.ContentType) // Enforce correct content type

	// 5. CALL MINIO SDK PUT METHOD
	presignedURL, err := m.MinioClient.MinioClient.PresignedPutObject(ctx, m.MinioClient.BucketName, tmpPath, expiry)
	if err != nil {
		log.Printf("Error generating upload URL: %v", err)
		return nil, fmt.Errorf("cannot generate upload URL: %w", err)
	}

	// 6. Return DTO to Frontend
	return &model.GenerateURLResponse{
		UploadURL: presignedURL.String(),
		TmpPath:   tmpPath,
		PerPath:   perPath,
	}, nil
}

// CommitFile implements [MediaService].
func (m *mediaService) CommitFile(ctx context.Context, tmpPath string, perPath string) error {
	panic("unimplemented")
}

// DeleteFile implements [MediaService].
func (m *mediaService) DeleteFile(ctx context.Context, objectPath string) error {
	panic("unimplemented")
}

func NewMediaService(
	MinioClient *storage.MinioClient,
) MediaService {
	return &mediaService{
		MinioClient: MinioClient,
	}
}
