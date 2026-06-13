package model

type StreamStatus string

const (
	StatusPending StreamStatus = "pending"
	StatusActive  StreamStatus = "active"
	StatusEnded   StreamStatus = "ended"
	StatusFailed  StreamStatus = "failed"
)

type StreamSource string

const (
	SourceDirect  StreamSource = "direct"  // SDK (PC/Mobile webcam)
	SourceIngress StreamSource = "ingress" // OBS/Encoder
)

type CreateStreamRequest struct {
	UserID   string       `json:"user_id"`
	RoomName string       `json:"room_name"`
	Protocol string       `json:"protocol"` // rtmp | whip
	Source   StreamSource `json:"source"`
}

type StreamResponse struct {
	StreamID string       `json:"stream_id"` // Record ID in Database
	RoomName string       `json:"room_name"` // Identifier on LiveKit
	Status   StreamStatus `json:"status"`    // PENDING status
	WSURL    string       `json:"ws_url"`    // WebSocket URL cho LiveKit

	// Granted to Mobile/Webcam (Source: Direct)
	PublisherToken string `json:"publisher_token,omitempty"`

	// Granted to OBS (Source: Ingress)
	IngressURL string `json:"ingress_url,omitempty"`
	StreamKey  string `json:"stream_key,omitempty"`
}

type IngressResult struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	StreamKey string `json:"stream_key"`
}
