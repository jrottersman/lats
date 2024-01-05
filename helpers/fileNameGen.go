package helpers

import (
	"fmt"

	"github.com/google/uuid"
)

// RandomStateFileName generates a gob file name starting with a UUID
func RandomStateFileName() *string {
	u := uuid.New()
	filename := fmt.Sprintf("%s.gob", u)
	return &filename
}

// InstanceName generates a new random instance name for when one isn't provided
func InstanceName() *string {
	u := uuid.New()
	instanceName := fmt.Sprintf("instance-%s", u)
	return &instanceName
}
