package grpc

import (
	"fmt"
	"net"

	media "github.com/himbo22/xoxz/common-service/protobuf/media"
	"github.com/himbo22/xoxz/media-service/internal/adapter/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// MediaServer executes MediaService logic
type MediaServer struct {
	server  *grpc.Server
	storage *storage.MinioClient
}

func NewMediaServer(storage *storage.MinioClient) *MediaServer {
	grpcServer := grpc.NewServer()

	// register handlers
	handler := NewMediaHandler(storage)
	media.RegisterMediaServiceServer(grpcServer, handler)

	// Enable reflection for debugging (e.g., using grpcurl)
	reflection.Register(grpcServer)

	return &MediaServer{
		storage: storage,
		server:  grpcServer,
	}
}

func (m *MediaServer) Start(port string) error {
	// open tcp port
	listenAddr := fmt.Sprintf(":%s", port)
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", listenAddr, err)
	}

	return m.server.Serve(listener)
}

// Stop gracefully shuts down the server
func (m *MediaServer) Stop() {
	if m.server != nil {
		m.server.GracefulStop()
	}
}
