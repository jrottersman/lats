package mock

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
)

// MockRDSClient is a struct to mock an RDS client testing is fun :)
type MockRDSClient struct{}

// DescribeDBClusters mock get a db cluster
func (m MockRDSClient) DescribeDBClusters(ctx context.Context, params *rds.DescribeDBClustersInput, optFns ...func(*rds.Options)) (*rds.DescribeDBClustersOutput, error) {
	f := "foo"
	r := &rds.DescribeDBClustersOutput{
		DBClusters: []types.DBCluster{{DBClusterIdentifier: &f, Status: aws.String("creating"), DBClusterMembers: []types.DBClusterMember{{DBInstanceIdentifier: &f}}}},
	}
	return r, nil
}

// DescribeDBClusterSnapshots mock get a db cluster snapshot
func (m MockRDSClient) DescribeDBClusterSnapshots(ctx context.Context, params *rds.DescribeDBClusterSnapshotsInput, optFns ...func(*rds.Options)) (*rds.DescribeDBClusterSnapshotsOutput, error) {
	snapshots := []types.DBClusterSnapshot{}
	snap := types.DBClusterSnapshot{
		DBClusterSnapshotIdentifier: aws.String("foo"),
		DBClusterSnapshotArn:        aws.String("foo"),
		Status:                      aws.String("Completed"),
	}
	snapshots = append(snapshots, snap)

	return &rds.DescribeDBClusterSnapshotsOutput{
		DBClusterSnapshots: snapshots,
	}, nil
}

// DescribeDBInstances mock get for a db instance
func (m MockRDSClient) DescribeDBInstances(ctx context.Context, input *rds.DescribeDBInstancesInput, optFns ...func(*rds.Options)) (*rds.DescribeDBInstancesOutput, error) {
	f := "foo"
	status := "availaible"
	r := &rds.DescribeDBInstancesOutput{
		DBInstances: []types.DBInstance{{DBInstanceIdentifier: &f, DBInstanceStatus: &status}},
	}
	return r, nil
}

// DescribeDBSnapshots mock get for a db snapshot
func (m MockRDSClient) DescribeDBSnapshots(ctx context.Context, params *rds.DescribeDBSnapshotsInput, optFns ...func(*rds.Options)) (*rds.DescribeDBSnapshotsOutput, error) {
	snapshots := []types.DBSnapshot{}
	snap := types.DBSnapshot{
		DBSnapshotIdentifier: aws.String("foo"),
		DBSnapshotArn:        aws.String("foo"),
		Status:               aws.String("Completed"),
		PercentProgress:      aws.Int32(100),
	}
	snapshots = append(snapshots, snap)
	return &rds.DescribeDBSnapshotsOutput{
		DBSnapshots: snapshots,
	}, nil
}

// CreateDBInstance mock create a db instance
func (m MockRDSClient) CreateDBInstance(ctx context.Context, params *rds.CreateDBInstanceInput, optFns ...func(*rds.Options)) (*rds.CreateDBInstanceOutput, error) {
	return &rds.CreateDBInstanceOutput{}, nil
}

// CreateDBCluster mock create a db cluster
func (m MockRDSClient) CreateDBSnapshot(ctx context.Context, params *rds.CreateDBSnapshotInput, optFns ...func(*rds.Options)) (*rds.CreateDBSnapshotOutput, error) {
	var store int32
	store = 1000
	r := &rds.CreateDBSnapshotOutput{
		DBSnapshot: &types.DBSnapshot{
			AllocatedStorage: &store,
		},
	}
	return r, nil
}

// CreateDBSubnetGroup mock create a db subnet group
func (m MockRDSClient) CreateDBSubnetGroup(ctx context.Context, params *rds.CreateDBSubnetGroupInput, optFns ...func(*rds.Options)) (*rds.CreateDBSubnetGroupOutput, error) {
	return &rds.CreateDBSubnetGroupOutput{}, nil
}

// CreateDBParameterGroup mock create a db parameter group
func (m MockRDSClient) CreateDBParameterGroup(ctx context.Context, params *rds.CreateDBParameterGroupInput, optFns ...func(*rds.Options)) (*rds.CreateDBParameterGroupOutput, error) {
	parameters := types.DBParameterGroup{
		DBParameterGroupArn:    aws.String("foo"),
		DBParameterGroupFamily: aws.String("foo"),
		DBParameterGroupName:   aws.String("foo"),
	}
	return &rds.CreateDBParameterGroupOutput{DBParameterGroup: &parameters}, nil
}

