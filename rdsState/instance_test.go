package rdsState

import (
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/jrottersman/lats/state"
)

func TestGenerateRDSInstanceStack(t *testing.T) {
	type args struct {
		r    state.RDSRestorationStore
		name string
		fn   *string
	}
	r := state.RDSRestorationStore{
		Snapshot: &types.DBSnapshot{DBSnapshotIdentifier: aws.String("boo")},
		Instance: &types.DBInstance{},
	}
	arg := args{
		r:    r,
		name: "bar",
		fn:   aws.String("/tmp/foo.gob"),
	}
	defer os.Remove("/tmp/foo.gob")

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
			got, err := GenerateRDSInstanceStack(tt.args.r, tt.args.name, tt.args.fn)
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
