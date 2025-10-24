package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ctonew/mockapi/internal/models"
)

// ErrOrganisationNotFound indicates the requested organisation does not exist or is inaccessible.
var ErrOrganisationNotFound = errors.New("organisation not found")

// OrganisationRepository manages persistence for organisations.
type OrganisationRepository struct {
	db *sql.DB
}

// NewOrganisationRepository creates an OrganisationRepository.
func NewOrganisationRepository(db *sql.DB) *OrganisationRepository {
	return &OrganisationRepository{db: db}
}

// Create inserts a new organisation record.
func (r *OrganisationRepository) Create(ctx context.Context, name string, ownerUserID, actorID int64) (*models.Organisation, error) {
	query := `
		INSERT INTO organisations (name, owner_user_id, created_by, updated_by)
		VALUES ($1, $2, $3, $3)
		RETURNING id, uuid, name, owner_user_id, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
	`

	row := r.db.QueryRowContext(ctx, query, name, ownerUserID, actorID)
	organisation := &models.Organisation{}
	if err := scanOrganisation(row, organisation); err != nil {
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

	rows, err := r.db.QueryContext(ctx, query, ownerUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var organisations []models.Organisation
	for rows.Next() {
		var organisation models.Organisation
		if err := rows.Scan(
			&organisation.ID,
			&organisation.UUID,
			&organisation.Name,
			&organisation.OwnerUserID,
			&organisation.CreatedAt,
			&organisation.CreatedBy,
			&organisation.UpdatedAt,
			&organisation.UpdatedBy,
			&organisation.DeletedAt,
			&organisation.DeletedBy,
		); err != nil {
			return nil, err
		}
		organisations = append(organisations, organisation)
	}

	return organisations, rows.Err()
}

// GetByUUIDAndOwner ensures an organisation belongs to the provided user before returning it.
func (r *OrganisationRepository) GetByUUIDAndOwner(ctx context.Context, organisationUUID string, ownerUserID int64) (*models.Organisation, error) {
	query := `
		SELECT id, uuid, name, owner_user_id, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
		FROM organisations
		WHERE uuid = $1 AND owner_user_id = $2 AND deleted_at IS NULL
	`

	row := r.db.QueryRowContext(ctx, query, organisationUUID, ownerUserID)
	organisation := &models.Organisation{}
	if err := scanOrganisation(row, organisation); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrOrganisationNotFound
		}
		return nil, err
	}

	return organisation, nil
}

func scanOrganisation(row *sql.Row, organisation *models.Organisation) error {
	return row.Scan(
		&organisation.ID,
		&organisation.UUID,
		&organisation.Name,
		&organisation.OwnerUserID,
		&organisation.CreatedAt,
		&organisation.CreatedBy,
		&organisation.UpdatedAt,
		&organisation.UpdatedBy,
		&organisation.DeletedAt,
		&organisation.DeletedBy,
	)
}
