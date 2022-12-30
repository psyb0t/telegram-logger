package v1

import "github.com/google/uuid"

// generateUserToken creates a unique user token
func generateUserToken() string {
	return uuid.New().String()
}
