package uuid

import (
	"fmt"

	"github.com/google/uuid"
)

func GenerateUUID() (string, error) {
	// Generate a new UUID and return it as a string.
	newUUID, err := uuid.NewRandom()
	if err != nil {
		fmt.Printf("Error generating UUID: %v\n", err)
		return "", err
	}
	return newUUID.String(), nil
}
