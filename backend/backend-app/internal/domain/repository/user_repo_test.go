package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"obscura.app/backend/internal/domain/models"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// Подключаемся к тестовой базе данных
	dsn := "host=db port=5432 user=postgres password=postgres dbname=obscura_test sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)

	// Очищаем таблицы перед тестом
	db.Migrator().DropTable(&models.User{}, &models.Session{})
	db.AutoMigrate(&models.User{}, &models.Session{})

	return db
}

func TestUserRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormUserRepository(db)
	ctx := context.Background()

	// Тест создания пользователя
	t.Run("Create User", func(t *testing.T) {
		user := &models.User{
			Email:    "test@example.com",
			Password: "password123",
			Name:     "Test User",
		}

		err := repo.Create(ctx, user)
		assert.NoError(t, err)
		assert.NotZero(t, user.ID)
	})

	// Тест получения пользователя по email
	t.Run("Get User by Email", func(t *testing.T) {
		user, err := repo.GetByEmail(ctx, "test@example.com")
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "test@example.com", user.Email)
		assert.Equal(t, "Test User", user.Name)
		assert.True(t, user.CheckPassword("password123"))
	})

	// Тест обновления пользователя
	t.Run("Update User", func(t *testing.T) {
		user, err := repo.GetByEmail(ctx, "test@example.com")
		assert.NoError(t, err)

		user.Name = "Updated Name"
		err = repo.Update(ctx, user)
		assert.NoError(t, err)

		updatedUser, err := repo.GetByEmail(ctx, "test@example.com")
		assert.NoError(t, err)
		assert.Equal(t, "Updated Name", updatedUser.Name)
	})

	// Тест получения списка пользователей
	t.Run("List Users", func(t *testing.T) {
		// Создаем еще одного пользователя
		user2 := &models.User{
			Email:    "test2@example.com",
			Password: "password123",
			Name:     "Test User 2",
		}
		err := repo.Create(ctx, user2)
		assert.NoError(t, err)

		users, err := repo.List(ctx, 0, 10)
		assert.NoError(t, err)
		assert.Len(t, users, 2)
	})

	// Тест удаления пользователя
	t.Run("Delete User", func(t *testing.T) {
		user, err := repo.GetByEmail(ctx, "test@example.com")
		assert.NoError(t, err)

		err = repo.Delete(ctx, user.ID)
		assert.NoError(t, err)

		_, err = repo.GetByEmail(ctx, "test@example.com")
		assert.Error(t, err)
	})
}
