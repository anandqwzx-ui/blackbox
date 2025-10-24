package dto

import "time"

// CreateProjectRequest captures the data needed to create a project.
type CreateProjectRequest struct {
	OrganisationUUID string `json:"organisationUuid" binding:"required"`
	Name             string `json:"name" binding:"required"`
	Description      string `json:"description"`
}

// UpdateProjectRequest represents the allowable fields for updating a project.
type UpdateProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// ProjectResponse is the payload returned to clients for project resources.
type ProjectResponse struct {
	UUID              string     `json:"uuid"`
	OrganisationUUID  string     `json:"organisationUuid"`
	Name              string     `json:"name"`
	Description       *string    `json:"description"`
	Code              string     `json:"code"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
}
