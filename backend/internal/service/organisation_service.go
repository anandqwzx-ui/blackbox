package service

import (
	"context"

	"github.com/ctonew/mockapi/internal/models"
	"github.com/ctonew/mockapi/internal/repository"
)

// OrganisationService coordinates organisation-related operations with necessary access checks.
type OrganisationService struct {
	repo *repository.OrganisationRepository
}

// NewOrganisationService constructs an OrganisationService.
func NewOrganisationService(repo *repository.OrganisationRepository) *OrganisationService {
	return &OrganisationService{repo: repo}
}

// Create creates a new organisation for the supplied user.
func (s *OrganisationService) Create(ctx context.Context, userID int64, name string) (*models.Organisation, error) {
	return s.repo.Create(ctx, name, userID, userID)
}

// List returns all organisations belonging to the supplied user.
func (s *OrganisationService) List(ctx context.Context, userID int64) ([]models.Organisation, error) {
	return s.repo.ListByUser(ctx, userID)
}

// Get ensures the user has access to the requested organisation.
func (s *OrganisationService) Get(ctx context.Context, userID int64, organisationUUID string) (*models.Organisation, error) {
	return s.repo.GetByUUIDAndOwner(ctx, organisationUUID, userID)
}
