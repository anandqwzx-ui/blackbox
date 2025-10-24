package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/crudbox/crudbox/internal/dto"
	"github.com/crudbox/crudbox/internal/requestcontext"
	"github.com/crudbox/crudbox/internal/service"
)

// ProjectHandler manages project-related HTTP endpoints.
type ProjectHandler struct {
	projects *service.ProjectService
}

// NewProjectHandler constructs a ProjectHandler.
func NewProjectHandler(projects *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{projects: projects}
}

// RegisterRoutes registers project endpoints under the provided router group.
func (h *ProjectHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/projects", h.create)
	router.GET("/projects", h.list)
	router.PUT("/projects/:projectUuid", h.update)
	router.DELETE("/projects/:projectUuid", h.delete)
}

func (h *ProjectHandler) create(c *gin.Context) {
	userID, ok := requestcontext.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorised"})
		return
	}

	var request dto.CreateProjectRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.projects.Create(c.Request.Context(), userID, request.OrganisationUUID, request.Name, request.Description)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "unable to create project"})
		return
	}

	c.JSON(http.StatusCreated, projectToDTO(result))
}

func (h *ProjectHandler) list(c *gin.Context) {
	userID, ok := requestcontext.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorised"})
		return
	}

	projects, err := h.projects.List(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch projects"})
		return
	}

	responses := make([]dto.ProjectResponse, 0, len(projects))
	for i := range projects {
		responses = append(responses, projectToDTO(&projects[i]))
	}

	c.JSON(http.StatusOK, responses)
}

func (h *ProjectHandler) update(c *gin.Context) {
	userID, ok := requestcontext.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorised"})
		return
	}

	projectUUID := c.Param("projectUuid")
	if projectUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing project identifier"})
		return
	}

	var request dto.UpdateProjectRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.projects.Update(c.Request.Context(), userID, projectUUID, request.Name, request.Description)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "unable to update project"})
		return
	}

	c.JSON(http.StatusOK, projectToDTO(result))
}

func (h *ProjectHandler) delete(c *gin.Context) {
	userID, ok := requestcontext.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorised"})
		return
	}

	projectUUID := c.Param("projectUuid")
	if projectUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing project identifier"})
		return
	}

	if err := h.projects.Delete(c.Request.Context(), userID, projectUUID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "unable to delete project"})
		return
	}

	c.Status(http.StatusNoContent)
}

func projectToDTO(result *service.ProjectResult) dto.ProjectResponse {
	var description *string
	if result.Project.Description.Valid {
		description = &result.Project.Description.String
	}

	return dto.ProjectResponse{
		UUID:             result.Project.UUID,
		OrganisationUUID: result.OrganisationUUID,
		Name:             result.Project.Name,
		Description:      description,
		Code:             result.Project.Code,
		CreatedAt:        result.Project.CreatedAt,
		UpdatedAt:        result.Project.UpdatedAt,
	}
}
