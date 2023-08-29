package state

import (
	"encoding/gob"
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

func Test_Encoder(t *testing.T) {
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

	stack := NewStack(name, roname, objects)

	r, err := stack.Encoder()
	if err != nil {
		t.Errorf("encode error: %s", err)
	}
	var result Stack
	dec := gob.NewDecoder(r)
	err = dec.Decode(&result)
	if err != nil {
		t.Errorf("decode error: %s", err)
	}
	if result.Name != stack.Name {
		t.Errorf("got %s expected %s", result.Name, stack.Name)
	}
}
