package state

import (
	"encoding/json"
	"fmt"
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

func TestSyncState(t *testing.T) {
	var mu sync.Mutex
	var s []stateKV
	sm := StateManager{
		mu,
		s,
	}

	filename := "/tmp/foo"
	obj := "boo"
	ot := "baz"
	defer os.Remove(filename)
	sm.UpdateState(obj, filename, ot)
	sm.SyncState(filename)
	f, err := os.ReadFile(filename)
	if err != nil {
		t.Errorf("couldn't read file got %s", err)
	}
	var sf []stateKV
	err = json.Unmarshal(f, &sf)
	if err != nil {
		fmt.Printf("Couldn't unmarshall the json %s", err)
	}
	if sf[0].Object != obj {
		t.Errorf("got %s expected %s", sf[0].Object, obj)
	}

}
