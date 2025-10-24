package models

import (
    "database/sql"
    "time"
)

// Endpoint represents a mockable API endpoint associated with a project.
type Endpoint struct {
    ID              int64          `db:"id"`
    UUID            string         `db:"uuid"`
    ProjectID       int64          `db:"project_id"`
    Name            sql.NullString `db:"name"`
    Method          string         `db:"method"`
    Path            string         `db:"path"`
    ResponseStatus  int            `db:"response_status"`
    ResponseBody    string         `db:"response_body"`
    ResponseHeaders string         `db:"response_headers"`
    Enabled         bool           `db:"enabled"`
    CreatedAt       time.Time      `db:"created_at"`
    CreatedBy       sql.NullInt64  `db:"created_by"`
    UpdatedAt       time.Time      `db:"updated_at"`
    UpdatedBy       sql.NullInt64  `db:"updated_by"`
    DeletedAt       sql.NullTime   `db:"deleted_at"`
    DeletedBy       sql.NullInt64  `db:"deleted_by"`
}
