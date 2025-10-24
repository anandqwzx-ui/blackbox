package models

import (
	"database/sql"
	"time"
)

// Organisation represents an organisation a user owns and manages projects within.
type Organisation struct {
	ID           int64
	UUID         string
	Name         string
	OwnerUserID  int64
	CreatedAt    time.Time
	CreatedBy    sql.NullInt64
	UpdatedAt    time.Time
	UpdatedBy    sql.NullInt64
	DeletedAt    sql.NullTime
	DeletedBy    sql.NullInt64
}
