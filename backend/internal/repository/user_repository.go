package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ctonew/mockapi/internal/models"
)

// ErrUserNotFound is returned when the requested user cannot be located.
var ErrUserNotFound = errors.New("user not found")

// UserRepository provides access to user persistence operations.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository constructs a new UserRepository instance.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a new user record.
func (r *UserRepository) Create(ctx context.Context, email, passwordHash string) (*models.User, error) {
	query := `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, uuid, email, password_hash, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
	`

	row := r.db.QueryRowContext(ctx, query, email, passwordHash)
	user := &models.User{}
	if err := scanUser(row, user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetByEmail returns a user by email if not soft deleted.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, uuid, email, password_hash, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`

	row := r.db.QueryRowContext(ctx, query, email)
	user := &models.User{}
	if err := scanUser(row, user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

// GetByUUID returns a user by UUID if not soft deleted.
func (r *UserRepository) GetByUUID(ctx context.Context, uuid string) (*models.User, error) {
	query := `
		SELECT id, uuid, email, password_hash, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
		FROM users
		WHERE uuid = $1 AND deleted_at IS NULL
	`

	row := r.db.QueryRowContext(ctx, query, uuid)
	user := &models.User{}
	if err := scanUser(row, user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func scanUser(row *sql.Row, user *models.User) error {
	return row.Scan(
		&user.ID,
		&user.UUID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.CreatedBy,
		&user.UpdatedAt,
		&user.UpdatedBy,
		&user.DeletedAt,
		&user.DeletedBy,
	)
}
