package model

import "github.com/google/uuid"

type GenerateURLRequest struct {
	UserID      uuid.UUID `json:"user_id"`
	ContentType string    `json:"type"`      // image/jpeg , image/png , ...
	FileName    string    `json:"file_name"` // abc.jpeg
}

type GenerateURLResponse struct {
	UploadURL string // Direct link to Nginx/MinIO (For FE to call PUT)
	TmpPath   string // Temporary path (For FE to store, then send to Backend for confirmation)
	PerPath   string // Permanent path (Pre-generated for Frontend reference)
}
