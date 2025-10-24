package service

import (
    "context"
    "encoding/json"
    "errors"
    "strings"

    "github.com/crudbox/crudbox/internal/models"
    "github.com/crudbox/crudbox/internal/repository"
)

// ErrEndpointForbidden indicates the user attempted to access an endpoint outside their ownership.
var ErrEndpointForbidden = errors.New("endpoint access forbidden")

// ErrEndpointDuplicate indicates an endpoint with the same method and path already exists for the project.
var ErrEndpointDuplicate = errors.New("endpoint already exists for this method and path")

// EndpointResult bundles an endpoint with its parent project UUID for response mapping.
type EndpointResult struct {
    Endpoint    *models.Endpoint
    ProjectUUID string
}

// EndpointService manages endpoint operations with proper authorisation checks.
type EndpointService struct {
    repo        *repository.EndpointRepository
    projectRepo *repository.ProjectRepository
}

// NewEndpointService constructs an EndpointService.
func NewEndpointService(repo *repository.EndpointRepository, projectRepo *repository.ProjectRepository) *EndpointService {
    return &EndpointService{repo: repo, projectRepo: projectRepo}
}

// Create creates an endpoint within a project owned by the user.
func (s *EndpointService) Create(ctx context.Context, userID int64, projectUUID, name, method, path string, status int, body string, headers map[string]string, enabled bool) (*EndpointResult, error) {
    record, err := s.projectRepo.GetByUUIDAndOwner(ctx, projectUUID, userID)
    if err != nil {
        return nil, err
    }

    normalisedPath := normalisePath(path)
    normalisedMethod := strings.ToUpper(method)
    headersJSON, err := encodeHeaders(headers)
    if err != nil {
        return nil, err
    }

    endpoint, err := s.repo.Create(ctx, record.Project.ID, name, normalisedMethod, normalisedPath, status, body, headersJSON, enabled, userID)
    if err != nil {
        if errors.Is(err, repository.ErrEndpointConflict) {
            return nil, ErrEndpointDuplicate
        }
        return nil, err
    }

    return &EndpointResult{Endpoint: endpoint, ProjectUUID: record.Project.UUID}, nil
}

// List returns endpoints for a project owned by the user.
func (s *EndpointService) List(ctx context.Context, userID int64, projectUUID string) ([]EndpointResult, error) {
    record, err := s.projectRepo.GetByUUIDAndOwner(ctx, projectUUID, userID)
    if err != nil {
        return nil, err
    }

    endpoints, err := s.repo.ListByProject(ctx, record.Project.ID)
    if err != nil {
        return nil, err
    }

    results := make([]EndpointResult, 0, len(endpoints))
    for i := range endpoints {
        endpoint := endpoints[i]
        results = append(results, EndpointResult{Endpoint: &endpoint, ProjectUUID: record.Project.UUID})
    }

    return results, nil
}

// Update modifies an endpoint owned by the user.
func (s *EndpointService) Update(ctx context.Context, userID int64, endpointUUID, name, method, path string, status int, body string, headers map[string]string, enabled bool) (*EndpointResult, error) {
    endpoint, err := s.repo.GetByUUID(ctx, endpointUUID)
    if err != nil {
        return nil, err
    }

    project, err := s.projectRepo.GetByID(ctx, endpoint.ProjectID)
    if err != nil {
        return nil, err
    }

    if project.OwnerUserID != userID {
        return nil, ErrEndpointForbidden
    }

    normalisedPath := normalisePath(path)
    normalisedMethod := strings.ToUpper(method)
    headersJSON, err := encodeHeaders(headers)
    if err != nil {
        return nil, err
    }

    updated, err := s.repo.Update(ctx, endpoint.ID, name, normalisedMethod, normalisedPath, status, body, headersJSON, enabled, userID)
    if err != nil {
        if errors.Is(err, repository.ErrEndpointConflict) {
            return nil, ErrEndpointDuplicate
        }
        return nil, err
    }

    return &EndpointResult{Endpoint: updated, ProjectUUID: project.UUID}, nil
}

// Delete removes an endpoint owned by the user via soft delete.
func (s *EndpointService) Delete(ctx context.Context, userID int64, endpointUUID string) error {
    endpoint, err := s.repo.GetByUUID(ctx, endpointUUID)
    if err != nil {
        return err
    }

    project, err := s.projectRepo.GetByID(ctx, endpoint.ProjectID)
    if err != nil {
        return err
    }

    if project.OwnerUserID != userID {
        return ErrEndpointForbidden
    }

    return s.repo.SoftDelete(ctx, endpoint.ID, userID)
}

// ResolveMock fetches the endpoint response definition for a project code and request path/method.
func (s *EndpointService) ResolveMock(ctx context.Context, projectCode, method, path string) (*models.Endpoint, error) {
    project, err := s.projectRepo.GetByCode(ctx, projectCode)
    if err != nil {
        return nil, err
    }

    return s.repo.FindByProjectAndPath(ctx, project.ID, strings.ToUpper(method), normalisePath(path))
}

func encodeHeaders(headers map[string]string) (string, error) {
    if headers == nil {
        headers = map[string]string{}
    }
    bytes, err := json.Marshal(headers)
    if err != nil {
        return "", err
    }
    return string(bytes), nil
}

func normalisePath(path string) string {
    if path == "" {
        return "/"
    }
    if !strings.HasPrefix(path, "/") {
        return "/" + path
    }
    return path
}
