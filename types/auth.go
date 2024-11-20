package types

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required,min=2"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	User  UserResponse `json:"user"`
	Token TokenPair    `json:"tokens"`
}

// ErrorResponse Custom error responses
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
