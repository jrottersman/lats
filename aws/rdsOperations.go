package aws

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
)

// Client is used for mocking the AWS RDS instance for testing
type Client interface {
	DescribeDBInstances(ctx context.Context, input *rds.DescribeDBInstancesInput, optFns ...func(*rds.Options)) (*rds.DescribeDBInstancesOutput, error)
	CreateDBSnapshot(ctx context.Context, params *rds.CreateDBSnapshotInput, optFns ...func(*rds.Options)) (*rds.CreateDBSnapshotOutput, error)
	DescribeDBParameterGroups(ctx context.Context, params *rds.DescribeDBParameterGroupsInput, optFns ...func(*rds.Options)) (*rds.DescribeDBParameterGroupsOutput, error)
	CopyDBSnapshot(ctx context.Context, params *rds.CopyDBSnapshotInput, optFns ...func(*rds.Options)) (*rds.CopyDBSnapshotOutput, error)
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
	} else {
		return &output.DBInstances[0], nil
	}
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

// Get a Parameter Group we will use this for moving custom parameter groups around jumped the gun here but oh well
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
	} else {
		return &output.DBParameterGroups[0], err
	}
}
