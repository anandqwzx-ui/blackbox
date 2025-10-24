package middleware

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"

    "github.com/crudbox/crudbox/internal/repository"
    "github.com/crudbox/crudbox/internal/requestcontext"
    "github.com/crudbox/crudbox/internal/utils"
)

// AuthMiddleware ensures requests are authenticated using Bearer JWT tokens.
type AuthMiddleware struct {
    tokens *utils.TokenManager
    users  *repository.UserRepository
}

// NewAuthMiddleware constructs an AuthMiddleware instance.
func NewAuthMiddleware(tokens *utils.TokenManager, users *repository.UserRepository) *AuthMiddleware {
    return &AuthMiddleware{tokens: tokens, users: users}
}

// Handler returns the Gin middleware handler function.
func (m *AuthMiddleware) Handler() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
            return
        }

        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
            return
        }

        claims, err := m.tokens.ValidateToken(parts[1])
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            return
        }

        user, err := m.users.GetByUUID(c.Request.Context(), claims.UserUUID)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
            return
        }

        requestcontext.SetUser(c, user)

        c.Next()
    }
}
