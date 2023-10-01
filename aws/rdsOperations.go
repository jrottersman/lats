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

// MaxConcurrentJobs max number of operations to hit AWS with at the same time
const MaxConcurrentJobs = 3

// Client is used for mocking the AWS RDS instance for testing
type Client interface {
	DescribeDBClusters(ctx context.Context, params *rds.DescribeDBClustersInput, optFns ...func(*rds.Options)) (*rds.DescribeDBClustersOutput, error)
	DescribeDBClusterSnapshots(ctx context.Context, params *rds.DescribeDBClusterSnapshotsInput, optFns ...func(*rds.Options)) (*rds.DescribeDBClusterSnapshotsOutput, error)
	DescribeDBInstances(ctx context.Context, input *rds.DescribeDBInstancesInput, optFns ...func(*rds.Options)) (*rds.DescribeDBInstancesOutput, error)
	DescribeDBSnapshots(ctx context.Context, params *rds.DescribeDBSnapshotsInput, optFns ...func(*rds.Options)) (*rds.DescribeDBSnapshotsOutput, error)
	CreateDBParameterGroup(ctx context.Context, params *rds.CreateDBParameterGroupInput, optFns ...func(*rds.Options)) (*rds.CreateDBParameterGroupOutput, error)
	CreateDBInstance(ctx context.Context, params *rds.CreateDBInstanceInput, optFns ...func(*rds.Options)) (*rds.CreateDBInstanceOutput, error)
	CreateDBSnapshot(ctx context.Context, params *rds.CreateDBSnapshotInput, optFns ...func(*rds.Options)) (*rds.CreateDBSnapshotOutput, error)
	CreateDBClusterSnapshot(ctx context.Context, params *rds.CreateDBClusterSnapshotInput, optFns ...func(*rds.Options)) (*rds.CreateDBClusterSnapshotOutput, error)
	ModifyDBParameterGroup(ctx context.Context, params *rds.ModifyDBParameterGroupInput, optFns ...func(*rds.Options)) (*rds.ModifyDBParameterGroupOutput, error)
	DescribeDBClusterParameters(ctx context.Context, params *rds.DescribeDBClusterParametersInput, optFns ...func(*rds.Options)) (*rds.DescribeDBClusterParametersOutput, error)
	DescribeDBClusterParameterGroups(ctx context.Context, params *rds.DescribeDBClusterParameterGroupsInput, optFns ...func(*rds.Options)) (*rds.DescribeDBClusterParameterGroupsOutput, error)
	DescribeDBParameters(ctx context.Context, params *rds.DescribeDBParametersInput, optFns ...func(*rds.Options)) (*rds.DescribeDBParametersOutput, error)
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

//GetCluster describes an RDS cluster
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

//GetInstancesFromCluster get's the instaces associated with a database cluster
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

//CreateClusterFromStack creates an RDS cluster from a stack
func (instances *DbInstances) CreateClusterFromStack(s *state.Stack) error {

	// get the one which is the cluster and create it
	first := s.Objects[1]
	if len(first) != 1 {
		return fmt.Errorf("Multiple clusters and there should only be one")
	}
	for _, v := range first {
		b := v.ReadObject()
		dbi := b.(*rds.RestoreDBClusterFromSnapshotInput)
		_, err := instances.RestoreSnapshotCluster(*dbi) // we might need to do something with the output in which case this changes
		if err != nil {
			return err
		}
	}

	// get two which is the instances create them in parrallel
	second := s.Objects[2]
	waitChan := make(chan struct{}, MaxConcurrentJobs)
	for _, i := range second {
		waitChan <- struct{}{}
		go func(inst state.Object) {
			o := inst.ReadObject()
			ins := o.(*rds.CreateDBInstanceInput)
			_, err := instances.RestoreInstanceForCluster(*ins)
			if err != nil {
				fmt.Printf("error creating instance %s", err)
			}
		}(i)
	}
	return nil
}

//CreateInstanceFromStack creates an RDS instance from a stack object
func (instances *DbInstances) CreateInstanceFromStack(s *state.Stack) error {
	instance := s.Objects[1]
	if len(instance) != 1 {
		return fmt.Errorf("There should only be a single instance")
	}
	for _, v := range instance {
		b := v.ReadObject()
		ins := b.(*rds.RestoreDBInstanceFromDBSnapshotInput)
		_, err := instances.RestoreSnapshotInstance(*ins)
		if err != nil {
			return err
		}
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

//CreateClusterSnapshot so it turns out AWS is annoying and makes us create snapshots seperatly for clusters and instaces how fun!
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

//CopyClusterSnaphot see CopySnapshot now for a Cluster
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

//GetClusterParameterGroup get the cluster parameter group so we can make a new one in a new region or you know store it for restoration (actually we won't need to do that cause the data is stored on the snapshot :P)
func (instances *DbInstances) GetClusterParameterGroup(ParameterGroupName string) (
	*types.DBClusterParameterGroup, error) {
	output, err := instances.RdsClient.DescribeDBClusterParameterGroups(context.TODO(), &rds.DescribeDBClusterParameterGroupsInput{
		DBClusterParameterGroupName: aws.String(ParameterGroupName),
	})
	if err != nil {
		var notFoundError *types.DBClusterParameterGroupNotFoundFault
		if errors.As(err, &notFoundError) {
			log.Printf("Parameter group %v does not exist.\n", ParameterGroupName)
			err = nil
		} else {
			log.Printf("Error getting parameter group %v: %v\n", ParameterGroupName, err)
		}
		return nil, err
	}
	return &output.DBClusterParameterGroups[0], nil
}

//GetParametersForClusterParameterGroup returns parameters for a cluster group
func (instances *DbInstances) GetParametersForClusterParameterGroup(ParameterGroupName string) (*[]types.Parameter, error) {
	output, err := instances.RdsClient.DescribeDBClusterParameters(context.TODO(), &rds.DescribeDBClusterParametersInput{
		DBClusterParameterGroupName: aws.String(ParameterGroupName),
	})
	if err != nil {
		log.Printf("Error getting parameters %s", err)
		return nil, err
	}
	parameters := output.Parameters
	for {
		if output.Marker == nil {
			break
		}
		output, err := instances.RdsClient.DescribeDBClusterParameters(context.TODO(), &rds.DescribeDBClusterParametersInput{
			DBClusterParameterGroupName: aws.String(ParameterGroupName),
			Marker:                      output.Marker,
		})
		if err != nil {
			log.Printf("Error getting parameters %s", err)
			return nil, err
		}
		parameters = append(parameters, output.Parameters...)
	}
	return &parameters, nil
}

// GetParameterGroup we will use this for moving custom parameter groups
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

//GetParametersForGroup returns the parameters for a parameter group
func (instances *DbInstances) GetParametersForGroup(ParameterGroupName string) (*[]types.Parameter, error) {
	output, err := instances.RdsClient.DescribeDBParameters(context.TODO(), &rds.DescribeDBParametersInput{
		DBParameterGroupName: aws.String(ParameterGroupName),
	})
	if err != nil {
		log.Printf("Error getting parameters %s", err)
		return nil, err
	}
	parameters := output.Parameters
	for {
		if output.Marker == nil {
			break
		}
		output, err := instances.RdsClient.DescribeDBParameters(context.TODO(), &rds.DescribeDBParametersInput{
			DBParameterGroupName: aws.String(ParameterGroupName),
			Marker:               output.Marker,
		})
		if err != nil {
			log.Printf("Error getting parameters %s", err)
			return nil, err
		}
		parameters = append(parameters, output.Parameters...)

	}
	return &parameters, nil
}

//CreateParameterGroup creates a pararmeter group for a DB instance
func (instances *DbInstances) CreateParameterGroup(p *types.DBParameterGroup) (*rds.CreateDBParameterGroupOutput, error) {
	input := rds.CreateDBParameterGroupInput{
		DBParameterGroupFamily: p.DBParameterGroupFamily,
		DBParameterGroupName:   p.DBParameterGroupName,
		Description:            p.Description,
	}
	output, err := instances.RdsClient.CreateDBParameterGroup(context.TODO(), &input)
	if err != nil {
		log.Printf("error creating parameter group %s", err)
		return output, err
	}
	return output, err
}

//ModifyParameterGroup adds all the parameters to a db parameter group
func (instances *DbInstances) ModifyParameterGroup(pg string, parameters []types.Parameter) error {
	//batch this thing
	batchSize := 20
	batches := make([][]types.Parameter, 0, (len(parameters)+batchSize-1)/batchSize)

	for batchSize < len(parameters) {
		parameters, batches = parameters[batchSize:], append(batches, parameters[0:batchSize:batchSize])
	}
	batches = append(batches, parameters)

	for _, batch := range batches {
		_, err := instances.RdsClient.ModifyDBParameterGroup(context.TODO(), &rds.ModifyDBParameterGroupInput{
			DBParameterGroupName: aws.String(pg),
			Parameters:           batch,
		})
		if err != nil {
			fmt.Printf("error updating parameters %s", err)
		}
	}
	return nil
}

//RestoreSnapshotCluster takes a snapshot turns it into a DB Cluster fun fact the cluster won't be ready from just this there will be no instances
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

//RestoreInstanceForCluster our cluster has no instances by default it need's instances to be usable this makes them exist
func (instances *DbInstances) RestoreInstanceForCluster(input rds.CreateDBInstanceInput) (*rds.CreateDBInstanceOutput, error) {
	output, err := instances.RdsClient.CreateDBInstance(context.TODO(), &input)
	if err != nil {
		log.Printf("error creating instance %s", err)
		return nil, err
	}
	return output, nil
}

//GetInstanceSnapshotARN get the arn for an instance snapshot
func (instances *DbInstances) GetInstanceSnapshotARN(name string, marker *string) (*string, error) {
	fmt.Printf("in instance snapshots\n\n")
	output, err := instances.RdsClient.DescribeDBSnapshots(context.TODO(), &rds.DescribeDBSnapshotsInput{
		// Marker: marker,
		DBSnapshotIdentifier: aws.String(name),
	})
	if err != nil {
		fmt.Printf("error with snapshots %s", err)
		return nil, fmt.Errorf("error retreiving snapshot: %s", err)
	}
	for _, v := range output.DBSnapshots {
		fmt.Printf("%v\n", v)
		if *v.DBSnapshotIdentifier == name {
			return v.DBSnapshotArn, nil
		}
	}
	if output.Marker != nil {
		instances.GetInstanceSnapshotARN(name, output.Marker)

	}
	return nil, fmt.Errorf("snapshot not found")
}

//GetClusterSnapshotARN get's the cluster snapshot arn from snapshot name
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

//GetSnapshotARN get's the snapshot ARN from the snapshot name
func (instances *DbInstances) GetSnapshotARN(name string, cluster bool) (*string, error) {
	fmt.Printf("getting snapshot ARN\n\n")
	if cluster {
		snap, err := instances.GetClusterSnapshotARN(name, nil)
		if err != nil {
			return nil, fmt.Errorf("cluster snapshot error: %s", err)
		}
		return snap, nil
	}
	snap, err := instances.GetInstanceSnapshotARN(name, nil)
	if err != nil {
		return nil, fmt.Errorf("Instance Snapshot error %s", err)
	}
	return snap, nil
}
