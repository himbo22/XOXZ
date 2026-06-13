package grpc

import (
	"context"
	"fmt"
	"time"

	media "github.com/himbo22/xoxz/common-service/protobuf/media"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type MediaClient interface {
	// Commit a valid file by moving it from tmp/ to per/.
	CommitFile(ctx context.Context, in *media.CommitFileRequest) (*media.CommitFileResponse, error)
	// Delete the physical file and clear cache.
	DeleteFile(ctx context.Context, in *media.DeleteFileRequest) (*media.DeleteFileResponse, error)

	Close() error
}

type mediaClient struct {
	client media.MediaServiceClient
	conn   *grpc.ClientConn
}

// Close implements [MediaClient].
func (m *mediaClient) Close() error {
	if m.conn != nil {
		return m.conn.Close()
	}
	return nil
}

// CommitFile implements [MediaClient].
func (m *mediaClient) CommitFile(ctx context.Context, in *media.CommitFileRequest) (*media.CommitFileResponse, error) {
	// Microservice calls must have a timeout so the profile flow does not hang with media-service.
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &media.CommitFileRequest{
		TmpPath: in.TmpPath,
		PerPath: in.PerPath,
	}

	// Call over the internal network.
	response, err := m.client.CommitFile(ctxTimeout, req)
	if err != nil {
		// Translate gRPC status codes into readable errors.
		st, ok := status.FromError(err)
		if ok {
			return nil, fmt.Errorf("media-service rejected the request (code: %s): %s", st.Code(), st.Message())
		}
		return nil, fmt.Errorf("could not connect to media-service: %w", err)
	}

	return response, nil
}

// DeleteFile implements [MediaClient].
func (m *mediaClient) DeleteFile(ctx context.Context, in *media.DeleteFileRequest) (*media.DeleteFileResponse, error) {
	req := &media.DeleteFileRequest{
		ObjectPath: in.ObjectPath,
	}

	response, err := m.client.DeleteFile(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			return nil, fmt.Errorf("media-service rejected file deletion (code: %s): %s", st.Code(), st.Message())
		}
		return nil, fmt.Errorf("could not connect to media-service to delete the file: %w", err)
	}

	return response, nil
}

func NewMediaClient(address string) (MediaClient, error) {
	// For production, use secure credentials (TLS). Here we use insecure for internal microservice traffic if within private net.
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("connect to grpc es-svc: %w", err)
	}

	client := media.NewMediaServiceClient(conn)

	return &mediaClient{
		client: client,
		conn:   conn,
	}, nil
}
