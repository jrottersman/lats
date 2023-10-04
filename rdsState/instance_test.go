package rdsState

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
		r       state.RDSRestorationStore
		name    string
		fn      *string
		paramfn *string
		pg      []laws.ParameterGroup
	}
	r := state.RDSRestorationStore{
		Snapshot: &types.DBSnapshot{DBSnapshotIdentifier: aws.String("boo")},
		Instance: &types.DBInstance{},
	}
	arg := args{
		r:       r,
		name:    "bar",
		fn:      aws.String("/tmp/foo.gob"),
		paramfn: aws.String("/tmp/bar.gob"),
		pg:      []laws.ParameterGroup{},
	}
	defer os.Remove("/tmp/foo.gob")
	defer os.Remove("/tmp/bar.gob")

	obj := state.Object{
		FileName: "/tmp/foo.gob",
		Order:    1,
		ObjType:  state.LoneInstance,
	}
	oMap := make(map[int][]state.Object)
	oMap[1] = []state.Object{obj}
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
			got, err := GenerateRDSInstanceStack(tt.args.r, tt.args.name, tt.args.fn, tt.args.paramfn, tt.args.pg)
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
