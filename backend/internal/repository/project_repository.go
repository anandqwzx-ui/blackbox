package repository

import (
    "context"
    "database/sql"
    "errors"

    "github.com/jmoiron/sqlx"

    "github.com/crudbox/crudbox/internal/models"
)

// ErrProjectNotFound indicates a project record was not found or is inaccessible for the user.
var ErrProjectNotFound = errors.New("project not found")

// ProjectWithOrganisation couples a project with organisation metadata required by higher layers.
type ProjectWithOrganisation struct {
    Project           models.Project
    OrganisationUUID string
}

type projectWithOrganisationRow struct {
    models.Project    `db:""`
    OrganisationUUID string `db:"organisation_uuid"`
}

// ProjectRepository handles project persistence operations.
type ProjectRepository struct {
    db *sqlx.DB
}

// NewProjectRepository constructs a ProjectRepository.
func NewProjectRepository(db *sqlx.DB) *ProjectRepository {
    return &ProjectRepository{db: db}
}

// Create inserts a new project record.
func (r *ProjectRepository) Create(ctx context.Context, organisationID, ownerUserID int64, name, description, code string, actorID int64) (*models.Project, error) {
    query := `
        INSERT INTO projects (organisation_id, owner_user_id, name, description, code, created_by, updated_by)
        VALUES ($1, $2, $3, NULLIF($4, ''), $5, $6, $6)
        RETURNING id, uuid, organisation_id, owner_user_id, name, description, code,
                  created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
    `

    row := r.db.QueryRowxContext(ctx, query, organisationID, ownerUserID, name, description, code, actorID)
    project := &models.Project{}
    if err := row.StructScan(project); err != nil {
        return nil, err
    }

    return project, nil
}

// ListByUser returns active projects owned by a user with organisation identifiers.
func (r *ProjectRepository) ListByUser(ctx context.Context, ownerUserID int64) ([]ProjectWithOrganisation, error) {
    query := `
        SELECT p.id, p.uuid, p.organisation_id, p.owner_user_id, p.name, p.description, p.code,
               p.created_at, p.created_by, p.updated_at, p.updated_by, p.deleted_at, p.deleted_by,
               o.uuid
        FROM projects p
        JOIN organisations o ON o.id = p.organisation_id
        WHERE p.owner_user_id = $1 AND p.deleted_at IS NULL AND o.deleted_at IS NULL
        ORDER BY p.created_at DESC
    `

    var rowsData []projectWithOrganisationRow
    if err := r.db.SelectContext(ctx, &rowsData, query, ownerUserID); err != nil {
        return nil, err
    }

    projects := make([]ProjectWithOrganisation, 0, len(rowsData))
    for _, row := range rowsData {
        projects = append(projects, ProjectWithOrganisation{
            Project:           row.Project,
            OrganisationUUID: row.OrganisationUUID,
        })
    }

    return projects, nil
}

// GetByUUIDAndOwner returns a project owned by the specified user.
func (r *ProjectRepository) GetByUUIDAndOwner(ctx context.Context, projectUUID string, ownerUserID int64) (*ProjectWithOrganisation, error) {
    query := `
        SELECT p.id, p.uuid, p.organisation_id, p.owner_user_id, p.name, p.description, p.code,
               p.created_at, p.created_by, p.updated_at, p.updated_by, p.deleted_at, p.deleted_by,
               o.uuid
        FROM projects p
        JOIN organisations o ON o.id = p.organisation_id
        WHERE p.uuid = $1 AND p.owner_user_id = $2 AND p.deleted_at IS NULL AND o.deleted_at IS NULL
    `

    var row projectWithOrganisationRow
    if err := r.db.GetContext(ctx, &row, query, projectUUID, ownerUserID); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrProjectNotFound
        }
        return nil, err
    }

    return &ProjectWithOrganisation{Project: row.Project, OrganisationUUID: row.OrganisationUUID}, nil
}

// Update applies changes to an existing project that the owner has access to.
func (r *ProjectRepository) Update(ctx context.Context, projectID int64, name, description string, actorID int64) (*models.Project, error) {
    query := `
        UPDATE projects
        SET name = $1,
            description = NULLIF($2, ''),
            updated_at = CURRENT_TIMESTAMP,
            updated_by = $3
        WHERE id = $4
        RETURNING id, uuid, organisation_id, owner_user_id, name, description, code,
                  created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
    `

    row := r.db.QueryRowxContext(ctx, query, name, description, actorID, projectID)
    project := &models.Project{}
    if err := row.StructScan(project); err != nil {
        return nil, err
    }

    return project, nil
}

// SoftDelete marks a project as deleted.
func (r *ProjectRepository) SoftDelete(ctx context.Context, projectID int64, actorID int64) error {
    query := `
        UPDATE projects
        SET deleted_at = CURRENT_TIMESTAMP,
            deleted_by = $1
        WHERE id = $2 AND deleted_at IS NULL
    `

    result, err := r.db.ExecContext(ctx, query, actorID, projectID)
    if err != nil {
        return err
    }

    affected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if affected == 0 {
        return ErrProjectNotFound
    }

    return nil
}

// GetByCode returns a project by its short code.
func (r *ProjectRepository) GetByCode(ctx context.Context, code string) (*models.Project, error) {
    query := `
        SELECT id, uuid, organisation_id, owner_user_id, name, description, code,
               created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
        FROM projects
        WHERE code = $1 AND deleted_at IS NULL
    `

    project := &models.Project{}
    if err := r.db.GetContext(ctx, project, query, code); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrProjectNotFound
        }
        return nil, err
    }

    return project, nil
}

// GetByID returns a project by internal ID.
func (r *ProjectRepository) GetByID(ctx context.Context, id int64) (*models.Project, error) {
    query := `
        SELECT id, uuid, organisation_id, owner_user_id, name, description, code,
               created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
        FROM projects
        WHERE id = $1 AND deleted_at IS NULL
    `

    project := &models.Project{}
    if err := r.db.GetContext(ctx, project, query, id); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrProjectNotFound
        }
        return nil, err
    }

    return project, nil
}
