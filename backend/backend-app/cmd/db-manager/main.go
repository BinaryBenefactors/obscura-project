package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"obscura.app/backend/internal/config"
	"obscura.app/backend/internal/database"
	"obscura.app/backend/internal/domain/models"
	"obscura.app/backend/internal/domain/repository"
)

func main() {
	// Парсим флаги
	command := flag.String("cmd", "", "Command to execute")
	userID := flag.Uint("id", 0, "User ID for commands")
	email := flag.String("email", "", "User email")
	name := flag.String("name", "", "User name")
	password := flag.String("password", "", "User password")
	query := flag.String("query", "", "Search query")
	active := flag.Bool("active", true, "User active status")
	role := flag.String("role", "user", "User role")
	flag.Parse()

	// Загружаем конфигурацию
	cfg := config.NewConfig()

	// Устанавливаем переменные окружения для подключения к БД
	os.Setenv("DB_HOST", "db")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("DB_NAME", "appdb")

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
			fmt.Printf("ID: %d, Email: %s, Name: %s, Role: %s, Active: %v, Created: %s\n",
				user.ID, user.Email, user.Name, user.Role, user.Active, user.CreatedAt.Format(time.RFC3339))
		}

	case "search-users":
		if *query == "" {
			log.Fatal("Search query is required")
		}
		users, err := userRepo.List(ctx, 0, 100)
		if err != nil {
			log.Fatalf("Failed to list users: %v", err)
		}
		fmt.Printf("Search results for '%s':\n", *query)
		found := false
		for _, user := range users {
			if strings.Contains(strings.ToLower(user.Email), strings.ToLower(*query)) ||
				strings.Contains(strings.ToLower(user.Name), strings.ToLower(*query)) {
				fmt.Printf("ID: %d, Email: %s, Name: %s, Role: %s, Active: %v\n",
					user.ID, user.Email, user.Name, user.Role, user.Active)
				found = true
			}
		}
		if !found {
			fmt.Println("No users found")
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
		fmt.Printf("User details:\nID: %d\nEmail: %s\nName: %s\nRole: %s\nActive: %v\nCreated: %s\n",
			user.ID, user.Email, user.Name, user.Role, user.Active, user.CreatedAt.Format(time.RFC3339))

	case "create-user":
		if *email == "" || *password == "" {
			log.Fatal("Email and password are required")
		}
		user := &models.User{
			Email:    *email,
			Name:     *name,
			Password: *password,
			Role:     "user",
			Active:   true,
		}
		if err := userRepo.Create(ctx, user); err != nil {
			log.Fatalf("Failed to create user: %v", err)
		}
		fmt.Printf("User created successfully:\nID: %d\nEmail: %s\nName: %s\nRole: %s\n",
			user.ID, user.Email, user.Name, user.Role)

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
		fmt.Printf("User updated successfully:\nID: %d\nEmail: %s\nName: %s\nRole: %s\n",
			user.ID, user.Email, user.Name, user.Role)

	case "toggle-user":
		if *userID == 0 {
			log.Fatal("User ID is required")
		}
		user, err := userRepo.GetByID(ctx, *userID)
		if err != nil {
			log.Fatalf("Failed to get user: %v", err)
		}
		user.Active = *active
		if err := userRepo.Update(ctx, user); err != nil {
			log.Fatalf("Failed to update user: %v", err)
		}
		fmt.Printf("User %d %s\n", *userID, map[bool]string{true: "activated", false: "deactivated"}[*active])

	case "change-role":
		if *userID == 0 {
			log.Fatal("User ID is required")
		}
		user, err := userRepo.GetByID(ctx, *userID)
		if err != nil {
			log.Fatalf("Failed to get user: %v", err)
		}
		user.Role = *role
		if err := userRepo.Update(ctx, user); err != nil {
			log.Fatalf("Failed to update user: %v", err)
		}
		fmt.Printf("User %d role changed to %s\n", *userID, *role)

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

	case "delete-session":
		if *userID == 0 {
			log.Fatal("Session ID is required")
		}
		if err := sessionRepo.Delete(ctx, *userID); err != nil {
			log.Fatalf("Failed to delete session: %v", err)
		}
		fmt.Printf("Session %d deleted successfully\n", *userID)

	case "delete-user-sessions":
		if *userID == 0 {
			log.Fatal("User ID is required")
		}
		sessions, err := sessionRepo.List(ctx, 0, 100)
		if err != nil {
			log.Fatalf("Failed to list sessions: %v", err)
		}
		deleted := 0
		for _, session := range sessions {
			if session.UserID == *userID {
				if err := sessionRepo.Delete(ctx, session.ID); err != nil {
					log.Printf("Failed to delete session %d: %v", session.ID, err)
					continue
				}
				deleted++
			}
		}
		fmt.Printf("Deleted %d sessions for user %d\n", deleted, *userID)

	case "cleanup-sessions":
		sessions, err := sessionRepo.List(ctx, 0, 100)
		if err != nil {
			log.Fatalf("Failed to list sessions: %v", err)
		}
		deleted := 0
		now := time.Now()
		for _, session := range sessions {
			if session.ExpiresAt.Before(now) {
				if err := sessionRepo.Delete(ctx, session.ID); err != nil {
					log.Printf("Failed to delete session %d: %v", session.ID, err)
					continue
				}
				deleted++
			}
		}
		fmt.Printf("Cleaned up %d expired sessions\n", deleted)

	case "show-stats":
		users, err := userRepo.List(ctx, 0, 100)
		if err != nil {
			log.Fatalf("Failed to list users: %v", err)
		}
		sessions, err := sessionRepo.List(ctx, 0, 100)
		if err != nil {
			log.Fatalf("Failed to list sessions: %v", err)
		}
		activeUsers := 0
		for _, user := range users {
			if user.Active {
				activeUsers++
			}
		}
		fmt.Printf("Statistics:\n")
		fmt.Printf("Total users: %d\n", len(users))
		fmt.Printf("Active users: %d\n", activeUsers)
		fmt.Printf("Total sessions: %d\n", len(sessions))
		fmt.Printf("Active sessions: %d\n", len(sessions))

	default:
		fmt.Println("Available commands:")
		fmt.Println("  list-users - List all users")
		fmt.Println("  search-users -query <query> - Search users by name or email")
		fmt.Println("  get-user -id <user_id> or -email <email> - Get user details")
		fmt.Println("  create-user -email <email> -name <name> -password <password> - Create new user")
		fmt.Println("  delete-user -id <user_id> - Delete user")
		fmt.Println("  update-user -id <user_id> [-name <name>] [-password <password>] - Update user")
		fmt.Println("  toggle-user -id <user_id> -active <true|false> - Activate/deactivate user")
		fmt.Println("  change-role -id <user_id> -role <role> - Change user role")
		fmt.Println("  list-sessions - List all active sessions")
		fmt.Println("  delete-session -id <session_id> - Delete session")
		fmt.Println("  delete-user-sessions -user-id <user_id> - Delete all user sessions")
		fmt.Println("  cleanup-sessions - Clean up expired sessions")
		fmt.Println("  show-stats - Show database statistics")
		os.Exit(1)
	}
}
