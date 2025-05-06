package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"obscura.app/backend/internal/domain/models"
	"obscura.app/backend/internal/domain/repository"
)

type FileService struct {
	repo       repository.FileRepository
	uploadPath string
	maxSize    int64
}

func NewFileService(repo repository.FileRepository, uploadPath string, maxSize int64) *FileService {
	return &FileService{
		repo:       repo,
		uploadPath: uploadPath,
		maxSize:    maxSize,
	}
}

func (s *FileService) UploadFile(ctx context.Context, file *multipart.FileHeader, userID string) (*models.UploadedFile, error) {
	if file.Size > s.maxSize {
		return nil, fmt.Errorf("file too large (max %dMB)", s.maxSize/(1024*1024))
	}

	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Создаем директорию если не существует
	if err := os.MkdirAll(s.uploadPath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Генерируем уникальное имя файла
	fileID := fmt.Sprintf("%d", time.Now().UnixNano())
	newFileName := fmt.Sprintf("%s%s", fileID, filepath.Ext(file.Filename))
	filePath := filepath.Join(s.uploadPath, newFileName)

	// Создаем файл
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// Копируем содержимое
	if _, err = io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Создаем запись в репозитории
	uploadedFile := &models.UploadedFile{
		ID:         fileID,
		Original:   newFileName,
		Processed:  fmt.Sprintf("/processed/%s_result.mp4", newFileName),
		UploadedAt: time.Now(),
		Status:     string(models.StatusProcessing),
		UserID:     userID,
	}

	if err := s.repo.SaveFile(ctx, uploadedFile); err != nil {
		// Если не удалось сохранить в репозиторий, удаляем файл
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save file info: %w", err)
	}

	return uploadedFile, nil
}

func (s *FileService) GetFileStatus(ctx context.Context, fileID string) (*models.UploadedFile, error) {
	return s.repo.GetFileByID(ctx, fileID)
}

func (s *FileService) GetUserFiles(ctx context.Context, userID string) ([]models.UploadedFile, error) {
	return s.repo.GetUserFiles(ctx, userID)
}

func (s *FileService) UpdateFileStatus(ctx context.Context, fileID string, status string) error {
	return s.repo.UpdateFileStatus(ctx, fileID, status)
}
