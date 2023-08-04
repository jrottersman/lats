package state

import (
	"os"
	"sync"
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

func TestUpdateState(t *testing.T) {
	var mu sync.Mutex
	var s []stateKV
	sm := StateManager{
		mu,
		s,
	}
	filename := "/tmp/foo"
	obj := "boo"
	sm.UpdateState(obj, filename)
	if sm.StateLocations[0].Object != "boo" {
		t.Errorf("got %s expected %s\n", sm.StateLocations[0].Object, obj)
	}
}
