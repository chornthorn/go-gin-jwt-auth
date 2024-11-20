package types

// AuthenticatedUser represents the authenticated user in the context
type AuthenticatedUser struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}
