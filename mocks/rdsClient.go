package mock

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
)

type MockRDSClient struct{}

func (m MockRDSClient) DescribeDBClusters(ctx context.Context, params *rds.DescribeDBClustersInput, optFns ...func(*rds.Options)) (*rds.DescribeDBClustersOutput, error) {
	f := "foo"
	r := &rds.DescribeDBClustersOutput{
		DBClusters: []types.DBCluster{{DBClusterIdentifier: &f, Status: aws.String("creating"), DBClusterMembers: []types.DBClusterMember{{DBInstanceIdentifier: &f}}}},
	}
	return r, nil
}

func (m MockRDSClient) DescribeDBClusterSnapshots(ctx context.Context, params *rds.DescribeDBClusterSnapshotsInput, optFns ...func(*rds.Options)) (*rds.DescribeDBClusterSnapshotsOutput, error) {
	snapshots := []types.DBClusterSnapshot{}
	snap := types.DBClusterSnapshot{
		DBClusterSnapshotIdentifier: aws.String("foo"),
		DBClusterSnapshotArn:        aws.String("foo"),
	}
	snapshots = append(snapshots, snap)

	return &rds.DescribeDBClusterSnapshotsOutput{
		DBClusterSnapshots: snapshots,
	}, nil
}

func (m MockRDSClient) DescribeDBInstances(ctx context.Context, input *rds.DescribeDBInstancesInput, optFns ...func(*rds.Options)) (*rds.DescribeDBInstancesOutput, error) {
	f := "foo"
	r := &rds.DescribeDBInstancesOutput{
		DBInstances: []types.DBInstance{{DBInstanceIdentifier: &f}},
	}
	return r, nil
}

func (m MockRDSClient) DescribeDBSnapshots(ctx context.Context, params *rds.DescribeDBSnapshotsInput, optFns ...func(*rds.Options)) (*rds.DescribeDBSnapshotsOutput, error) {
	snapshots := []types.DBSnapshot{}
	snap := types.DBSnapshot{
		DBSnapshotIdentifier: aws.String("foo"),
		DBSnapshotArn:        aws.String("foo"),
	}
	snapshots = append(snapshots, snap)
	return &rds.DescribeDBSnapshotsOutput{
		DBSnapshots: snapshots,
	}, nil
}

func (m MockRDSClient) CreateDBInstance(ctx context.Context, params *rds.CreateDBInstanceInput, optFns ...func(*rds.Options)) (*rds.CreateDBInstanceOutput, error) {
	return &rds.CreateDBInstanceOutput{}, nil
}

func (m MockRDSClient) CreateDBSnapshot(ctx context.Context, params *rds.CreateDBSnapshotInput, optFns ...func(*rds.Options)) (*rds.CreateDBSnapshotOutput, error) {
	r := &rds.CreateDBSnapshotOutput{
		DBSnapshot: &types.DBSnapshot{
			AllocatedStorage: 1000,
		},
	}
	return r, nil
}

func (m MockRDSClient) DescribeDBParameterGroups(ctx context.Context, params *rds.DescribeDBParameterGroupsInput, optFns ...func(*rds.Options)) (*rds.DescribeDBParameterGroupsOutput, error) {
	f := "foo"
	r := rds.DescribeDBParameterGroupsOutput{
		DBParameterGroups: []types.DBParameterGroup{{DBParameterGroupName: &f}},
	}
	return &r, nil
}

func (m MockRDSClient) DescribeDBParameters(ctx context.Context, params *rds.DescribeDBParametersInput, optFns ...func(*rds.Options)) (*rds.DescribeDBParametersOutput, error) {
	x := []types.Parameter{}
	return &rds.DescribeDBParametersOutput{Parameters: x}, nil
}

func (m MockRDSClient) CopyDBClusterSnapshot(ctx context.Context, params *rds.CopyDBClusterSnapshotInput, optFns ...func(*rds.Options)) (*rds.CopyDBClusterSnapshotOutput, error) {
	return &rds.CopyDBClusterSnapshotOutput{
		DBClusterSnapshot: &types.DBClusterSnapshot{},
	}, nil
}

func (m MockRDSClient) CopyDBSnapshot(ctx context.Context, params *rds.CopyDBSnapshotInput, optFns ...func(*rds.Options)) (*rds.CopyDBSnapshotOutput, error) {
	r := &rds.CopyDBSnapshotOutput{
		DBSnapshot: &types.DBSnapshot{
			AllocatedStorage: 1000,
		},
	}
	return r, nil
}

func (m MockRDSClient) RestoreDBClusterFromSnapshot(ctx context.Context, params *rds.RestoreDBClusterFromSnapshotInput, optFns ...func(*rds.Options)) (*rds.RestoreDBClusterFromSnapshotOutput, error) {
	r := &rds.RestoreDBClusterFromSnapshotOutput{DBCluster: &types.DBCluster{Status: aws.String("creating")}}
	return r, nil
}

func (m MockRDSClient) RestoreDBInstanceFromDBSnapshot(ctx context.Context, params *rds.RestoreDBInstanceFromDBSnapshotInput, optFns ...func(*rds.Options)) (*rds.RestoreDBInstanceFromDBSnapshotOutput, error) {
	r := &rds.RestoreDBInstanceFromDBSnapshotOutput{}
	return r, nil
}

func (m MockRDSClient) CreateDBClusterSnapshot(ctx context.Context, params *rds.CreateDBClusterSnapshotInput, optFns ...func(*rds.Options)) (*rds.CreateDBClusterSnapshotOutput, error) {
	r := rds.CreateDBClusterSnapshotOutput{
		DBClusterSnapshot: &types.DBClusterSnapshot{
			AllocatedStorage: 1000,
		},
	}
	return &r, nil
}
