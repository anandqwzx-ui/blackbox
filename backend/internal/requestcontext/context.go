package requestcontext

import (
	"github.com/gin-gonic/gin"

	"github.com/ctonew/mockapi/internal/models"
)

const (
	userIDKey   = "auth.user.id"
	userUUIDKey = "auth.user.uuid"
	userEmailKey = "auth.user.email"
)

// SetUser stores authenticated user details in the Gin context for downstream handlers.
func SetUser(c *gin.Context, user *models.User) {
	c.Set(userIDKey, user.ID)
	c.Set(userUUIDKey, user.UUID)
	c.Set(userEmailKey, user.Email)
}

// UserID retrieves the authenticated user's numeric identifier.
func UserID(c *gin.Context) (int64, bool) {
	value, exists := c.Get(userIDKey)
	if !exists {
		return 0, false
	}

	id, ok := value.(int64)
	return id, ok
}

// UserUUID retrieves the authenticated user's UUID.
func UserUUID(c *gin.Context) (string, bool) {
	value, exists := c.Get(userUUIDKey)
	if !exists {
		return "", false
	}

	uuid, ok := value.(string)
	return uuid, ok
}

// UserEmail retrieves the authenticated user's email address.
func UserEmail(c *gin.Context) (string, bool) {
	value, exists := c.Get(userEmailKey)
	if !exists {
		return "", false
	}

	email, ok := value.(string)
	return email, ok
}
