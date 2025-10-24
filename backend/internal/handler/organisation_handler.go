package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ctonew/mockapi/internal/dto"
	"github.com/ctonew/mockapi/internal/requestcontext"
	"github.com/ctonew/mockapi/internal/service"
)

// OrganisationHandler manages organisation endpoints.
type OrganisationHandler struct {
	organisations *service.OrganisationService
}

// NewOrganisationHandler creates an OrganisationHandler instance.
func NewOrganisationHandler(organisations *service.OrganisationService) *OrganisationHandler {
	return &OrganisationHandler{organisations: organisations}
}

// RegisterRoutes registers organisation routes under the provided group.
func (h *OrganisationHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/organisations", h.create)
	router.GET("/organisations", h.list)
}

func (h *OrganisationHandler) create(c *gin.Context) {
	userID, ok := requestcontext.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorised"})
		return
	}

	var request dto.CreateOrganisationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	organisation, err := h.organisations.Create(c.Request.Context(), userID, request.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to create organisation"})
		return
	}

	response := dto.OrganisationResponse{
		UUID:      organisation.UUID,
		Name:      organisation.Name,
		CreatedAt: organisation.CreatedAt,
		UpdatedAt: organisation.UpdatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

func (h *OrganisationHandler) list(c *gin.Context) {
	userID, ok := requestcontext.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorised"})
		return
	}

	organisations, err := h.organisations.List(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch organisations"})
		return
	}

	responses := make([]dto.OrganisationResponse, 0, len(organisations))
	for _, organisation := range organisations {
		responses = append(responses, dto.OrganisationResponse{
			UUID:      organisation.UUID,
			Name:      organisation.Name,
			CreatedAt: organisation.CreatedAt,
			UpdatedAt: organisation.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, responses)
}
