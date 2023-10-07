package rdsstate

import (
	"encoding/gob"
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	laws "github.com/jrottersman/lats/aws"
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
		R:                 r,
		StackName:         "bar",
		InstanceFileName:  "/tmp/foo.gob",
		ParameterFileName: "/tmp/bar.gob",
		ParameterGroups:   []laws.ParameterGroup{},
	}
	arg := args{
		i: inputs,
	}
	defer os.Remove("/tmp/foo.gob")
	defer os.Remove("/tmp/bar.gob")

	obj := state.Object{
		FileName: "/tmp/foo.gob",
		Order:    2,
		ObjType:  state.LoneInstance,
	}
	pobj := state.Object{
		FileName: "/tmp/bar.gob",
		Order:    1,
		ObjType:  state.DBParameterGroup,
	}
	oMap := make(map[int][]state.Object)
	oMap[2] = []state.Object{obj}
	oMap[1] = []state.Object{pobj}
	expected := state.Stack{
		Name:                  "bar",
		RestorationObjectName: state.LoneInstance,
		Objects:               oMap,
	}

	tests := []struct {
		name    string
		args    args
		want    *state.Stack
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

func Test_encodeParameterGroups(t *testing.T) {
	pg := laws.ParameterGroup{}
	pgs := []laws.ParameterGroup{}
	pgs = append(pgs, pg)
	r := encodeParameterGroups(pgs)
	var result []laws.ParameterGroup
	dec := gob.NewDecoder(&r)
	err := dec.Decode(&result)
	if err != nil {
		t.Errorf("decode error: %s", err)
	}
	if len(pgs) != len(result) {
		t.Errorf("got %d expected %d", len(result), len(pgs))
	}
}

func Test_DecodeParameterGroups(t *testing.T) {
	pg := laws.ParameterGroup{}
	pgs := []laws.ParameterGroup{}
	pgs = append(pgs, pg)
	r := encodeParameterGroups(pgs)
	result := DecodeParameterGroups(r)
	if len(pgs) != len(result) {
		t.Errorf("got %d expected %d", len(result), len(pgs))
	}
}
