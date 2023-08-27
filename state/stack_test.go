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
