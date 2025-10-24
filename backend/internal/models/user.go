package models

import (
    "database/sql"
    "time"
)

// User represents the database model for the users table.
type User struct {
    ID           int64         `db:"id"`
    UUID         string        `db:"uuid"`
    Email        string        `db:"email"`
    PasswordHash string        `db:"password_hash"`
    CreatedAt    time.Time     `db:"created_at"`
    CreatedBy    sql.NullInt64 `db:"created_by"`
    UpdatedAt    time.Time     `db:"updated_at"`
    UpdatedBy    sql.NullInt64 `db:"updated_by"`
    DeletedAt    sql.NullTime  `db:"deleted_at"`
    DeletedBy    sql.NullInt64 `db:"deleted_by"`
}
