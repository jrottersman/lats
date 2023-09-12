package aws

import (
	"context"
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
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

func (m mockRDSClient) DescribeDBInstances(ctx context.Context, input *rds.DescribeDBInstancesInput, optFns ...func(*rds.Options)) (*rds.DescribeDBInstancesOutput, error) {
	f := "foo"
	r := &rds.DescribeDBInstancesOutput{
		DBInstances: []types.DBInstance{{DBInstanceIdentifier: &f}},
	}
	return r, nil
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

func TestGetCluster(t *testing.T) {
	expected := "foo"
	c := mockRDSClient{}
	dbi := DbInstances{
		RdsClient: c,
	}
	resp, err := dbi.GetCluster(expected)
	if err != nil {
		t.Errorf("got error %s", err)
	}
	if *resp.DBClusterIdentifier != expected {
		t.Errorf("got %s expected %s", *resp.DBClusterIdentifier, expected)
	}
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

func TestGetInstanceFromCluster(t *testing.T) {
	expected := "foo"
	c := mockRDSClient{}
	dbi := DbInstances{
		RdsClient: c,
	}
	resp, err := dbi.GetCluster(expected)
	if err != nil {
		t.Errorf("got error %s", err)
	}
	r, err := dbi.GetInstancesFromCluster(resp)
	if err != nil {
		t.Errorf("got error %s", err)
	}
	if *r[0].DBInstanceIdentifier != expected {
		t.Errorf("got %s expected %s", *r[0].DBInstanceIdentifier, expected)
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

func TestCreateClusterSnapshot(t *testing.T) {
	c := mockRDSClient{}
	dbi := DbInstances{
		RdsClient: c,
	}
	resp, err := dbi.CreateClusterSnapshot("foo", "bar")
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

func TestCopySnapshot(t *testing.T) {
	c := mockRDSClient{}
	dbi := DbInstances{
		RdsClient: c,
	}
	resp, err := dbi.CopySnapshot("foo", "bar", "us-east-1", "keyArn")
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if resp.AllocatedStorage != 1000 {
		t.Errorf("got %d expected 1000", resp.AllocatedStorage)
	}
}

func TestRestoreSnapshotInstance(t *testing.T) {

	filename := "/tmp/foo"
	snap := types.DBSnapshot{
		AllocatedStorage:     1000,
		Encrypted:            true,
		PercentProgress:      100,
		DBInstanceIdentifier: aws.String("foobar"),
		DBSnapshotIdentifier: aws.String("foo"),
	}

	defer os.Remove(filename)
	r := state.EncodeRDSSnapshotOutput(&snap)
	_, err := state.WriteOutput(filename, r)
	if err != nil {
		t.Errorf("error writing file %s", err)
	}

	var s []state.StateKV
	kv := state.StateKV{
		Object:       "foo",
		FileLocation: filename,
		ObjectType:   "RDSSnapshot",
	}
	s = append(s, kv)

	filename2 := "/tmp/foo2"
	dbz := types.DBInstance{
		AllocatedStorage:     1000,
		DBInstanceIdentifier: aws.String("foobar"),
	}

	defer os.Remove(filename2)
	r2 := state.EncodeRDSDatabaseOutput(&dbz)
	_, err = state.WriteOutput(filename2, r2)
	if err != nil {
		t.Errorf("error writing file %s", err)
	}
	kv2 := state.StateKV{
		Object:       "foobar",
		FileLocation: filename2,
		ObjectType:   state.RdsInstanceType,
	}
	s = append(s, kv2)
	sm := state.StateManager{
		Mu:             &sync.Mutex{},
		StateLocations: s,
	}

	resp, err := state.RDSRestorationStoreBuilder(sm, "foo")
	if err != nil {
		t.Errorf("error writing file %s", err)
	}

	c := mockRDSClient{}
	dbi := DbInstances{
		RdsClient: c,
	}
	input := state.GenerateRestoreDBInstanceFromDBSnapshotInput(*resp)
	resp2, err := dbi.RestoreSnapshotInstance(*input)
	if err != nil {
		t.Errorf("got error: %s", err)
	}
	if reflect.TypeOf(resp2) != reflect.TypeOf(&rds.RestoreDBInstanceFromDBSnapshotOutput{}) {
		t.Error()
	}
}

func TestRestoreSnapshotCluster(t *testing.T) {

	snap := types.DBClusterSnapshot{
		AllocatedStorage:            1000,
		PercentProgress:             100,
		DBClusterIdentifier:         aws.String("foobar"),
		DBClusterSnapshotIdentifier: aws.String("foo"),
	}

	dbz := types.DBCluster{
		DBClusterIdentifier: aws.String("foobar"),
	}

	store := state.RDSRestorationStore{
		Cluster:         &dbz,
		ClusterSnapshot: &snap,
	}

	c := mockRDSClient{}
	dbi := DbInstances{
		RdsClient: c,
	}
	input := state.GenerateRestoreDBInstanceFromDBClusterSnapshotInput(store)
	resp, err := dbi.RestoreSnapshotInstance(*input)
	if err != nil {
		t.Errorf("got error: %s", err)
	}
	if reflect.TypeOf(resp) != reflect.TypeOf(&rds.RestoreDBInstanceFromDBSnapshotOutput{}) {
		t.Errorf("got %T expected %T", resp, &rds.RestoreDBInstanceFromDBSnapshotOutput{})
	}
}
