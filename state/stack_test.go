package state

import (
	"encoding/gob"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

func Test_ReadObject(t *testing.T) {
	filename := "/tmp/foo"
	order := 5
	objType := LoneInstance

	defer os.Remove(filename)
	db := rds.RestoreDBInstanceFromDBSnapshotInput{
		DBInstanceIdentifier: aws.String("foo"),
	}
	r := EncodeRestoreDBInstanceFromDBSnapshotInput(&db)
	_, err := WriteOutput(filename, r)
	if err != nil {
		t.Errorf("failed to write output, %s", err)
	}

	resp := NewObject(filename, order, objType)
	i := resp.ReadObject()
	_, ok := i.(*rds.RestoreDBInstanceFromDBSnapshotInput)
	if !ok {
		t.Errorf("this should have been ok")
	}
}

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

func Test_Write(t *testing.T) {
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
	filename := "/tmp/foobar"
	defer os.Remove(filename)
	err := stack.Write(filename)
	if err != nil {
		t.Errorf("got %s expected nil", err)
	}
}

func Test_ReadStack(t *testing.T) {
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
	filename := "/tmp/foobar"
	defer os.Remove(filename)
	err := stack.Write(filename)
	if err != nil {
		t.Errorf("writing got %s expected nil", err)
	}
	sp, err := ReadStack(filename)
	if err != nil {
		t.Errorf("error reading stack %s", err)
	}
	if sp.Name != stack.Name {
		t.Errorf("got %s expected %s", sp.Name, stack.Name)
	}
}

func Test_DeleteStack(t *testing.T) {
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
	filename := "/tmp/foobar"
	err := stack.Write(filename)
	if err != nil {
		t.Errorf("writing got %s expected nil", err)
	}
	fe := DeleteStack(filename)
	if fe != nil {
		t.Errorf("delete stack error %s", fe)
	}
}
