package models

import (
	"database/sql"
	"time"
)

// Project represents a grouping of mocked endpoints belonging to a user and organisation.
type Project struct {
	ID            int64
	UUID          string
	OrganisationID int64
	OwnerUserID   int64
	Name          string
	Description   sql.NullString
	Code          string
	CreatedAt     time.Time
	CreatedBy     sql.NullInt64
	UpdatedAt     time.Time
	UpdatedBy     sql.NullInt64
	DeletedAt     sql.NullTime
	DeletedBy     sql.NullInt64
}
