package aws

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/jrottersman/lats/state"
)

// Client is used for mocking the AWS RDS instance for testing
type Client interface {
	DescribeDBClusters(ctx context.Context, params *rds.DescribeDBClustersInput, optFns ...func(*rds.Options)) (*rds.DescribeDBClustersOutput, error)
	DescribeDBInstances(ctx context.Context, input *rds.DescribeDBInstancesInput, optFns ...func(*rds.Options)) (*rds.DescribeDBInstancesOutput, error)
	CreateDBSnapshot(ctx context.Context, params *rds.CreateDBSnapshotInput, optFns ...func(*rds.Options)) (*rds.CreateDBSnapshotOutput, error)
	DescribeDBParameterGroups(ctx context.Context, params *rds.DescribeDBParameterGroupsInput, optFns ...func(*rds.Options)) (*rds.DescribeDBParameterGroupsOutput, error)
	CopyDBSnapshot(ctx context.Context, params *rds.CopyDBSnapshotInput, optFns ...func(*rds.Options)) (*rds.CopyDBSnapshotOutput, error)
	RestoreDBClusterFromSnapshot(ctx context.Context, params *rds.RestoreDBClusterFromSnapshotInput, optFns ...func(*rds.Options)) (*rds.RestoreDBClusterFromSnapshotOutput, error)
	RestoreDBInstanceFromDBSnapshot(ctx context.Context, params *rds.RestoreDBInstanceFromDBSnapshotInput, optFns ...func(*rds.Options)) (*rds.RestoreDBInstanceFromDBSnapshotOutput, error)
}

// DbInstances holds our RDS client that allows for operations in AWS
type DbInstances struct {
	RdsClient Client
}

// GetInstance describes an RDS instance and returns it's output
func (instances *DbInstances) GetInstance(instanceName string) (
	*types.DBInstance, error) {
	output, err := instances.RdsClient.DescribeDBInstances(context.TODO(),
		&rds.DescribeDBInstancesInput{
			DBInstanceIdentifier: aws.String(instanceName),
		})
	if err != nil {
		var notFoundError *types.DBInstanceNotFoundFault
		if errors.As(err, &notFoundError) {
			log.Printf("DB instance %v does not exist.\n", instanceName)
			err = nil
		} else {
			log.Printf("Couldn't get instance %v: %v\n", instanceName, err)
		}
		return nil, err
	}
	return &output.DBInstances[0], nil
}

func (instances *DbInstances) GetCluster(clusterName string) (*types.DBCluster, error) {
	output, err := instances.RdsClient.DescribeDBClusters(context.TODO(),
		&rds.DescribeDBClustersInput{
			DBClusterIdentifier: aws.String(clusterName),
		})
	if err != nil {
		var notFoundError *types.DBClusterNotFoundFault
		if errors.As(err, &notFoundError) {
			log.Printf("DB cluster %v does not exist.\n", clusterName)
			err = nil
		} else {
			log.Printf("Couldn't get DB cluster %v: %v\n", clusterName, err)
		}
		return nil, err
	}
	return &output.DBClusters[0], err
}

// CreateSnapshot cretaes an AWS snapshot
// :instanceName - name of the database we want to backup
// :snapShotName name of the backup we are creating
func (instances *DbInstances) CreateSnapshot(instanceName string, snapshotName string) (
	*types.DBSnapshot, error) {
	output, err := instances.RdsClient.CreateDBSnapshot(context.TODO(), &rds.CreateDBSnapshotInput{
		DBInstanceIdentifier: aws.String(instanceName),
		DBSnapshotIdentifier: aws.String(snapshotName),
	})
	if err != nil {
		log.Printf("Couldn't create snapshot %v: %v\n", snapshotName, err)
		return nil, err
	}
	return output.DBSnapshot, nil
}

// CopySnapshot copies a snapshot to a new region note it needs to run from the destination region so it needs a different client then CreateSnapshot!
func (instances *DbInstances) CopySnapshot(originalSnapshotName string, NewSnapshotName string, sourceRegion string, KmsKey string) (
	*types.DBSnapshot, error) {
	output, err := instances.RdsClient.CopyDBSnapshot(context.TODO(), &rds.CopyDBSnapshotInput{
		SourceDBSnapshotIdentifier: aws.String(originalSnapshotName),
		TargetDBSnapshotIdentifier: aws.String(NewSnapshotName),
		SourceRegion:               aws.String(sourceRegion), // this generates a presigned URL under the hood which enables cross region copies
		KmsKeyId:                   aws.String(KmsKey),
	})
	if err != nil {
		log.Printf("Couldn't copy snapshot %s: %s\n", NewSnapshotName, err)
		return nil, err
	}
	return output.DBSnapshot, nil
}

// TODO Copy Option Group

// GetParameterGroup we will use this for moving custom parameter groups around jumped the gun here but oh well
func (instances *DbInstances) GetParameterGroup(parameterGroupName string) (
	*types.DBParameterGroup, error) {
	output, err := instances.RdsClient.DescribeDBParameterGroups(
		context.TODO(), &rds.DescribeDBParameterGroupsInput{
			DBParameterGroupName: aws.String(parameterGroupName),
		})
	if err != nil {
		var notFoundError *types.DBParameterGroupNotFoundFault
		if errors.As(err, &notFoundError) {
			log.Printf("Parameter group %v does not exist.\n", parameterGroupName)
			err = nil
		} else {
			log.Printf("Error getting parameter group %v: %v\n", parameterGroupName, err)
		}
		return nil, err
	}
	return &output.DBParameterGroups[0], err

}

func (instances *DbInstances) restoreSnapshotCluster(store state.RDSRestorationStore) (*rds.RestoreDBClusterFromSnapshotOutput, error) {
	backupClusterIden := fmt.Sprintf("%s-backup", *store.Cluster.DBClusterIdentifier)
	input := rds.RestoreDBClusterFromSnapshotInput{
		DBClusterIdentifier: aws.String(backupClusterIden),
		SnapshotIdentifier:  store.GetClusterSnapshotIdentifier(),
		Engine:              store.GetClusterEngine(),
		AvailabilityZones:   *store.GetClusterAZs(),
	}
	output, err := instances.RdsClient.RestoreDBClusterFromSnapshot(context.TODO(), &input)
	if err != nil {
		log.Printf("error creating snapshot cluster")
		return nil, err
	}
	return output, nil
}

// RestoreSnapshotInstance restores a single db instance from a snapshot
func (instances *DbInstances) RestoreSnapshotInstance(store state.RDSRestorationStore) (*rds.RestoreDBInstanceFromDBSnapshotOutput, error) {

	backupDbIden := fmt.Sprintf("%s-backup", *store.GetInstanceIdentifier())
	input := rds.RestoreDBInstanceFromDBSnapshotInput{
		DBInstanceIdentifier: &backupDbIden,
		DBSnapshotIdentifier: store.GetSnapshotIdentifier(),
		AllocatedStorage:     store.GetAllocatedStorage(),
	}

	output, err := instances.RdsClient.RestoreDBInstanceFromDBSnapshot(context.TODO(), &input)
	if err != nil {
		log.Printf("error creating instance from snapshot %s", err)
		return nil, err
	}
	return output, nil
}
