package models

import (
    "database/sql"
    "time"
)

// Project represents a grouping of mocked endpoints belonging to a user and organisation.
type Project struct {
    ID             int64          `db:"id"`
    UUID           string         `db:"uuid"`
    OrganisationID int64          `db:"organisation_id"`
    OwnerUserID    int64          `db:"owner_user_id"`
    Name           string         `db:"name"`
    Description    sql.NullString `db:"description"`
    Code           string         `db:"code"`
    CreatedAt      time.Time      `db:"created_at"`
    CreatedBy      sql.NullInt64  `db:"created_by"`
    UpdatedAt      time.Time      `db:"updated_at"`
    UpdatedBy      sql.NullInt64  `db:"updated_by"`
    DeletedAt      sql.NullTime   `db:"deleted_at"`
    DeletedBy      sql.NullInt64  `db:"deleted_by"`
}
