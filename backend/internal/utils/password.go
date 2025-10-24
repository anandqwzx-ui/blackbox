package utils

import "golang.org/x/crypto/bcrypt"

// PasswordManager encapsulates password hashing concerns.
type PasswordManager struct{}

// NewPasswordManager returns a PasswordManager instance.
func NewPasswordManager() *PasswordManager {
	return &PasswordManager{}
}

// HashPassword hashes the supplied plain text password.
func (p *PasswordManager) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Compare verifies a password matches the stored hash.
func (p *PasswordManager) Compare(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
