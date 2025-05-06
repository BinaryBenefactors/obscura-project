package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"obscura.app/backend/internal/domain/models"
	"obscura.app/backend/internal/handlers/middleware"
	"obscura.app/backend/internal/services"
)

type FileHandler struct {
	fileService *services.FileService
}

func NewFileHandler(fileService *services.FileService) *FileHandler {
	return &FileHandler{
		fileService: fileService,
	}
}

func (h *FileHandler) UploadFile(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return middleware.SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	userID, _ := r.Context().Value(middleware.UserIDKey).(string)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return middleware.SendError(w, "Failed to parse form data", http.StatusBadRequest)
	}

	file, fileHeader, err := r.FormFile("mediaFile")
	if err != nil {
		return middleware.SendError(w, "No file uploaded", http.StatusBadRequest)
	}
	defer file.Close()

	uploadedFile, err := h.fileService.UploadFile(r.Context(), fileHeader, userID)
	if err != nil {
		return err
	}

	response := models.UploadResponse{
		SuccessResponse: models.SuccessResponse{
			FileURL:  "/files/" + uploadedFile.Original,
			FileType: fileHeader.Header.Get("Content-Type"),
			Message:  "File is being processed",
		},
		FileID: uploadedFile.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(response)
}

func (h *FileHandler) GetFileStatus(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return middleware.SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	fileID := strings.TrimPrefix(r.URL.Path, "/api/upload/")
	if fileID == "" {
		return middleware.SendError(w, "File ID required", http.StatusBadRequest)
	}

	file, err := h.fileService.GetFileStatus(r.Context(), fileID)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(file)
}

func (h *FileHandler) GetUserFiles(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return middleware.SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	userID, _ := r.Context().Value(middleware.UserIDKey).(string)
	if strings.HasPrefix(userID, "anon_") {
		return middleware.SendError(w, "Authentication required", http.StatusUnauthorized)
	}

	files, err := h.fileService.GetUserFiles(r.Context(), userID)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(files)
}
