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
	DescribeDBClusterSnapshots(ctx context.Context, params *rds.DescribeDBClusterSnapshotsInput, optFns ...func(*rds.Options)) (*rds.DescribeDBClusterSnapshotsOutput, error)
	DescribeDBInstances(ctx context.Context, input *rds.DescribeDBInstancesInput, optFns ...func(*rds.Options)) (*rds.DescribeDBInstancesOutput, error)
	DescribeDBSnapshots(ctx context.Context, params *rds.DescribeDBSnapshotsInput, optFns ...func(*rds.Options)) (*rds.DescribeDBSnapshotsOutput, error)
	CreateDBSnapshot(ctx context.Context, params *rds.CreateDBSnapshotInput, optFns ...func(*rds.Options)) (*rds.CreateDBSnapshotOutput, error)
	CreateDBClusterSnapshot(ctx context.Context, params *rds.CreateDBClusterSnapshotInput, optFns ...func(*rds.Options)) (*rds.CreateDBClusterSnapshotOutput, error)
	DescribeDBParameterGroups(ctx context.Context, params *rds.DescribeDBParameterGroupsInput, optFns ...func(*rds.Options)) (*rds.DescribeDBParameterGroupsOutput, error)
	CopyDBSnapshot(ctx context.Context, params *rds.CopyDBSnapshotInput, optFns ...func(*rds.Options)) (*rds.CopyDBSnapshotOutput, error)
	CopyDBClusterSnapshot(ctx context.Context, params *rds.CopyDBClusterSnapshotInput, optFns ...func(*rds.Options)) (*rds.CopyDBClusterSnapshotOutput, error)
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

func (instances *DbInstances) GetInstancesFromCluster(c *types.DBCluster) ([]types.DBInstance, error) {
	if c.DBClusterMembers == nil {
		return nil, nil
	}
	dbs := []types.DBInstance{}
	for _, v := range c.DBClusterMembers {
		db, err := instances.GetInstance(*v.DBInstanceIdentifier)
		if err != nil {
			log.Printf("error with get instance %s", err)
		}
		dbs = append(dbs, *db)
	}
	return dbs, nil
}

func (instaces *DbInstances) CreateClusterFromStack(s *state.Stack) error {
	// sort keys
	// sorted := s.SortStack()
	// get the one which is the cluster and create it
	first := s.Objects[1]
	if len(first) != 1 {
		return fmt.Errorf("Multiple clusters and there should only be one")
	}
	for _, v := range first {
		b := v.ReadObject()
		dbi := b.(*rds.RestoreDBClusterFromSnapshotInput)
		_, err := instaces.RestoreSnapshotCluster(*dbi) // we might need to do something with the output in which case this changes
		if err != nil {
			return err
		}
	}

	// get two which is the instances create them in parrallel
	return nil
}

func (instances *DbInstances) CreateInstanceFromStack(s *state.Stack) error {
	instance := s.Objects[1]
	if len(instance) != 1 {
		return fmt.Errorf("There should only be a single instance")
	}
	return nil
}

