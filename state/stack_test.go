package state

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/rds"
)

func Test_NewObject(t *testing.T) {
	obj := rds.RestoreDBClusterFromSnapshotInput{}
	order := 5
	objType := "rdsInstance"

	resp := NewObject(obj, order, objType)
	if resp.Order != order {
		t.Errorf("NewObject order expected %d got %d", resp.Order, order)
	}
}

func Test_NewStack(t *testing.T) {
	name := "foo"
	roname := "bar"

	o1 := Object{
		rds.RestoreDBClusterFromSnapshotInput{},
		1,
		"RDSCluster",
	}

	o2 := Object{
		rds.RestoreDBClusterFromS3Input{},
		1,
		"RDSCluster",
	}
	objects := []Object{o1, o2}

	resp := NewStack(name, roname, objects)
	if len(resp.Objects[1]) != 2 {
		t.Errorf("expected 2 got %d", len(resp.Objects[1]))
	}
}
