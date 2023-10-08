package stack_test

import (
	"encoding/gob"
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/jrottersman/lats/stack"
	"github.com/jrottersman/lats/state"
)

func Test_ReadObject(t *testing.T) {
	filename := "/tmp/foo"
	order := 5
	objType := stack.LoneInstance

	defer os.Remove(filename)
	db := rds.RestoreDBInstanceFromDBSnapshotInput{
		DBInstanceIdentifier: aws.String("foo"),
	}
	r := state.EncodeRestoreDBInstanceFromDBSnapshotInput(&db)
	_, err := state.WriteOutput(filename, r)
	if err != nil {
		t.Errorf("failed to write output, %s", err)
	}

	resp := stack.NewObject(filename, order, objType)
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

	resp := stack.NewObject(filename, order, objType)
	if resp.Order != order {
		t.Errorf("NewObject order expected %d got %d", resp.Order, order)
	}
}

func Test_NewStack(t *testing.T) {
	name := "foo"
	roname := "bar"

	o1 := stack.Object{
		"tmp/foo",
		1,
		"RDSCluster",
	}

	o2 := stack.Object{
		"tmp/bar",
		1,
		"RDSCluster",
	}
	objects := []stack.Object{o1, o2}

	resp := stack.NewStack(name, roname, objects)
	if len(resp.Objects[1]) != 2 {
		t.Errorf("expected 2 got %d", len(resp.Objects[1]))
	}
}

func Test_Encoder(t *testing.T) {
	name := "foo"
	roname := "bar"

	o1 := stack.Object{
		"tmp/foo",
		1,
		"RDSCluster",
	}

	o2 := stack.Object{
		"tmp/bar",
		1,
		"RDSCluster",
	}
	objects := []stack.Object{o1, o2}

	mStack := stack.NewStack(name, roname, objects)

	r, err := mStack.Encoder()
	if err != nil {
		t.Errorf("encode error: %s", err)
	}
	var result stack.Stack
	dec := gob.NewDecoder(r)
	err = dec.Decode(&result)
	if err != nil {
		t.Errorf("decode error: %s", err)
	}
	if result.Name != mStack.Name {
		t.Errorf("got %s expected %s", result.Name, mStack.Name)
	}
}

func Test_Write(t *testing.T) {
	name := "foo"
	roname := "bar"

	o1 := stack.Object{
		"tmp/foo",
		1,
		"RDSCluster",
	}

	o2 := stack.Object{
		"tmp/bar",
		1,
		"RDSCluster",
	}
	objects := []stack.Object{o1, o2}

	mStack := stack.NewStack(name, roname, objects)
	filename := "/tmp/foobar"
	defer os.Remove(filename)
	err := mStack.Write(filename)
	if err != nil {
		t.Errorf("got %s expected nil", err)
	}
}

func Test_ReadStack(t *testing.T) {
	name := "foo"
	roname := "bar"

	o1 := stack.Object{
		"tmp/foo",
		1,
		"RDSCluster",
	}

	o2 := stack.Object{
		"tmp/bar",
		1,
		"RDSCluster",
	}
	objects := []stack.Object{o1, o2}

	mStack := stack.NewStack(name, roname, objects)
	filename := "/tmp/foobar"
	defer os.Remove(filename)
	err := mStack.Write(filename)
	if err != nil {
		t.Errorf("writing got %s expected nil", err)
	}
	sp, err := stack.ReadStack(filename)
	if err != nil {
		t.Errorf("error reading stack %s", err)
	}
	if sp.Name != mStack.Name {
		t.Errorf("got %s expected %s", sp.Name, mStack.Name)
	}
}

func Test_DeleteStack(t *testing.T) {
	name := "foo"
	roname := "bar"

	o1 := stack.Object{
		"tmp/foo",
		1,
		"RDSCluster",
	}

	o2 := stack.Object{
		"tmp/bar",
		1,
		"RDSCluster",
	}
	objects := []stack.Object{o1, o2}

	mStack := stack.NewStack(name, roname, objects)
	filename := "/tmp/foobar"
	err := mStack.Write(filename)
	if err != nil {
		t.Errorf("writing got %s expected nil", err)
	}
	fe := stack.DeleteStack(filename)
	if fe != nil {
		t.Errorf("delete stack error %s", fe)
	}
}

func TestStack_SortStack(t *testing.T) {
	type fields struct {
		Name                  string
		RestorationObjectName string
		Objects               map[int][]stack.Object
	}
	objs := make(map[int][]stack.Object)
	objs[1] = []stack.Object{}
	objs[2] = []stack.Object{}
	objs[3] = []stack.Object{}
	expected := []int{1, 2, 3}
	tests := []struct {
		name   string
		fields fields
		want   *[]int
	}{
		{name: "Test", fields: fields{Objects: objs}, want: &expected},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := stack.Stack{
				Name:                  tt.fields.Name,
				RestorationObjectName: tt.fields.RestorationObjectName,
				Objects:               tt.fields.Objects,
			}
			if got := s.SortStack(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Stack.SortStack() = %v, want %v", got, tt.want)
			}
		})
	}
}
