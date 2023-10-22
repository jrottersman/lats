package helpers

import (
	"strings"
	"testing"
)

func TestRandomStateFileName(t *testing.T) {
	s := RandomStateFileName()
	if !strings.Contains(*s, "gob") {
		t.Errorf("string should contain gob instead looks like: %s", *s)
	}
}
