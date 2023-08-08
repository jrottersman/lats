package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
)

type mockRDSClient struct{}

func (m mockRDSClient) DescribeDBInstances(ctx context.Context, input *rds.DescribeDBInstancesInput, optFns ...func(*rds.Options)) (*rds.DescribeDBInstancesOutput, error) {
	f := "foo"
	r := &rds.DescribeDBInstancesOutput{
		DBInstances: []types.DBInstance{types.DBInstance{DBInstanceIdentifier: &f}},
	}
	return r, nil
}

// TODO create mock CreateSnapshot function
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
		DBParameterGroups: []types.DBParameterGroup{types.DBParameterGroup{DBParameterGroupName: &f}},
	}
	return &r, nil
}

func (m mockRDSClient) CopyDBSnapshot(ctx context.Context, params *rds.CopyDBSnapshotInput, optFns ...func(*rds.Options)) (*rds.CopyDBSnapshotOutput, error) {
	return nil, nil
}

func TestGetInstance(t *testing.T) {
	expected := "foo"
	c := mockRDSClient{}
	dbi := DbInstances{
		RdsClient: c,
	}
	resp, err := dbi.GetInstance("foo")
	if err != nil {
		t.Errorf("got the following error %v", err)
	}
	if *resp.DBInstanceIdentifier != expected {
		t.Errorf("got %s expected %s", *resp.DBInstanceIdentifier, expected)
	}
}

func TestCreateSnapshot(t *testing.T) {
	c := mockRDSClient{}
	dbi := DbInstances{
		RdsClient: c,
	}
	resp, err := dbi.CreateSnapshot("foo", "bar")
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if resp.AllocatedStorage != 1000 {
		t.Errorf("got %d expected 1000", resp.AllocatedStorage)
	}
}

func TestDescribeParameterGroup(t *testing.T) {
	c := mockRDSClient{}
	dbi := DbInstances{
		RdsClient: c,
	}
	resp, err := dbi.GetParameterGroup("foo")
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if *resp.DBParameterGroupName != "foo" {
		t.Errorf("expected foo got %s", *resp.DBParameterGroupName)
	}
}
