package handler

import (
    "encoding/json"
    "errors"
    "net/http"

    "github.com/gin-gonic/gin"

    "github.com/crudbox/crudbox/internal/dto"
    "github.com/crudbox/crudbox/internal/repository"
    "github.com/crudbox/crudbox/internal/requestcontext"
    "github.com/crudbox/crudbox/internal/service"
)

// EndpointHandler manages CRUD operations for mock endpoints.
type EndpointHandler struct {
    endpoints *service.EndpointService
}

// NewEndpointHandler constructs an EndpointHandler.
func NewEndpointHandler(endpoints *service.EndpointService) *EndpointHandler {
    return &EndpointHandler{endpoints: endpoints}
}

// RegisterRoutes wires endpoint routes into the router group.
func (h *EndpointHandler) RegisterRoutes(router *gin.RouterGroup) {
    router.POST("/projects/:projectUuid/endpoints", h.create)
    router.GET("/projects/:projectUuid/endpoints", h.list)
    router.PUT("/endpoints/:endpointUuid", h.update)
    router.DELETE("/endpoints/:endpointUuid", h.delete)
}

func (h *EndpointHandler) create(c *gin.Context) {
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

    var request dto.CreateEndpointRequest
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    enabled := true
    if request.Enabled != nil {
        enabled = *request.Enabled
    }

    result, err := h.endpoints.Create(
        c.Request.Context(),
        userID,
        projectUUID,
        request.Name,
        request.Method,
        request.Path,
        request.ResponseStatus,
        request.ResponseBody,
        request.ResponseHeaders,
        enabled,
    )
    if err != nil {
        h.handleServiceError(c, err)
        return
    }

    c.JSON(http.StatusCreated, endpointToDTO(result))
}

func (h *EndpointHandler) list(c *gin.Context) {
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

    results, err := h.endpoints.List(c.Request.Context(), userID, projectUUID)
    if err != nil {
        h.handleServiceError(c, err)
        return
    }

    responses := make([]dto.EndpointResponse, 0, len(results))
    for i := range results {
        responses = append(responses, endpointToDTO(&results[i]))
    }

    c.JSON(http.StatusOK, responses)
}

func (h *EndpointHandler) update(c *gin.Context) {
    userID, ok := requestcontext.UserID(c)
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorised"})
        return
    }

    endpointUUID := c.Param("endpointUuid")
    if endpointUUID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "missing endpoint identifier"})
        return
    }

    var request dto.UpdateEndpointRequest
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    enabled := true
    if request.Enabled != nil {
        enabled = *request.Enabled
    }

    result, err := h.endpoints.Update(
        c.Request.Context(),
        userID,
        endpointUUID,
        request.Name,
        request.Method,
        request.Path,
        request.ResponseStatus,
        request.ResponseBody,
        request.ResponseHeaders,
        enabled,
    )
    if err != nil {
        h.handleServiceError(c, err)
        return
    }

    c.JSON(http.StatusOK, endpointToDTO(result))
}

func (h *EndpointHandler) delete(c *gin.Context) {
    userID, ok := requestcontext.UserID(c)
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorised"})
        return
    }

    endpointUUID := c.Param("endpointUuid")
    if endpointUUID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "missing endpoint identifier"})
        return
    }

    if err := h.endpoints.Delete(c.Request.Context(), userID, endpointUUID); err != nil {
        h.handleServiceError(c, err)
        return
    }

    c.Status(http.StatusNoContent)
}

func (h *EndpointHandler) handleServiceError(c *gin.Context, err error) {
    switch {
    case errors.Is(err, service.ErrEndpointForbidden):
        c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
    case errors.Is(err, service.ErrEndpointDuplicate):
        c.JSON(http.StatusConflict, gin.H{"error": "endpoint already exists for this method and path"})
    case errors.Is(err, repository.ErrProjectNotFound), errors.Is(err, repository.ErrEndpointNotFound):
        c.JSON(http.StatusNotFound, gin.H{"error": "resource not found"})
    default:
        c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
    }
}

func endpointToDTO(result *service.EndpointResult) dto.EndpointResponse {
    headers := map[string]string{}
    if result.Endpoint.ResponseHeaders != "" {
        _ = json.Unmarshal([]byte(result.Endpoint.ResponseHeaders), &headers)
    }

    var name *string
    if result.Endpoint.Name.Valid {
        name = &result.Endpoint.Name.String
    }

    return dto.EndpointResponse{
        UUID:           result.Endpoint.UUID,
        ProjectUUID:    result.ProjectUUID,
        Name:           name,
        Method:         result.Endpoint.Method,
        Path:           result.Endpoint.Path,
        ResponseStatus: result.Endpoint.ResponseStatus,
        ResponseBody:   result.Endpoint.ResponseBody,
        ResponseHeaders: headers,
        Enabled:        result.Endpoint.Enabled,
        CreatedAt:      result.Endpoint.CreatedAt,
        UpdatedAt:      result.Endpoint.UpdatedAt,
    }
}
