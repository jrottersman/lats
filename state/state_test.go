package state

import (
	"os"
	"testing"
)

func TestInitState(t *testing.T) {
	filename := "/tmp/.stateConf.json"
	defer os.Remove(filename)

	err := initState(filename)
	if err != nil {
		t.Errorf("error creating state file %s", err)
	}
}
