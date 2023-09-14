package cmd

import (
	"sync"
	"testing"

	"github.com/jrottersman/lats/state"
)

func TestFindStack(t *testing.T) {
	sm := state.StateManager{
		Mu:             &sync.Mutex{},
		StateLocations: []state.StateKV{},
	}
	filename := "/tmp/foo"
	obj := "boo"
	ot := "bar"
	sm.UpdateState(obj, filename, ot)
	resp := FindStack(sm, "baz")
	if resp != nil {
		t.Errorf("Expected nil got %v", resp)
	}

}
