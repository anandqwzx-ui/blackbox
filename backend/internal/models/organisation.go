package models

import (
    "database/sql"
    "time"
)

// Organisation represents an organisation a user owns and manages projects within.
type Organisation struct {
    ID          int64         `db:"id"`
    UUID        string        `db:"uuid"`
    Name        string        `db:"name"`
    OwnerUserID int64         `db:"owner_user_id"`
    CreatedAt   time.Time     `db:"created_at"`
    CreatedBy   sql.NullInt64 `db:"created_by"`
    UpdatedAt   time.Time     `db:"updated_at"`
    UpdatedBy   sql.NullInt64 `db:"updated_by"`
    DeletedAt   sql.NullTime  `db:"deleted_at"`
    DeletedBy   sql.NullInt64 `db:"deleted_by"`
}
