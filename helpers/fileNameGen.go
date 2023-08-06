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
