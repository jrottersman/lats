package helpers

import (
	"fmt"

	"github.com/google/uuid"
)

func RandomStateFileName() *string {
	u := uuid.New()
	filename := fmt.Sprintf("%s.gob", u)
	return &filename
}

func InstanceName() *string {
	u := uuid.New()
	instanceName := fmt.Sprintf("instance-%s", u)
	return &instanceName
}
