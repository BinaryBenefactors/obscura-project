package models

import "time"

type UploadedFile struct {
	ID         string    `json:"id"`
	Original   string    `json:"original_name"`
	Processed  string    `json:"processed_path"`
	UploadedAt time.Time `json:"uploaded_at"`
	Status     string    `json:"status"`
	UserID     string    `json:"user_id"`
}

type FileStatus string

const (
	StatusProcessing FileStatus = "processing"
	StatusCompleted  FileStatus = "completed"
	StatusFailed     FileStatus = "failed"
)

type FileType string

const (
	TypeImage FileType = "image"
	TypeVideo FileType = "video"
)

type FileMetadata struct {
	Type     FileType
	MimeType string
	Size     int64
	Path     string
}
