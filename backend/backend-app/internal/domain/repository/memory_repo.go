package repository

import (
	"context"
	"sync"

	"obscura.app/backend/internal/domain/models"
)

// MemoryFileRepository - временная реализация репозитория в памяти
type MemoryFileRepository struct {
	files map[string][]models.UploadedFile
	mu    sync.RWMutex
}

// NewMemoryFileRepository создает новый репозиторий в памяти
func NewMemoryFileRepository() *MemoryFileRepository {
	return &MemoryFileRepository{
		files: make(map[string][]models.UploadedFile),
	}
}

func (r *MemoryFileRepository) SaveFile(ctx context.Context, file *models.UploadedFile) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.files[file.UserID] = append(r.files[file.UserID], *file)
	return nil
}

func (r *MemoryFileRepository) GetFileByID(ctx context.Context, fileID string) (*models.UploadedFile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, files := range r.files {
		for _, file := range files {
			if file.ID == fileID {
				return &file, nil
			}
		}
	}
	return nil, nil
}

func (r *MemoryFileRepository) GetUserFiles(ctx context.Context, userID string) ([]models.UploadedFile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.files[userID], nil
}

func (r *MemoryFileRepository) UpdateFileStatus(ctx context.Context, fileID string, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for userID, files := range r.files {
		for i, file := range files {
			if file.ID == fileID {
				r.files[userID][i].Status = status
				return nil
			}
		}
	}
	return nil
}

func (r *MemoryFileRepository) DeleteFile(ctx context.Context, fileID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for userID, files := range r.files {
		for i, file := range files {
			if file.ID == fileID {
				r.files[userID] = append(files[:i], files[i+1:]...)
				return nil
			}
		}
	}
	return nil
}
