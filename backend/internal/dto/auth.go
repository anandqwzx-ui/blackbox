package dto

// SignupRequest captures the fields required to register a new user.
type SignupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginRequest captures the login credentials for a user.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse represents the payload returned after successful authentication.
type AuthResponse struct {
	UserUUID string `json:"userUuid"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}