func (instances *DbInstances) getClusterStatus(name string) (*string, error) {
	cluster, err := instances.GetCluster(name)
	if err != nil {
		return nil, err
	}
	return cluster.Status, nil

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

func (instances *DbInstances) CreateClusterSnapshot(clusterName string, snapshotName string) (*types.DBClusterSnapshot, error) {
	output, err := instances.RdsClient.CreateDBClusterSnapshot(context.TODO(), &rds.CreateDBClusterSnapshotInput{
		DBClusterIdentifier:         aws.String(clusterName),
		DBClusterSnapshotIdentifier: aws.String(snapshotName),
	})
	if err != nil {
		log.Printf("Couldn't create snapshot %s: because of %s\n", snapshotName, err)
		return nil, err
	}
	return output.DBClusterSnapshot, nil
}

// CopySnapshot copies a snapshot to a new region note it needs to run from the destination region so it needs a different client then CreateSnapshot!
func (instances *DbInstances) CopySnapshot(originalSnapshotName string, newSnapshotName string, sourceRegion string, KmsKey string) (
	*types.DBSnapshot, error) {
	output, err := instances.RdsClient.CopyDBSnapshot(context.TODO(), &rds.CopyDBSnapshotInput{
		SourceDBSnapshotIdentifier: aws.String(originalSnapshotName),
		TargetDBSnapshotIdentifier: aws.String(newSnapshotName),
		SourceRegion:               aws.String(sourceRegion), // this generates a presigned URL under the hood which enables cross region copies
		KmsKeyId:                   aws.String(KmsKey),
	})
	if err != nil {
		log.Printf("Couldn't copy snapshot %s: %s\n", newSnapshotName, err)
		return nil, err
	}
	return output.DBSnapshot, nil
}

func (instances *DbInstances) CopyClusterSnaphot(originalSnapshotName string, newSnapshotName string, sourceRegion string, kmsKey string) (
	*types.DBClusterSnapshot, error) {
	output, err := instances.RdsClient.CopyDBClusterSnapshot(context.TODO(), &rds.CopyDBClusterSnapshotInput{
		SourceDBClusterSnapshotIdentifier: aws.String(originalSnapshotName),
		TargetDBClusterSnapshotIdentifier: aws.String(newSnapshotName),
		SourceRegion:                      aws.String(sourceRegion),
		KmsKeyId:                          aws.String(kmsKey),
	})
	if err != nil {
		log.Printf("Couldn't copy snapshot %s: %s\n", newSnapshotName, err)
		return nil, err
	}
	return output.DBClusterSnapshot, nil
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

func (instances *DbInstances) RestoreSnapshotCluster(input rds.RestoreDBClusterFromSnapshotInput) (*rds.RestoreDBClusterFromSnapshotOutput, error) {
	output, err := instances.RdsClient.RestoreDBClusterFromSnapshot(context.TODO(), &input)
	if err != nil {
		log.Printf("error creating snapshot cluster")
		return nil, err
	}
	return output, nil
}

// RestoreSnapshotInstance restores a single db instance from a snapshot
func (instances *DbInstances) RestoreSnapshotInstance(input rds.RestoreDBInstanceFromDBSnapshotInput) (*rds.RestoreDBInstanceFromDBSnapshotOutput, error) {

	output, err := instances.RdsClient.RestoreDBInstanceFromDBSnapshot(context.TODO(), &input)
	if err != nil {
		log.Printf("error creating instance from snapshot %s", err)
		return nil, err
	}
	return output, nil
}

func (instances *DbInstances) GetInstanceSnapshotARN(name string, marker *string) (*string, error) {
	output, err := instances.RdsClient.DescribeDBSnapshots(context.TODO(), &rds.DescribeDBSnapshotsInput{
		Marker: marker,
	})
	if err != nil {
		return nil, err
	}
	for _, v := range output.DBSnapshots {
		if *v.DBSnapshotIdentifier == name {
			return v.DBSnapshotArn, nil
		}
	}
	if output.Marker != nil {
		instances.GetInstanceSnapshotARN(name, output.Marker)

	}
	return nil, fmt.Errorf("snapshot not found")
}

func (instances *DbInstances) GetClusterSnapshotARN(name string, marker *string) (*string, error) {
	output, err := instances.RdsClient.DescribeDBClusterSnapshots(context.TODO(), &rds.DescribeDBClusterSnapshotsInput{
		Marker: marker,
	})
	if err != nil {
		return nil, err
	}
	for _, v := range output.DBClusterSnapshots {
		if *v.DBClusterSnapshotIdentifier == name {
			return v.DBClusterSnapshotArn, nil
		}
	}
	if output.Marker != nil {
		instances.GetClusterSnapshotARN(name, output.Marker)
	}
	return nil, fmt.Errorf("cluster snapshot not found")
}

func (instances *DbInstances) GetSnapshotARN(name string, cluster bool) (*string, error) {
	if cluster {
		snap, err := instances.GetClusterSnapshotARN(name, nil)
		return snap, err
	}
	snap, err := instances.GetInstanceSnapshotARN(name, nil)
	return snap, err
}
