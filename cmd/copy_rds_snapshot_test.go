package cmd

import (
	"os"
	"sync"
	"testing"

	"github.com/jrottersman/lats/stack"
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
	resp, _ := FindStack(sm, "baz")
	if resp != nil {
		t.Errorf("Expected nil got %v", resp)
	}
	stk := stack.Stack{
		Name:                  "foo",
		RestorationObjectName: "stack",
	}
	stk.Write(filename)
	defer os.Remove(filename)
	sm.UpdateState("foo", filename, "stack")
	exp, err := FindStack(sm, "foo")
	if err != nil {
		t.Errorf("got %v as an error", err)
	}
	if exp.Name != stk.Name {
		t.Errorf("got %s expected %s", exp.Name, stk.Name)
	}
}
