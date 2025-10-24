package repository

import (
    "context"
    "database/sql"
    "errors"

    "github.com/jmoiron/sqlx"

    "github.com/crudbox/crudbox/internal/models"
)

// ErrOrganisationNotFound indicates the requested organisation does not exist or is inaccessible.
var ErrOrganisationNotFound = errors.New("organisation not found")

// OrganisationRepository manages persistence for organisations.
type OrganisationRepository struct {
    db *sqlx.DB
}

// NewOrganisationRepository creates an OrganisationRepository.
func NewOrganisationRepository(db *sqlx.DB) *OrganisationRepository {
    return &OrganisationRepository{db: db}
}

// Create inserts a new organisation record.
func (r *OrganisationRepository) Create(ctx context.Context, name string, ownerUserID, actorID int64) (*models.Organisation, error) {
    query := `
        INSERT INTO organisations (name, owner_user_id, created_by, updated_by)
        VALUES ($1, $2, $3, $3)
        RETURNING id, uuid, name, owner_user_id, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
    `

    row := r.db.QueryRowxContext(ctx, query, name, ownerUserID, actorID)
    organisation := &models.Organisation{}
    if err := row.StructScan(organisation); err != nil {
        return nil, err
    }

    return organisation, nil
}

// ListByUser returns the active organisations owned by the supplied user.
func (r *OrganisationRepository) ListByUser(ctx context.Context, ownerUserID int64) ([]models.Organisation, error) {
    query := `
        SELECT id, uuid, name, owner_user_id, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
        FROM organisations
        WHERE owner_user_id = $1 AND deleted_at IS NULL
        ORDER BY created_at DESC
    `

    var organisations []models.Organisation
    if err := r.db.SelectContext(ctx, &organisations, query, ownerUserID); err != nil {
        return nil, err
    }

    return organisations, nil
}

// GetByUUIDAndOwner ensures an organisation belongs to the provided user before returning it.
func (r *OrganisationRepository) GetByUUIDAndOwner(ctx context.Context, organisationUUID string, ownerUserID int64) (*models.Organisation, error) {
    query := `
        SELECT id, uuid, name, owner_user_id, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
        FROM organisations
        WHERE uuid = $1 AND owner_user_id = $2 AND deleted_at IS NULL
    `

    organisation := &models.Organisation{}
    if err := r.db.GetContext(ctx, organisation, query, organisationUUID, ownerUserID); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrOrganisationNotFound
        }
        return nil, err
    }

    return organisation, nil
}
