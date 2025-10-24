package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/crudbox/crudbox/internal/models"
)

// TokenManager generates and validates JWT tokens for authenticated requests.
type TokenManager struct {
	secret string
}

// TokenClaims represents the custom claims encoded within issued JWTs.
type TokenClaims struct {
	UserID   int64  `json:"userId"`
	UserUUID string `json:"userUuid"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// NewTokenManager creates a new TokenManager instance with the provided secret key.
func NewTokenManager(secret string) *TokenManager {
	return &TokenManager{secret: secret}
}

// GenerateToken signs a JWT for the supplied user.
func (m *TokenManager) GenerateToken(user *models.User) (string, error) {
	claims := TokenClaims{
		UserID:   user.ID,
		UserUUID: user.UUID,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.UUID,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secret))
}

// ValidateToken parses and validates an encoded JWT string.
func (m *TokenManager) ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
