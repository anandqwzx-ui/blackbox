package models

import (
	"database/sql"
	"time"
)

// User represents the database model for the users table.
type User struct {
	ID           int64
	UUID         string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	CreatedBy    sql.NullInt64
	UpdatedAt    time.Time
	UpdatedBy    sql.NullInt64
	DeletedAt    sql.NullTime
	DeletedBy    sql.NullInt64
}
