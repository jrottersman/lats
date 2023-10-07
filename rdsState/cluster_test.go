package rdsstate_test

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/jrottersman/lats/aws"
	mock "github.com/jrottersman/lats/mocks"
	"github.com/jrottersman/lats/rdsstate"
	"github.com/jrottersman/lats/state"
)

func TestGenerateRDSClusterStack(t *testing.T) {
	type args struct {
		i rdsstate.ClusterStackInput
	}

	clusterId := "foo"
	snapshotId := "bar"
	filen := "/tmp/bar"
	pFileName := "/tmp/foo"
	c := mock.MockRDSClient{}
	i := aws.DbInstances{RdsClient: c}

	resto := state.RDSRestorationStore{
		Cluster:         &types.DBCluster{DBClusterIdentifier: &clusterId},
		ClusterSnapshot: &types.DBClusterSnapshot{DBClusterSnapshotIdentifier: &snapshotId},
	}
	input := rdsstate.ClusterStackInput{
		R:                 resto,
		StackName:         "foo",
		Filename:          filen,
		ParameterFileName: pFileName,
		ParameterGroups:   []aws.ParameterGroup{},
		Client:            i,
		Folder:            "/tmp",
	}
	arg := args{
		i: input,
	}

	objs := make(map[int][]state.Object)
	pObjs := []state.Object{}
	pObj := state.Object{
		FileName: pFileName,
		Order:    1,
		ObjType:  state.DBClusterParameterGroup,
	}
	pObjs = append(pObjs, pObj)
	tObjs := []state.Object{}
	obj := state.Object{
		FileName: "/tmp/bar",
		Order:    2,
		ObjType:  state.Cluster,
	}
	tObjs = append(tObjs, obj)
	objs[1] = pObjs
	objs[2] = tObjs
	objs[3] = nil
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
			got, err := rdsstate.GenerateRDSClusterStack(tt.args.i)
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
