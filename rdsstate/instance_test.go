package rdsstate

import (
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/jrottersman/lats/pgstate"
	"github.com/jrottersman/lats/stack"
	"github.com/jrottersman/lats/state"
)

func TestGenerateRDSInstanceStack(t *testing.T) {
	type args struct {
		i InstanceStackInputs
	}
	r := state.RDSRestorationStore{
		Snapshot: &types.DBSnapshot{DBSnapshotIdentifier: aws.String("boo")},
		Instance: &types.DBInstance{},
	}
	inputs := InstanceStackInputs{
		R:                      r,
		StackName:              "bar",
		InstanceFileName:       "/tmp/foo.gob",
		ParameterFileName:      "/tmp/bar.gob",
		SecurityGroupsFileName: "/tmp/bat.gob",
		ParameterGroups:        []pgstate.ParameterGroup{},
	}
	arg := args{
		i: inputs,
	}
	defer os.Remove("/tmp/foo.gob")
	defer os.Remove("/tmp/bar.gob")

	obj := stack.Object{
		FileName: "/tmp/foo.gob",
		Order:    2,
		ObjType:  stack.LoneInstance,
	}
	pobj := stack.Object{
		FileName: "/tmp/bar.gob",
		Order:    1,
		ObjType:  stack.DBParameterGroup,
	}
	oMap := make(map[int][]stack.Object)
	oMap[2] = []stack.Object{obj}
	oMap[1] = []stack.Object{pobj}
	expected := stack.Stack{
		Name:                  "bar",
		RestorationObjectName: stack.LoneInstance,
		Objects:               oMap,
	}

	tests := []struct {
		name    string
		args    args
		want    *stack.Stack
		wantErr bool
	}{
		{"first", arg, &expected, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateRDSInstanceStack(tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateRDSInstanceStack() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateRDSInstanceStack() = %v, want %v", got, tt.want)
			}
		})
	}
}
