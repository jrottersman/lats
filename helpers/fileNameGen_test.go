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

func TestSnapshotName(t *testing.T) {
	expected := "foo"
	s := SnapshotName(expected)
	if !strings.Contains(s, expected) {
		t.Errorf("string is %s should contain %s", s, expected)
	}
}
