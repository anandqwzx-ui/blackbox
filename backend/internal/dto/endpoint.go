package dto

import "time"

// CreateEndpointRequest captures the information needed to create a mock endpoint.
type CreateEndpointRequest struct {
	Name            string            `json:"name"`
	Method          string            `json:"method" binding:"required"`
	Path            string            `json:"path" binding:"required"`
	ResponseStatus  int               `json:"responseStatus" binding:"required"`
	ResponseBody    string            `json:"responseBody"`
	ResponseHeaders map[string]string `json:"responseHeaders"`
	Enabled         *bool             `json:"enabled"`
}

// UpdateEndpointRequest holds updateable fields for an endpoint.
type UpdateEndpointRequest struct {
	Name            string            `json:"name"`
	Method          string            `json:"method" binding:"required"`
	Path            string            `json:"path" binding:"required"`
	ResponseStatus  int               `json:"responseStatus" binding:"required"`
	ResponseBody    string            `json:"responseBody"`
	ResponseHeaders map[string]string `json:"responseHeaders"`
	Enabled         *bool             `json:"enabled"`
}

// EndpointResponse is returned after endpoint operations.
type EndpointResponse struct {
	UUID            string            `json:"uuid"`
	ProjectUUID     string            `json:"projectUuid"`
	Name            *string           `json:"name"`
	Method          string            `json:"method"`
	Path            string            `json:"path"`
	ResponseStatus  int               `json:"responseStatus"`
	ResponseBody    string            `json:"responseBody"`
	ResponseHeaders map[string]string `json:"responseHeaders"`
	Enabled         bool              `json:"enabled"`
	CreatedAt       time.Time         `json:"createdAt"`
	UpdatedAt       time.Time         `json:"updatedAt"`
}
