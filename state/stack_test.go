package state

import (
	"testing"
)

func Test_NewObject(t *testing.T) {
	filename := "/tmp/foo"
	order := 5
	objType := "rdsInstance"

	resp := NewObject(filename, order, objType)
	if resp.Order != order {
		t.Errorf("NewObject order expected %d got %d", resp.Order, order)
	}
}

func Test_NewStack(t *testing.T) {
	name := "foo"
	roname := "bar"

	o1 := Object{
		"tmp/foo",
		1,
		"RDSCluster",
	}

	o2 := Object{
		"tmp/bar",
		1,
		"RDSCluster",
	}
	objects := []Object{o1, o2}

	resp := NewStack(name, roname, objects)
	if len(resp.Objects[1]) != 2 {
		t.Errorf("expected 2 got %d", len(resp.Objects[1]))
	}
}
