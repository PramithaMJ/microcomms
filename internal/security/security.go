package security

import (
	"fmt"
	"strings"
)

// Mock function to simulate JWT validation
func ValidateJWT(token string) (bool, error) {
	if strings.HasPrefix(token, "Bearer ") {
		return true, nil
	}
	return false, fmt.Errorf("invalid JWT token")
}
