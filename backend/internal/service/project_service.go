package service

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/ctonew/mockapi/internal/models"
	"github.com/ctonew/mockapi/internal/repository"
	"github.com/ctonew/mockapi/internal/utils"
)

// ProjectService coordinates project lifecycle operations while enforcing access controls.
type ProjectService struct {
	repo       *repository.ProjectRepository
	orgRepo    *repository.OrganisationRepository
}

// ProjectResult represents a project bundled with the owning organisation's UUID for response models.
type ProjectResult struct {
	Project           *models.Project
	OrganisationUUID string
}

// NewProjectService constructs a ProjectService.
func NewProjectService(repo *repository.ProjectRepository, orgRepo *repository.OrganisationRepository) *ProjectService {
	return &ProjectService{repo: repo, orgRepo: orgRepo}
}

// Create establishes a new project under the provided organisation for the acting user.
func (s *ProjectService) Create(ctx context.Context, userID int64, organisationUUID, name, description string) (*ProjectResult, error) {
	organisation, err := s.orgRepo.GetByUUIDAndOwner(ctx, organisationUUID, userID)
	if err != nil {
		return nil, err
	}

	var (
		project *models.Project
		code    string
	)

	const maxAttempts = 5
	for attempt := 0; attempt < maxAttempts; attempt++ {
		code, err = utils.GenerateProjectCode()
		if err != nil {
			return nil, err
		}

		project, err = s.repo.Create(ctx, organisation.ID, userID, name, description, code, userID)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
				continue
			}
			return nil, err
		}

		return &ProjectResult{Project: project, OrganisationUUID: organisation.UUID}, nil
	}

	return nil, errors.New("failed to generate unique project code")
}

// List fetches projects belonging to the user.
func (s *ProjectService) List(ctx context.Context, userID int64) ([]ProjectResult, error) {
	records, err := s.repo.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	results := make([]ProjectResult, 0, len(records))
	for _, record := range records {
		project := record.Project
		results = append(results, ProjectResult{
			Project:           &project,
			OrganisationUUID: record.OrganisationUUID,
		})
	}

	return results, nil
}

// Update modifies a project after validating ownership.
func (s *ProjectService) Update(ctx context.Context, userID int64, projectUUID, name, description string) (*ProjectResult, error) {
	record, err := s.repo.GetByUUIDAndOwner(ctx, projectUUID, userID)
	if err != nil {
		return nil, err
	}

	updated, err := s.repo.Update(ctx, record.Project.ID, name, description, userID)
	if err != nil {
		return nil, err
	}

	return &ProjectResult{Project: updated, OrganisationUUID: record.OrganisationUUID}, nil
}

// Delete soft deletes a project ensuring the caller is the owner.
func (s *ProjectService) Delete(ctx context.Context, userID int64, projectUUID string) error {
	record, err := s.repo.GetByUUIDAndOwner(ctx, projectUUID, userID)
	if err != nil {
		return err
	}

	return s.repo.SoftDelete(ctx, record.Project.ID, userID)
}

// GetByUUIDAndOwner returns the project after verifying ownership.
func (s *ProjectService) GetByUUIDAndOwner(ctx context.Context, userID int64, projectUUID string) (*ProjectResult, error) {
	record, err := s.repo.GetByUUIDAndOwner(ctx, projectUUID, userID)
	if err != nil {
		return nil, err
	}

	return &ProjectResult{Project: &record.Project, OrganisationUUID: record.OrganisationUUID}, nil
}

// GetByCode is used by the mock handler to resolve a project via its short code.
func (s *ProjectService) GetByCode(ctx context.Context, code string) (*models.Project, error) {
	return s.repo.GetByCode(ctx, code)
}
