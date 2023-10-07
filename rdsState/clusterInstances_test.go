package rdsstate_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/jrottersman/lats/aws"
	mock "github.com/jrottersman/lats/mocks"
	"github.com/jrottersman/lats/rdsstate"
	"github.com/jrottersman/lats/state"
)

func TestClusterInstancesToObjects(t *testing.T) {
	type args struct {
		t     *types.DBCluster
		c     aws.DbInstances
		f     string
		order int
	}

	// mock Client
	m := mock.MockRDSClient{}
	cl := aws.DbInstances{RdsClient: m}
	nilArg := args{
		t:     &types.DBCluster{},
		c:     cl,
		f:     "/tmp/foo",
		order: 2,
	}
	id := "foo"
	defer os.Remove("/tmp/foo.gob")
	// Create DB Cluster Member
	mem := []types.DBClusterMember{}
	one := types.DBClusterMember{
		DBInstanceIdentifier: &id,
	}
	mem = append(mem, one)

	//Want object
	objs := []state.Object{}
	fo := state.Object{
		FileName: "/tmp/foo.gob",
		Order:    2,
		ObjType:  state.RdsInstanceType,
	}
	objs = append(objs, fo)
	arg := args{
		t:     &types.DBCluster{DBClusterIdentifier: &id, DBClusterMembers: mem},
		c:     cl,
		f:     "/tmp",
		order: 2,
	}
	tests := []struct {
		name    string
		args    args
		want    []state.Object
		wantErr bool
	}{
		{name: "nil", args: nilArg, want: nil, wantErr: false},
		{name: "one", args: arg, want: objs, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rdsstate.ClusterInstancesToObjects(tt.args.t, tt.args.c, tt.args.f, tt.args.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClusterInstancesToObjects() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClusterInstancesToObjects() = %v, want %v", got, tt.want)
			}
		})
	}
}
