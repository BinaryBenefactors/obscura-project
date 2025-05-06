package models

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

type SuccessResponse struct {
	FileURL  string `json:"fileUrl"`
	FileType string `json:"fileType"`
	Message  string `json:"message"`
}

type UploadResponse struct {
	SuccessResponse
	FileID string `json:"file_id"`
}

type FileStatusResponse struct {
	ID         string `json:"id"`
	Status     string `json:"status"`
	Processed  string `json:"processed_path,omitempty"`
	Original   string `json:"original_name"`
	UploadedAt string `json:"uploaded_at"`
}
