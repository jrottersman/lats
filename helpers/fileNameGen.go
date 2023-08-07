package helpers

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func RandomStateFileName() *string {
	u := uuid.New()
	filename := fmt.Sprintf(".state/%s.gob", u)
	return &filename
}

func SnapshotName(db string) string {
	t := time.Now().UTC().String()
	return fmt.Sprintf("%s-%s", db, t)
}
