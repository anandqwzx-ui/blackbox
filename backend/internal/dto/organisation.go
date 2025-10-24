package dto

import "time"

// CreateOrganisationRequest represents the payload for creating a new organisation.
type CreateOrganisationRequest struct {
	Name string `json:"name" binding:"required"`
}

// OrganisationResponse is returned to clients representing their organisation data.
type OrganisationResponse struct {
	UUID      string    `json:"uuid"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
