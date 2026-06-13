package grpc

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	media "github.com/himbo22/xoxz/common-service/protobuf/media"
	"github.com/himbo22/xoxz/media-service/internal/adapter/storage"
	"github.com/minio/minio-go/v7"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MediaHandler struct {
	media.UnimplementedMediaServiceServer
	storage *storage.MinioClient
}

func NewMediaHandler(storage *storage.MinioClient) *MediaHandler {
	return &MediaHandler{
		storage: storage,
	}
}

func (h *MediaHandler) CommitFile(ctx context.Context, req *media.CommitFileRequest) (*media.CommitFileResponse, error) {
	// 1. Basic validation
	if req.TmpPath == "" || req.PerPath == "" {
		return nil, status.Error(codes.InvalidArgument, "tmp_path and per_path must not be empty")
	}

	// 2. Security checkpoint (Cross-Validation)
	// Source file must be in the tmp/ directory
	if !strings.HasPrefix(req.TmpPath, "tmp/") {
		return nil, status.Error(codes.PermissionDenied, "Only files from staging area (tmp/) can be committed")
	}
	// Destination must be in the per/ directory
	if !strings.HasPrefix(req.PerPath, "per/") {
		return nil, status.Error(codes.PermissionDenied, "Destination must be in permanent storage area (per/)")
	}

	// check first 512 bytes of the file
	obj, err := h.storage.MinioClient.GetObject(ctx, h.storage.BucketName, req.TmpPath, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()

	buffer := make([]byte, 512)
	_, err = obj.Read(buffer)
	// (Ignore EOF error if file is smaller than 512 bytes)

	// 2. Use Golang standard library to inspect Magic Bytes
	realContentType := http.DetectContentType(buffer)

	// 3. Compare actual content type
	if !strings.HasPrefix(realContentType, "image/") {
		// FRAUD DETECTED!
		// Uploaded .exe file disguised as an image.

		// Immediately remove this malicious file from tmp
		h.storage.MinioClient.RemoveObject(ctx, h.storage.BucketName, req.TmpPath, minio.RemoveObjectOptions{})

		return nil, fmt.Errorf("detected forged file format,已被 system destroyed")
	}

	// 3. Call down to Storage layer
	err = h.storage.CommitFile(ctx, req.TmpPath, req.PerPath)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "system error while moving file: %v", err)
	}

	// 4. Return success result
	return &media.CommitFileResponse{
		Success: true,
	}, nil
}

func (h *MediaHandler) DeleteFile(ctx context.Context, req *media.DeleteFileRequest) (*media.DeleteFileResponse, error) {
	if strings.TrimSpace(req.ObjectPath) == "" {
		return nil, status.Error(codes.InvalidArgument, "object_path must not be empty")
	}

	err := h.storage.MinioClient.RemoveObject(ctx, h.storage.BucketName, req.ObjectPath, minio.RemoveObjectOptions{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "system error while deleting file: %v", err)
	}

	return &media.DeleteFileResponse{Success: true}, nil
}
