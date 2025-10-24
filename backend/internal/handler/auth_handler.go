package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"

    "github.com/ctonew/mockapi/internal/dto"
    "github.com/ctonew/mockapi/internal/repository"
    "github.com/ctonew/mockapi/internal/service"
)

// AuthHandler exposes authentication related HTTP endpoints.
type AuthHandler struct {
    auth *service.AuthService
}

// NewAuthHandler constructs an AuthHandler.
func NewAuthHandler(auth *service.AuthService) *AuthHandler {
    return &AuthHandler{auth: auth}
}

// RegisterRoutes registers authentication routes within the provided router group.
func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
    router.POST("/signup", h.signup)
    router.POST("/login", h.login)
}

func (h *AuthHandler) signup(c *gin.Context) {
    var request dto.SignupRequest
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, token, err := h.auth.Signup(c.Request.Context(), request.Email, request.Password)
    if err != nil {
        switch err {
        case service.ErrEmailAlreadyExists:
            c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
        default:
            c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to process signup"})
        }
        return
    }

    response := dto.AuthResponse{
        UserUUID: user.UUID,
        Email:    user.Email,
        Token:    token,
    }

    c.JSON(http.StatusCreated, response)
}

func (h *AuthHandler) login(c *gin.Context) {
    var request dto.LoginRequest
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, token, err := h.auth.Login(c.Request.Context(), request.Email, request.Password)
    if err != nil {
        switch err {
        case service.ErrInvalidCredentials:
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        case repository.ErrUserNotFound:
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        default:
            c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to process login"})
        }
        return
    }

    response := dto.AuthResponse{
        UserUUID: user.UUID,
        Email:    user.Email,
        Token:    token,
    }

    c.JSON(http.StatusOK, response)
}
