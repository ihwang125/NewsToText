package services

import (
	"testing"

	"news-to-text/internal/models"
	"news-to-text/internal/repositories"
	"news-to-text/pkg/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.User{}, &models.Alert{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func setupTestRedis() *redis.Client {
	// Use Redis mock or in-memory alternative for testing
	// For simplicity, we'll use a real Redis instance running on default port
	// In production tests, you'd want to use a test-specific Redis instance
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   1, // Use a different DB for tests
	})

	// Clear test database
	client.FlushDB(redis.Context())

	return client
}

func TestAuthService_Register(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	redisClient := setupTestRedis()
	userRepo := repositories.NewUserRepository(db)
	authService := NewAuthService(userRepo, redisClient, "test-secret")

	tests := []struct {
		name    string
		request *models.UserCreateRequest
		wantErr bool
	}{
		{
			name: "Valid registration",
			request: &models.UserCreateRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "Duplicate email",
			request: &models.UserCreateRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "Invalid email",
			request: &models.UserCreateRequest{
				Email:    "invalid-email",
				Password: "password123",
			},
			wantErr: false, // Service doesn't validate email format, that's done at handler level
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, token, err := authService.Register(tt.request)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if user == nil {
				t.Errorf("Expected user but got nil")
				return
			}

			if token == "" {
				t.Errorf("Expected token but got empty string")
				return
			}

			if user.Email != tt.request.Email {
				t.Errorf("Expected email %s but got %s", tt.request.Email, user.Email)
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	redisClient := setupTestRedis()
	userRepo := repositories.NewUserRepository(db)
	authService := NewAuthService(userRepo, redisClient, "test-secret")

	// Create a test user
	hashedPassword, _ := utils.HashPassword("password123")
	testUser := &models.User{
		Email:    "test@example.com",
		Password: hashedPassword,
	}
	userRepo.Create(testUser)

	tests := []struct {
		name    string
		request *models.UserLoginRequest
		wantErr bool
	}{
		{
			name: "Valid login",
			request: &models.UserLoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "Invalid password",
			request: &models.UserLoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			wantErr: true,
		},
		{
			name: "Non-existent user",
			request: &models.UserLoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, token, err := authService.Login(tt.request)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if user == nil {
				t.Errorf("Expected user but got nil")
				return
			}

			if token == "" {
				t.Errorf("Expected token but got empty string")
				return
			}

			if user.Email != tt.request.Email {
				t.Errorf("Expected email %s but got %s", tt.request.Email, user.Email)
			}
		})
	}
}