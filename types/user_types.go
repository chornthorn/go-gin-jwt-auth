package types

type UpdateProfileRequest struct {
	Name string `json:"name" binding:"required"`
	// Add more fields as needed
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}
