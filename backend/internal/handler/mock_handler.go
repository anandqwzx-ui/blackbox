package handler

import (
    "encoding/json"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"

    "github.com/crudbox/crudbox/internal/repository"
    "github.com/crudbox/crudbox/internal/service"
)

// MockHandler resolves incoming requests into stored mock responses.
type MockHandler struct {
    endpoints *service.EndpointService
}

// NewMockHandler constructs a MockHandler.
func NewMockHandler(endpoints *service.EndpointService) *MockHandler {
    return &MockHandler{endpoints: endpoints}
}

// Handle looks up an endpoint definition by project code and incoming request path/method.
func (h *MockHandler) Handle(c *gin.Context) {
    projectCode := c.Param("projectCode")
    endpointPath := c.Param("endpointPath")
    endpointPath = strings.TrimPrefix(endpointPath, "/")

    endpoint, err := h.endpoints.ResolveMock(c.Request.Context(), projectCode, c.Request.Method, endpointPath)
    if err != nil {
        switch err {
        case repository.ErrProjectNotFound, repository.ErrEndpointNotFound:
            c.JSON(http.StatusNotFound, gin.H{"error": "mock endpoint not found"})
        default:
            c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resolve endpoint"})
        }
        return
    }

    headers := map[string]string{}
    if endpoint.ResponseHeaders != "" {
        _ = json.Unmarshal([]byte(endpoint.ResponseHeaders), &headers)
    }

    for key, value := range headers {
        c.Header(key, value)
    }

    if _, exists := headers["Content-Type"]; !exists {
        c.Header("Content-Type", "application/json")
    }

    c.Status(endpoint.ResponseStatus)
    if endpoint.ResponseBody != "" {
        _, _ = c.Writer.WriteString(endpoint.ResponseBody)
    }
    c.Abort()
}
