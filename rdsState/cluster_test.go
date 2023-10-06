package rdsState_test

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/jrottersman/lats/aws"
	mock "github.com/jrottersman/lats/mocks"
	"github.com/jrottersman/lats/rdsState"
	"github.com/jrottersman/lats/state"
)

func TestGenerateRDSClusterStack(t *testing.T) {
	type args struct {
		i rdsState.ClusterStackInput
	}

	clusterId := "foo"
	snapshotId := "bar"
	filen := "/tmp/bar"
	c := mock.MockRDSClient{}
	i := aws.DbInstances{RdsClient: c}

	resto := state.RDSRestorationStore{
		Cluster:         &types.DBCluster{DBClusterIdentifier: &clusterId},
		ClusterSnapshot: &types.DBClusterSnapshot{DBClusterSnapshotIdentifier: &snapshotId},
	}
	input := rdsState.ClusterStackInput{
		R:         resto,
		StackName: "foo",
		Filename:  filen,
		Client:    i,
		Folder:    "/tmp",
	}
	arg := args{
		i: input,
	}

	objs := make(map[int][]state.Object)
	tObjs := []state.Object{}
	obj := state.Object{
		FileName: "/tmp/bar",
		Order:    1,
		ObjType:  state.Cluster,
	}
	tObjs = append(tObjs, obj)
	objs[1] = tObjs
	objs[2] = nil
	wanted := state.Stack{
		Name:                  "foo",
		RestorationObjectName: state.Cluster,
		Objects:               objs,
	}

	tests := []struct {
		name    string
		args    args
		want    *state.Stack
		wantErr bool
	}{
		{"test", arg, &wanted, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rdsState.GenerateRDSClusterStack(tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateRDSClusterStack() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateRDSClusterStack() = %v, want %v", got, tt.want)
			}
		})
	}
}
