package constants

import "time"

const (
	// JWT Settings
	JWTTokenExpiration = time.Hour * 24

	// Auth Related Messages
	InvalidCredentials = "Invalid credentials"
	UnauthorizedAccess = "Unauthorized access"
)
