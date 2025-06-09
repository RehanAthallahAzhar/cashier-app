package helpers

import (
	"github.com/google/uuid"
)

// GenerateNewUserID menghasilkan UUID baru sebagai string
func GenerateNewUserID() string {
	return uuid.New().String()
}
