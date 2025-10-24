package service

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/ctonew/mockapi/internal/models"
	"github.com/ctonew/mockapi/internal/repository"
	"github.com/ctonew/mockapi/internal/utils"
)

// AuthServiceError provides sentinel errors the handler layer can reason about.
var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// AuthService orchestrates user registration and authentication flows.
type AuthService struct {
	users           *repository.UserRepository
	passwordManager *utils.PasswordManager
	tokenManager    *utils.TokenManager
}

// NewAuthService constructs a new AuthService.
func NewAuthService(users *repository.UserRepository, passwordManager *utils.PasswordManager, tokenManager *utils.TokenManager) *AuthService {
	return &AuthService{
		users:           users,
		passwordManager: passwordManager,
		tokenManager:    tokenManager,
	}
}

// Signup registers a new user and returns the created model alongside a JWT token.
func (s *AuthService) Signup(ctx context.Context, email, password string) (*models.User, string, error) {
	hash, err := s.passwordManager.HashPassword(password)
	if err != nil {
		return nil, "", err
	}

	user, err := s.users.Create(ctx, email, hash)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, "", ErrEmailAlreadyExists
		}
		return nil, "", err
	}

	token, err := s.tokenManager.GenerateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// Login authenticates a user and returns their model with a signed JWT token.
func (s *AuthService) Login(ctx context.Context, email, password string) (*models.User, string, error) {
	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, "", ErrInvalidCredentials
		}
		return nil, "", err
	}

	if err := s.passwordManager.Compare(user.PasswordHash, password); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	token, err := s.tokenManager.GenerateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}
