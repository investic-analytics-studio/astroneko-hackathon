package utils

import (
	"fmt"

	"github.com/google/uuid"
)

// ValidateAndParseUUID validates and parses a UUID string
func ValidateAndParseUUID(id string) (uuid.UUID, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid UUID format: %w", err)
	}
	return parsedID, nil
}

// IsValidUUID checks if a string is a valid UUID
func IsValidUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}

// GenerateUUID generates a new UUID
func GenerateUUID() uuid.UUID {
	return uuid.New()
}
