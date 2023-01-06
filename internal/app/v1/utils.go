package v1

import "github.com/google/uuid"

// generateUserID creates a unique user ID
func generateUserID() string {
	return uuid.New().String()
}
