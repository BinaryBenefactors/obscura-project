package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"obscura.app/backend/internal/config"
	"obscura.app/backend/internal/database"
	"obscura.app/backend/internal/domain/models"
	"obscura.app/backend/internal/domain/repository"
)

func main() {
	// Парсим флаги
	command := flag.String("cmd", "", "Command to execute (list-users, get-user, create-user, delete-user, update-user, list-sessions)")
	userID := flag.Uint("id", 0, "User ID for commands")
	email := flag.String("email", "", "User email")
	name := flag.String("name", "", "User name")
	password := flag.String("password", "", "User password")
	env := flag.String("env", "test", "Environment to use (test or prod)")
	flag.Parse()

	// Загружаем конфигурацию
	cfg := config.NewConfig()

	// Устанавливаем переменные окружения для подключения к БД
	os.Setenv("DB_HOST", "db")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "postgres")

	// Выбираем базу данных в зависимости от окружения
	dbName := "obscura_test"
	if *env == "prod" {
		dbName = "obscura"
		log.Println("WARNING: Using production database!")
	}
	os.Setenv("DB_NAME", dbName)

	// Подключаемся к БД
	db, err := database.NewPostgresDB(
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
	)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Создаем репозитории
	userRepo := repository.NewGormUserRepository(db)
	sessionRepo := repository.NewGormSessionRepository(db)
	ctx := context.Background()

	// Выполняем команду
	switch *command {
	case "list-users":
		users, err := userRepo.List(ctx, 0, 100)
		if err != nil {
			log.Fatalf("Failed to list users: %v", err)
		}
		fmt.Println("Users:")
		for _, user := range users {
			fmt.Printf("ID: %d, Email: %s, Name: %s, Created: %s\n",
				user.ID, user.Email, user.Name, user.CreatedAt.Format(time.RFC3339))
		}

	case "get-user":
		if *userID == 0 && *email == "" {
			log.Fatal("Either user ID or email is required")
		}
		var user *models.User
		var err error
		if *userID > 0 {
			user, err = userRepo.GetByID(ctx, *userID)
		} else {
			user, err = userRepo.GetByEmail(ctx, *email)
		}
		if err != nil {
			log.Fatalf("Failed to get user: %v", err)
		}
		fmt.Printf("User details:\nID: %d\nEmail: %s\nName: %s\nCreated: %s\n",
			user.ID, user.Email, user.Name, user.CreatedAt.Format(time.RFC3339))

	case "create-user":
		if *email == "" || *password == "" {
			log.Fatal("Email and password are required")
		}
		user := &models.User{
			Email:    *email,
			Name:     *name,
			Password: *password,
		}
		if err := userRepo.Create(ctx, user); err != nil {
			log.Fatalf("Failed to create user: %v", err)
		}
		fmt.Printf("User created successfully:\nID: %d\nEmail: %s\nName: %s\n",
			user.ID, user.Email, user.Name)

	case "delete-user":
		if *userID == 0 {
			log.Fatal("User ID is required")
		}
		if err := userRepo.Delete(ctx, *userID); err != nil {
			log.Fatalf("Failed to delete user: %v", err)
		}
		fmt.Printf("User %d deleted successfully\n", *userID)

	case "update-user":
		if *userID == 0 {
			log.Fatal("User ID is required")
		}
		user, err := userRepo.GetByID(ctx, *userID)
		if err != nil {
			log.Fatalf("Failed to get user: %v", err)
		}
		if *name != "" {
			user.Name = *name
		}
		if *password != "" {
			user.Password = *password
		}
		if err := userRepo.Update(ctx, user); err != nil {
			log.Fatalf("Failed to update user: %v", err)
		}
		fmt.Printf("User updated successfully:\nID: %d\nEmail: %s\nName: %s\n",
			user.ID, user.Email, user.Name)

	case "list-sessions":
		sessions, err := sessionRepo.List(ctx, 0, 100)
		if err != nil {
			log.Fatalf("Failed to list sessions: %v", err)
		}
		fmt.Println("Active sessions:")
		for _, session := range sessions {
			fmt.Printf("ID: %d, User ID: %d, Created: %s, Expires: %s\n",
				session.ID, session.UserID,
				session.CreatedAt.Format(time.RFC3339),
				session.ExpiresAt.Format(time.RFC3339))
		}

	default:
		fmt.Println("Available commands:")
		fmt.Println("  list-users - List all users")
		fmt.Println("  get-user -id <user_id> or -email <email> - Get user details")
		fmt.Println("  create-user -email <email> -name <name> -password <password> - Create new user")
		fmt.Println("  delete-user -id <user_id> - Delete user")
		fmt.Println("  update-user -id <user_id> [-name <name>] [-password <password>] - Update user")
		fmt.Println("  list-sessions - List all active sessions")
		fmt.Println("\nEnvironment options:")
		fmt.Println("  -env test  - Use test database (default)")
		fmt.Println("  -env prod  - Use production database")
		os.Exit(1)
	}
}
