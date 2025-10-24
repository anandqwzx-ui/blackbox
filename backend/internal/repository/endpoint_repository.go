package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ctonew/mockapi/internal/models"
)

// ErrEndpointNotFound signals an endpoint could not be located.
var ErrEndpointNotFound = errors.New("endpoint not found")

// EndpointRepository manages storage for mock endpoints.
type EndpointRepository struct {
	db *sql.DB
}

// NewEndpointRepository constructs an EndpointRepository.
func NewEndpointRepository(db *sql.DB) *EndpointRepository {
	return &EndpointRepository{db: db}
}

// Create inserts a new endpoint record.
func (r *EndpointRepository) Create(ctx context.Context, projectID int64, name, method, path string, responseStatus int, responseBody, headersJSON string, enabled bool, actorID int64) (*models.Endpoint, error) {
	query := `
		INSERT INTO endpoints (project_id, name, method, path, response_status, response_body, response_headers, enabled, created_by, updated_by)
		VALUES ($1, NULLIF($2, ''), $3, $4, $5, $6, $7::jsonb, $8, $9, $9)
		RETURNING id, uuid, project_id, name, method, path, response_status, response_body, response_headers, enabled,
		          created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
	`

	row := r.db.QueryRowContext(ctx, query, projectID, name, method, path, responseStatus, responseBody, headersJSON, enabled, actorID)
	endpoint := &models.Endpoint{}
	if err := scanEndpoint(row, endpoint); err != nil {
		return nil, err
	}

	return endpoint, nil
}

// ListByProject returns endpoints for a given project.
func (r *EndpointRepository) ListByProject(ctx context.Context, projectID int64) ([]models.Endpoint, error) {
	query := `
		SELECT id, uuid, project_id, name, method, path, response_status, response_body, response_headers, enabled,
		       created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
		FROM endpoints
		WHERE project_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var endpoints []models.Endpoint
	for rows.Next() {
		var endpoint models.Endpoint
		if err := rows.Scan(
			&endpoint.ID,
			&endpoint.UUID,
			&endpoint.ProjectID,
			&endpoint.Name,
			&endpoint.Method,
			&endpoint.Path,
			&endpoint.ResponseStatus,
			&endpoint.ResponseBody,
			&endpoint.ResponseHeaders,
			&endpoint.Enabled,
			&endpoint.CreatedAt,
			&endpoint.CreatedBy,
			&endpoint.UpdatedAt,
			&endpoint.UpdatedBy,
			&endpoint.DeletedAt,
			&endpoint.DeletedBy,
		); err != nil {
			return nil, err
		}
		endpoints = append(endpoints, endpoint)
	}

	return endpoints, rows.Err()
}

// GetByUUID loads an endpoint by UUID.
func (r *EndpointRepository) GetByUUID(ctx context.Context, endpointUUID string) (*models.Endpoint, error) {
	query := `
		SELECT id, uuid, project_id, name, method, path, response_status, response_body, response_headers, enabled,
		       created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
		FROM endpoints
		WHERE uuid = $1 AND deleted_at IS NULL
	`

	row := r.db.QueryRowContext(ctx, query, endpointUUID)
	endpoint := &models.Endpoint{}
	if err := scanEndpoint(row, endpoint); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEndpointNotFound
		}
		return nil, err
	}

	return endpoint, nil
}

// Update modifies endpoint details.
func (r *EndpointRepository) Update(ctx context.Context, endpointID int64, name, method, path string, responseStatus int, responseBody, headersJSON string, enabled bool, actorID int64) (*models.Endpoint, error) {
	query := `
		UPDATE endpoints
		SET name = NULLIF($1, ''),
		    method = $2,
		    path = $3,
		    response_status = $4,
		    response_body = $5,
		    response_headers = $6::jsonb,
		    enabled = $7,
		    updated_at = CURRENT_TIMESTAMP,
		    updated_by = $8
		WHERE id = $9 AND deleted_at IS NULL
		RETURNING id, uuid, project_id, name, method, path, response_status, response_body, response_headers, enabled,
		          created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
	`

	row := r.db.QueryRowContext(ctx, query, name, method, path, responseStatus, responseBody, headersJSON, enabled, actorID, endpointID)
	endpoint := &models.Endpoint{}
	if err := scanEndpoint(row, endpoint); err != nil {
		return nil, err
	}

	return endpoint, nil
}

// SoftDelete marks an endpoint as deleted.
func (r *EndpointRepository) SoftDelete(ctx context.Context, endpointID int64, actorID int64) error {
	query := `
		UPDATE endpoints
		SET deleted_at = CURRENT_TIMESTAMP,
		    deleted_by = $1
		WHERE id = $2 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, actorID, endpointID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrEndpointNotFound
	}

	return nil
}

// FindByProjectAndPath locates an endpoint by project, method, and path.
func (r *EndpointRepository) FindByProjectAndPath(ctx context.Context, projectID int64, method, path string) (*models.Endpoint, error) {
	query := `
		SELECT id, uuid, project_id, name, method, path, response_status, response_body, response_headers, enabled,
		       created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
		FROM endpoints
		WHERE project_id = $1 AND method = $2 AND path = $3 AND deleted_at IS NULL AND enabled = TRUE
	`

	row := r.db.QueryRowContext(ctx, query, projectID, method, path)
	endpoint := &models.Endpoint{}
	if err := scanEndpoint(row, endpoint); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEndpointNotFound
		}
		return nil, err
	}

	return endpoint, nil
}

func scanEndpoint(row *sql.Row, endpoint *models.Endpoint) error {
	return row.Scan(
		&endpoint.ID,
		&endpoint.UUID,
		&endpoint.ProjectID,
		&endpoint.Name,
		&endpoint.Method,
		&endpoint.Path,
		&endpoint.ResponseStatus,
		&endpoint.ResponseBody,
		&endpoint.ResponseHeaders,
		&endpoint.Enabled,
		&endpoint.CreatedAt,
		&endpoint.CreatedBy,
		&endpoint.UpdatedAt,
		&endpoint.UpdatedBy,
		&endpoint.DeletedAt,
		&endpoint.DeletedBy,
	)
}