func (m MockRDSClient) CreateDBClusterParameterGroup(ctx context.Context, params *rds.CreateDBClusterParameterGroupInput, optFns ...func(*rds.Options)) (*rds.CreateDBClusterParameterGroupOutput, error) {
	parameters := types.DBClusterParameterGroup{
		DBClusterParameterGroupArn:  aws.String("foo"),
		DBParameterGroupFamily:      aws.String("foo"),
		DBClusterParameterGroupName: aws.String("foo"),
	}
	return &rds.CreateDBClusterParameterGroupOutput{DBClusterParameterGroup: &parameters}, nil
}

func (m MockRDSClient) CreateOptionGroup(ctx context.Context, params *rds.CreateOptionGroupInput, optFns ...func(*rds.Options)) (*rds.CreateOptionGroupOutput, error) {
	return &rds.CreateOptionGroupOutput{}, nil
}

func (m MockRDSClient) DescribeDBClusterParameterGroups(ctx context.Context, params *rds.DescribeDBClusterParameterGroupsInput, optFns ...func(*rds.Options)) (*rds.DescribeDBClusterParameterGroupsOutput, error) {
	f := "foo"
	return &rds.DescribeDBClusterParameterGroupsOutput{
		DBClusterParameterGroups: []types.DBClusterParameterGroup{{DBClusterParameterGroupName: &f}},
	}, nil
}

func (m MockRDSClient) DescribeOptionGroups(ctx context.Context, params *rds.DescribeOptionGroupsInput, optFns ...func(*rds.Options)) (*rds.DescribeOptionGroupsOutput, error) {
	x := []types.OptionGroup{}
	x = append(x, types.OptionGroup{})
	return &rds.DescribeOptionGroupsOutput{
		OptionGroupsList: x,
	}, nil
}

func (m MockRDSClient) DescribeDBClusterParameters(ctx context.Context, params *rds.DescribeDBClusterParametersInput, optFns ...func(*rds.Options)) (*rds.DescribeDBClusterParametersOutput, error) {
	x := []types.Parameter{}
	return &rds.DescribeDBClusterParametersOutput{Parameters: x}, nil
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
func (m MockRDSClient) ModifyDBParameterGroup(ctx context.Context, params *rds.ModifyDBParameterGroupInput, optFns ...func(*rds.Options)) (*rds.ModifyDBParameterGroupOutput, error) {
	return &rds.ModifyDBParameterGroupOutput{}, nil
}

func (m MockRDSClient) ModifyOptionGroup(ctx context.Context, params *rds.ModifyOptionGroupInput, optFns ...func(*rds.Options)) (*rds.ModifyOptionGroupOutput, error) {
	return &rds.ModifyOptionGroupOutput{}, nil
}

func (m MockRDSClient) ModifyDBClusterParameterGroup(ctx context.Context, params *rds.ModifyDBClusterParameterGroupInput, optFns ...func(*rds.Options)) (*rds.ModifyDBClusterParameterGroupOutput, error) {
	return &rds.ModifyDBClusterParameterGroupOutput{}, nil
}

func (m MockRDSClient) CopyDBClusterSnapshot(ctx context.Context, params *rds.CopyDBClusterSnapshotInput, optFns ...func(*rds.Options)) (*rds.CopyDBClusterSnapshotOutput, error) {
	return &rds.CopyDBClusterSnapshotOutput{
		DBClusterSnapshot: &types.DBClusterSnapshot{},
	}, nil
}

func (m MockRDSClient) CopyDBSnapshot(ctx context.Context, params *rds.CopyDBSnapshotInput, optFns ...func(*rds.Options)) (*rds.CopyDBSnapshotOutput, error) {
	var store int32
	store = 1000
	r := &rds.CopyDBSnapshotOutput{
		DBSnapshot: &types.DBSnapshot{
			AllocatedStorage: &store,
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
	var store int32
	store = 1000
	r := rds.CreateDBClusterSnapshotOutput{
		DBClusterSnapshot: &types.DBClusterSnapshot{
			AllocatedStorage: &store,
		},
	}
	return &r, nil
}
