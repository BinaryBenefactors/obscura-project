package repository

import (
	"context"

	"obscura.app/backend/internal/domain/models"
)

type FileRepository interface {
	// Сохранение информации о загруженном файле
	SaveFile(ctx context.Context, file *models.UploadedFile) error

	// Получение информации о файле по ID
	GetFileByID(ctx context.Context, fileID string) (*models.UploadedFile, error)

	// Получение списка файлов пользователя
	GetUserFiles(ctx context.Context, userID string) ([]models.UploadedFile, error)

	// Обновление статуса файла
	UpdateFileStatus(ctx context.Context, fileID string, status string) error

	// Удаление файла
	DeleteFile(ctx context.Context, fileID string) error
}
