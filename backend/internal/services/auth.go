package services

import (
	"errors"
	"time"

	"news-to-text/internal/models"
	"news-to-text/internal/repositories"
	"news-to-text/pkg/auth"
	"news-to-text/pkg/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(req *models.UserCreateRequest) (*models.UserResponse, string, error)
	Login(req *models.UserLoginRequest) (*models.UserResponse, string, error)
	Logout(token string) error
	ValidateToken(token string) (*auth.Claims, error)
	GetUserByID(id uint) (*models.UserResponse, error)
}

type authService struct {
	userRepo   repositories.UserRepository
	jwtManager *auth.JWTManager
	redis      *redis.Client
}

func NewAuthService(userRepo repositories.UserRepository, redisClient *redis.Client, jwtSecret string) AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtManager: auth.NewJWTManager(jwtSecret),
		redis:      redisClient,
	}
}

func (s *authService) Register(req *models.UserCreateRequest) (*models.UserResponse, string, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, "", err
	}
	if existingUser != nil {
		return nil, "", errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, "", err
	}

	// Create user
	user := &models.User{
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, "", err
	}

	// Generate JWT token
	token, err := s.jwtManager.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, "", err
	}

	return user.ToResponse(), token, nil
}

func (s *authService) Login(req *models.UserLoginRequest) (*models.UserResponse, string, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", errors.New("invalid credentials")
		}
		return nil, "", err
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := s.jwtManager.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, "", err
	}

	return user.ToResponse(), token, nil
}

func (s *authService) Logout(token string) error {
	// Add token to blacklist in Redis with expiration
	ctx := redis.Context()
	return s.redis.Set(ctx, "blacklist:"+token, "true", 24*time.Hour).Err()
}

func (s *authService) ValidateToken(token string) (*auth.Claims, error) {
	// Check if token is blacklisted
	ctx := redis.Context()
	isBlacklisted, err := s.redis.Exists(ctx, "blacklist:"+token).Result()
	if err != nil {
		return nil, err
	}
	if isBlacklisted > 0 {
		return nil, errors.New("token is blacklisted")
	}

	return s.jwtManager.ValidateToken(token)
}

func (s *authService) GetUserByID(id uint) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return user.ToResponse(), nil
}