package store

import "obscura.app/backend/internal/model"

// UserRepository ...
type UserRepository interface {
	Create(*model.User) error
	FindByID(int) (*model.User, error)
	FindByEmail(string) (*model.User, error)
}
