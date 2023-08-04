package state

import (
	"os"
	"testing"
)

func TestInitState(t *testing.T) {
	filename := "/tmp/.stateConf.json"
	defer os.Remove(filename)

	err := InitState(filename)
	if err != nil {
		t.Errorf("error creating state file %s", err)
	}
}

func TestReadState(t *testing.T) {
	filename := "/tmp/.stateConf.json"
	defer os.Remove(filename)

	err := InitState(filename)
	if err != nil {
		t.Errorf("error creating state file %s", err)
	}
	s, err := ReadState(filename)
	if err != nil {
		t.Errorf("error creating state in memory: %s", err)
	}
	if len(s.StateLocations) > 0 {
		t.Errorf("expected length to be 0")
	}

}
