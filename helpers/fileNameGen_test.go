package helpers

import (
	"strings"
	"testing"
)

func TestSnapshotName(t *testing.T) {
	expected := "foo"
	s := SnapshotName(expected)
	if !strings.Contains(s, expected) {
		t.Errorf("string is %s should contain %s", s, expected)
	}
}
