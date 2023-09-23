package rdsState_test

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/jrottersman/lats/aws"
	"github.com/jrottersman/lats/rdsState"
	"github.com/jrottersman/lats/state"
)

type mockRDSClient struct{}

func (m mockRDSClient) DescribeDBClusters(ctx context.Context, params *rds.DescribeDBClustersInput, optFns ...func(*rds.Options)) (*rds.DescribeDBClustersOutput, error) {
	f := "foo"
	r := &rds.DescribeDBClustersOutput{
		DBClusters: []types.DBCluster{{DBClusterIdentifier: &f, DBClusterMembers: []types.DBClusterMember{{DBInstanceIdentifier: &f}}}},
	}
	return r, nil
}

func (m mockRDSClient) DescribeDBClusterSnapshots(ctx context.Context, params *rds.DescribeDBClusterSnapshotsInput, optFns ...func(*rds.Options)) (*rds.DescribeDBClusterSnapshotsOutput, error) {
	return &rds.DescribeDBClusterSnapshotsOutput{}, nil
}

func (m mockRDSClient) DescribeDBInstances(ctx context.Context, input *rds.DescribeDBInstancesInput, optFns ...func(*rds.Options)) (*rds.DescribeDBInstancesOutput, error) {
	f := "foo"
	r := &rds.DescribeDBInstancesOutput{
		DBInstances: []types.DBInstance{{DBInstanceIdentifier: &f}},
	}
	return r, nil
}

func (m mockRDSClient) DescribeDBSnapshots(ctx context.Context, params *rds.DescribeDBSnapshotsInput, optFns ...func(*rds.Options)) (*rds.DescribeDBSnapshotsOutput, error) {
	return &rds.DescribeDBSnapshotsOutput{}, nil
}

func (m mockRDSClient) CreateDBInstance(ctx context.Context, params *rds.CreateDBInstanceInput, optFns ...func(*rds.Options)) (*rds.CreateDBInstanceOutput, error) {
	return &rds.CreateDBInstanceOutput{}, nil
}

func (m mockRDSClient) CreateDBSnapshot(ctx context.Context, params *rds.CreateDBSnapshotInput, optFns ...func(*rds.Options)) (*rds.CreateDBSnapshotOutput, error) {
	r := &rds.CreateDBSnapshotOutput{
		DBSnapshot: &types.DBSnapshot{
			AllocatedStorage: 1000,
		},
	}
	return r, nil
}

func (m mockRDSClient) DescribeDBParameterGroups(ctx context.Context, params *rds.DescribeDBParameterGroupsInput, optFns ...func(*rds.Options)) (*rds.DescribeDBParameterGroupsOutput, error) {
	f := "foo"
	r := rds.DescribeDBParameterGroupsOutput{
		DBParameterGroups: []types.DBParameterGroup{{DBParameterGroupName: &f}},
	}
	return &r, nil
}

func (m mockRDSClient) CopyDBClusterSnapshot(ctx context.Context, params *rds.CopyDBClusterSnapshotInput, optFns ...func(*rds.Options)) (*rds.CopyDBClusterSnapshotOutput, error) {
	return &rds.CopyDBClusterSnapshotOutput{}, nil
}

func (m mockRDSClient) CopyDBSnapshot(ctx context.Context, params *rds.CopyDBSnapshotInput, optFns ...func(*rds.Options)) (*rds.CopyDBSnapshotOutput, error) {
	r := &rds.CopyDBSnapshotOutput{
		DBSnapshot: &types.DBSnapshot{
			AllocatedStorage: 1000,
		},
	}
	return r, nil
}

func (m mockRDSClient) RestoreDBClusterFromSnapshot(ctx context.Context, params *rds.RestoreDBClusterFromSnapshotInput, optFns ...func(*rds.Options)) (*rds.RestoreDBClusterFromSnapshotOutput, error) {
	r := &rds.RestoreDBClusterFromSnapshotOutput{}
	return r, nil
}

func (m mockRDSClient) RestoreDBInstanceFromDBSnapshot(ctx context.Context, params *rds.RestoreDBInstanceFromDBSnapshotInput, optFns ...func(*rds.Options)) (*rds.RestoreDBInstanceFromDBSnapshotOutput, error) {
	r := &rds.RestoreDBInstanceFromDBSnapshotOutput{}
	return r, nil
}

func (m mockRDSClient) CreateDBClusterSnapshot(ctx context.Context, params *rds.CreateDBClusterSnapshotInput, optFns ...func(*rds.Options)) (*rds.CreateDBClusterSnapshotOutput, error) {
	r := rds.CreateDBClusterSnapshotOutput{
		DBClusterSnapshot: &types.DBClusterSnapshot{
			AllocatedStorage: 1000,
		},
	}
	return &r, nil
}

func TestClusterInstancesToObjects(t *testing.T) {
	type args struct {
		t     *types.DBCluster
		c     aws.DbInstances
		f     string
		order int
	}

	// mock Client
	m := mockRDSClient{}
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
			got, err := rdsState.ClusterInstancesToObjects(tt.args.t, tt.args.c, tt.args.f, tt.args.order)
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
