package models

import (
	"database/sql"
	"time"
)

// Endpoint represents a mockable API endpoint associated with a project.
type Endpoint struct {
	ID             int64
	UUID           string
	ProjectID      int64
	Name           sql.NullString
	Method         string
	Path           string
	ResponseStatus int
	ResponseBody   string
	ResponseHeaders string
	Enabled        bool
	CreatedAt      time.Time
	CreatedBy      sql.NullInt64
	UpdatedAt      time.Time
	UpdatedBy      sql.NullInt64
	DeletedAt      sql.NullTime
	DeletedBy      sql.NullInt64
}
