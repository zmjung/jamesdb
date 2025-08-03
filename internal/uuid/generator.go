package uuid

import (
	"github.com/google/uuid"
)

func GenerateUUID() (string, error) {
	// Generate a new UUID and return it as a string.
	newUUID, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return newUUID.String(), nil
}

func GenerateShortID() (string, error) {
	// Generate a 8-character thread ID from a UUID.
	// This is a simplified version and is not guarenteed to be unique.
	uuid, err := GenerateUUID()
	if err != nil {
		return "", err
	}
	return uuid[:8], nil
}
