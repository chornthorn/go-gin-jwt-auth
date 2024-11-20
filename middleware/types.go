package middleware

// AuthenticatedUser represents the user data stored in gin context
type AuthenticatedUser struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// ContextKey type for context keys to avoid string collisions
type ContextKey string

const (
	// UserContextKey is the key used to store the user in the context
	UserContextKey ContextKey = "user"
	// TokenMetadataKey is the key used to store token metadata in the context
	TokenMetadataKey ContextKey = "token_metadata"
)
