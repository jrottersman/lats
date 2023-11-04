package aws

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/jrottersman/lats/pgstate"
	"github.com/jrottersman/lats/stack"
)

// MaxConcurrentJobs max number of operations to hit AWS with at the same time
const MaxConcurrentJobs = 3

// Client is used for mocking the AWS RDS instance for testing
type Client interface {
	DescribeDBClusters(ctx context.Context, params *rds.DescribeDBClustersInput, optFns ...func(*rds.Options)) (*rds.DescribeDBClustersOutput, error)
	DescribeDBClusterSnapshots(ctx context.Context, params *rds.DescribeDBClusterSnapshotsInput, optFns ...func(*rds.Options)) (*rds.DescribeDBClusterSnapshotsOutput, error)
	DescribeDBInstances(ctx context.Context, input *rds.DescribeDBInstancesInput, optFns ...func(*rds.Options)) (*rds.DescribeDBInstancesOutput, error)
	DescribeDBSnapshots(ctx context.Context, params *rds.DescribeDBSnapshotsInput, optFns ...func(*rds.Options)) (*rds.DescribeDBSnapshotsOutput, error)
	DescribeOptionGroups(ctx context.Context, params *rds.DescribeOptionGroupsInput, optFns ...func(*rds.Options)) (*rds.DescribeOptionGroupsOutput, error)
	CreateDBSubnetGroup(ctx context.Context, params *rds.CreateDBSubnetGroupInput, optFns ...func(*rds.Options)) (*rds.CreateDBSubnetGroupOutput, error)
	CreateDBParameterGroup(ctx context.Context, params *rds.CreateDBParameterGroupInput, optFns ...func(*rds.Options)) (*rds.CreateDBParameterGroupOutput, error)
	CreateDBClusterParameterGroup(ctx context.Context, params *rds.CreateDBClusterParameterGroupInput, optFns ...func(*rds.Options)) (*rds.CreateDBClusterParameterGroupOutput, error)
	CreateOptionGroup(ctx context.Context, params *rds.CreateOptionGroupInput, optFns ...func(*rds.Options)) (*rds.CreateOptionGroupOutput, error)
	CreateDBInstance(ctx context.Context, params *rds.CreateDBInstanceInput, optFns ...func(*rds.Options)) (*rds.CreateDBInstanceOutput, error)
	CreateDBSnapshot(ctx context.Context, params *rds.CreateDBSnapshotInput, optFns ...func(*rds.Options)) (*rds.CreateDBSnapshotOutput, error)
	CreateDBClusterSnapshot(ctx context.Context, params *rds.CreateDBClusterSnapshotInput, optFns ...func(*rds.Options)) (*rds.CreateDBClusterSnapshotOutput, error)
	ModifyDBParameterGroup(ctx context.Context, params *rds.ModifyDBParameterGroupInput, optFns ...func(*rds.Options)) (*rds.ModifyDBParameterGroupOutput, error)
	ModifyOptionGroup(ctx context.Context, params *rds.ModifyOptionGroupInput, optFns ...func(*rds.Options)) (*rds.ModifyOptionGroupOutput, error)
	ModifyDBClusterParameterGroup(ctx context.Context, params *rds.ModifyDBClusterParameterGroupInput, optFns ...func(*rds.Options)) (*rds.ModifyDBClusterParameterGroupOutput, error)
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
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	output, err := instances.RdsClient.DescribeDBInstances(ctx,
		&rds.DescribeDBInstancesInput{
			DBInstanceIdentifier: aws.String(instanceName),
		})
	if err != nil {
		var notFoundError *types.DBInstanceNotFoundFault
		if errors.As(err, &notFoundError) {
			slog.Warn("DB instance does not exist.\n", "Instance", instanceName)
			err = nil
		} else {
			slog.Warn("Couldn't get instance", "Instance", instanceName, "error", err)
		}
		return nil, err
	}
	return &output.DBInstances[0], nil
}

//GetCluster describes an RDS cluster
func (instances *DbInstances) GetCluster(clusterName string) (*types.DBCluster, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	output, err := instances.RdsClient.DescribeDBClusters(ctx,
		&rds.DescribeDBClustersInput{
			DBClusterIdentifier: aws.String(clusterName),
		})
	if err != nil {
		var notFoundError *types.DBClusterNotFoundFault
		if errors.As(err, &notFoundError) {
			slog.Info("DB cluster does not exist.", "Cluster", clusterName)
			err = nil
		} else {
			slog.Warn("Couldn't get DB cluster", "Cluster", clusterName, "Error", err)
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
			slog.Warn("error with get instance %s", err)
		}
		dbs = append(dbs, *db)
	}
	return dbs, nil
}

//GetOptionGroup get option group by name
func (instances *DbInstances) GetOptionGroup(OptionGroupName string) (*types.OptionGroup, error) {
	input := rds.DescribeOptionGroupsInput{
		OptionGroupName: &OptionGroupName,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	output, err := instances.RdsClient.DescribeOptionGroups(ctx, &input)
	if err != nil {
		return nil, fmt.Errorf("Error getting Option group %s", err)
	}
	return &output.OptionGroupsList[0], nil
}

// Create a DB subnetGroup
func (instances *DbInstances) CreateDBSubnetGroup(name string, description string, subnets []string) (*rds.CreateDBSubnetGroupOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	group, err := instances.RdsClient.CreateDBSubnetGroup(ctx, &rds.CreateDBSubnetGroupInput{
		DBSubnetGroupDescription: aws.String(description),
		DBSubnetGroupName:        aws.String(name),
		SubnetIds:                subnets,
	})
	return group, err
}

//CreateClusterFromStackInput creates a stack all inputs required
type CreateClusterFromStackInput struct {
	S             *stack.Stack
	DBSubnetGroup *string
}

//CreateClusterFromStack creates an RDS cluster from a stack
func (instances *DbInstances) CreateClusterFromStack(c CreateClusterFromStackInput) error {
	first := c.S.Objects[1]
	var pgName *string
	if len(first) < 1 {
		slog.Info("skipping the parameter set")
	} else {
		for _, p := range first {
			pb := p.ReadObject()
			switch pb.(type) {
			case *pgstate.ParameterGroup:
				pg := pb.(*pgstate.ParameterGroup)
				pgName = pg.ClusterParameterGroup.DBClusterParameterGroupName
				slog.Info("creating parameter group", "name", *pgName)
				_, err := instances.CreateClusterParameterGroup(&pg.ClusterParameterGroup)
				if err != nil {
					return err
				}
				batchSize := 20
				params := pg.Params
				batches := make([][]types.Parameter, 0, (len(params)+batchSize-1)/batchSize)
				for batchSize < len(params) {
					params, batches = params[batchSize:], append(batches, params[0:batchSize:batchSize])
				}
				batches = append(batches, params)
				for _, b := range batches {
					err = instances.ModifyClusterParameterGroup(*pg.ClusterParameterGroup.DBClusterParameterGroupName, b)
					if err != nil {
						return err
					}
				}
			case *types.OptionGroup:
				og := pb.(*types.OptionGroup)
				slog.Info("restoring option group")
				_, err := instances.RestoreOptionGroup(*og.EngineName, *og.MajorEngineVersion, *og.OptionGroupName, *og.OptionGroupDescription)
				if err != nil {
					return fmt.Errorf("error creating option group %s", err)
				}
				optConfigs := optionsToConfiguration(og.Options)
				err = instances.ModifyOptionGroup(*og.OptionGroupName, optConfigs)
				if err != nil {
					slog.Warn("error modifying option group", "error", err)
				}
			}
		}
		//Wait five minutes for parameter sets per aws docs
		for i := 0; i < 10; i++ {
			slog.Info("waiting for five minutes for Parameter group per AWS documentation", "seconds", 30*i)
			time.Sleep(30 * time.Second)
		}
	}

	// get the one which is the cluster and create it
	second := c.S.Objects[2]
	if len(second) != 1 {
		slog.Error("Multiple clusters and there should only be one")
		return fmt.Errorf("Multiple clusters and there should only be one")
	}
	for _, v := range second {
		b := v.ReadObject()
		dbi := b.(*rds.RestoreDBClusterFromSnapshotInput)
		if pgName != nil {
			dbi.DBClusterParameterGroupName = pgName
		}
		dbi.DBSubnetGroupName = c.DBSubnetGroup
		_, err := instances.RestoreSnapshotCluster(*dbi) // we might need to do something with the output in which case this changes
		if err != nil {
			return err
		}
	}

	// get three which is the instances create them in parallel
	third := c.S.Objects[3]
	waitChan := make(chan struct{}, MaxConcurrentJobs)
	for _, i := range third {
		waitChan <- struct{}{}
		go func(inst stack.Object) {
			o := inst.ReadObject()
			ins := o.(*rds.CreateDBInstanceInput)
			ins.DBSubnetGroupName = c.DBSubnetGroup
			_, err := instances.RestoreInstanceForCluster(*ins)
			if err != nil {
				slog.Error("error creating instance", "error", err)
			}
		}(i)
	}
	return nil
}

type CreateInstanceFromStackInput struct {
	Stack         *stack.Stack
	DBName        *string
	DBSubnetGroup *string
}

//CreateInstanceFromStack creates an RDS instance from a stack object
func (instances *DbInstances) CreateInstanceFromStack(c CreateInstanceFromStackInput) error {
	pgs := c.Stack.Objects[1]
	var pgName *string
	if len(pgs) == 0 {
		slog.Info("No parameter groups using the default parameter group")
	} else {
		for _, p := range pgs {
			pb := p.ReadObject()
			switch pb.(type) {
			case *pgstate.ParameterGroup:
				pg := pb.(*pgstate.ParameterGroup)
				pgName = pg.ParameterGroup.DBParameterGroupName
				_, err := instances.CreateParameterGroup(&pg.ParameterGroup)
				if err != nil {
					return err
				}
				batchSize := 20
				params := pg.Params
				batches := make([][]types.Parameter, 0, (len(params)+batchSize-1)/batchSize)
				for batchSize < len(params) {
					params, batches = params[batchSize:], append(batches, params[0:batchSize:batchSize])
				}
				batches = append(batches, params)
				for _, b := range batches {
					err = instances.ModifyParameterGroup(*pg.ParameterGroup.DBParameterGroupName, b)
					if err != nil {
						return err
					}
				}
			case *types.OptionGroup:
				og := pb.(*types.OptionGroup)
				_, err := instances.RestoreOptionGroup(*og.EngineName, *og.MajorEngineVersion, *og.OptionGroupName, *og.OptionGroupDescription)
				if err != nil {
					slog.Warn("failed to restore option group", "error", err)
				}
				optConfigs := optionsToConfiguration(og.Options)
				err = instances.ModifyOptionGroup(*og.OptionGroupName, optConfigs)
				if err != nil {
					slog.Warn("error modifying option group", "Error", err)
				}
			}
		}
		// Sleep for 5 minutes per AWS documentation to wait for a parameter group to be ready
		for i := 0; i < 10; i++ {
			slog.Info("waiting for five minutes for Parameter group per AWS documentation", "seconds", 30*i)
			time.Sleep(30 * time.Second)
		}
	}

	instance := c.Stack.Objects[2]
	slog.Info("starting to restore the instance")
	if len(instance) != 1 {
		slog.Error("No instances")
		return fmt.Errorf("There should only be a single instance")
	}
	for _, v := range instance {
		b := v.ReadObject()
		ins := b.(*rds.RestoreDBInstanceFromDBSnapshotInput)
		if pgName != nil {
			ins.DBParameterGroupName = pgName
		}
		if c.DBName != nil {
			ins.DBInstanceIdentifier = c.DBName
		}
		if c.DBSubnetGroup != nil {
			ins.DBSubnetGroupName = c.DBSubnetGroup
		}
		_, err := instances.RestoreSnapshotInstance(*ins)
		if err != nil {
			slog.Error("Failed to restore the instance", "error", err)
			return err
		}
		slog.Info("Database creation in progress")
	}
	return nil
}

func optionsToConfiguration(opts []types.Option) []types.OptionConfiguration {
	conf := []types.OptionConfiguration{}
	for _, v := range opts {
		dbsgs := []string{}
		for _, dbsg := range v.DBSecurityGroupMemberships {
			dbsgs = append(dbsgs, *dbsg.DBSecurityGroupName)
		}

		vpcsgs := []string{}
		for _, vpcsg := range v.VpcSecurityGroupMemberships {
			vpcsgs = append(vpcsgs, *vpcsg.VpcSecurityGroupId)
		}

		c := types.OptionConfiguration{
			OptionName:                  v.OptionName,
			DBSecurityGroupMemberships:  dbsgs,
			OptionSettings:              v.OptionSettings,
			OptionVersion:               v.OptionVersion,
			Port:                        v.Port,
			VpcSecurityGroupMemberships: vpcsgs,
		}
		conf = append(conf, c)
	}
	return conf
}

func (instances *DbInstances) getClusterStatus(name string) (*string, error) {
	cluster, err := instances.GetCluster(name)
	if err != nil {
		return nil, err
	}
	return cluster.Status, nil

}

// CreateSnapshot creates an AWS snapshot
// :instanceName - name of the database we want to backup
// :snapShotName name of the backup we are creating
func (instances *DbInstances) CreateSnapshot(instanceName string, snapshotName string) (
	*types.DBSnapshot, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	output, err := instances.RdsClient.CreateDBSnapshot(ctx, &rds.CreateDBSnapshotInput{
		DBInstanceIdentifier: aws.String(instanceName),
		DBSnapshotIdentifier: aws.String(snapshotName),
	})
	if err != nil {
		slog.Warn("Couldn't create snapshot", "snapshot", snapshotName, "error", err)
		return nil, err
	}
	return output.DBSnapshot, nil
}

//CreateClusterSnapshot so it turns out AWS is annoying and makes us create snapshots seperatly for clusters and instaces how fun!
func (instances *DbInstances) CreateClusterSnapshot(clusterName string, snapshotName string) (*types.DBClusterSnapshot, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	output, err := instances.RdsClient.CreateDBClusterSnapshot(ctx, &rds.CreateDBClusterSnapshotInput{
		DBClusterIdentifier:         aws.String(clusterName),
		DBClusterSnapshotIdentifier: aws.String(snapshotName),
	})
	if err != nil {
		slog.Warn("Couldn't create snapshot", "snapshot", snapshotName, "error", err)
		return nil, err
	}
	return output.DBClusterSnapshot, nil
}

// CopySnapshot copies a snapshot to a new region note it needs to run from the destination region so it needs a different client then CreateSnapshot!
func (instances *DbInstances) CopySnapshot(originalSnapshotName string, newSnapshotName string, sourceRegion string, KmsKey string) (
	*types.DBSnapshot, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	output, err := instances.RdsClient.CopyDBSnapshot(ctx, &rds.CopyDBSnapshotInput{
		SourceDBSnapshotIdentifier: aws.String(originalSnapshotName),
		TargetDBSnapshotIdentifier: aws.String(newSnapshotName),
		SourceRegion:               aws.String(sourceRegion), // this generates a presigned URL under the hood which enables cross region copies
		KmsKeyId:                   aws.String(KmsKey),
	})
	if err != nil {
		slog.Warn("Couldn't copy snapshot", "snapshot", originalSnapshotName, "error", err)
		return nil, err
	}
	return output.DBSnapshot, nil
}

//CopyClusterSnaphot see CopySnapshot now for a Cluster
func (instances *DbInstances) CopyClusterSnaphot(originalSnapshotName string, newSnapshotName string, sourceRegion string, kmsKey string) (
	*types.DBClusterSnapshot, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	output, err := instances.RdsClient.CopyDBClusterSnapshot(ctx, &rds.CopyDBClusterSnapshotInput{
		SourceDBClusterSnapshotIdentifier: aws.String(originalSnapshotName),
		TargetDBClusterSnapshotIdentifier: aws.String(newSnapshotName),
		SourceRegion:                      aws.String(sourceRegion),
		KmsKeyId:                          aws.String(kmsKey),
	})
	if err != nil {
		slog.Warn("Couldn't copy snapshot %s: %s\n", newSnapshotName, err)
		return nil, err
	}
	return output.DBClusterSnapshot, nil
}

//RestoreOptionGroup creates the option group for our database this is very optional
func (instances *DbInstances) RestoreOptionGroup(EngineName string, MajorEngineVersion string, OptionGroupName string, Description string) (
	*rds.CreateOptionGroupOutput, error) {
	input := rds.CreateOptionGroupInput{
		EngineName:             &EngineName,
		MajorEngineVersion:     &MajorEngineVersion,
		OptionGroupName:        &OptionGroupName,
		OptionGroupDescription: &Description,
	}
	out, err := instances.RdsClient.CreateOptionGroup(context.TODO(), &input)
	if err != nil {
		return nil, fmt.Errorf("Restore option group had an error %s", err)
	}
	return out, nil
}

//GetClusterParameterGroup get the cluster parameter group so we can make a new one in a new region or you know store it for restoration (actually we won't need to do that cause the data is stored on the snapshot :P)
func (instances *DbInstances) GetClusterParameterGroup(ParameterGroupName string) (
	*types.DBClusterParameterGroup, error) {
	output, err := instances.RdsClient.DescribeDBClusterParameterGroups(context.TODO(), &rds.DescribeDBClusterParameterGroupsInput{
		DBClusterParameterGroupName: aws.String(ParameterGroupName),
	})
	if err != nil {
		var notFoundError *types.DBClusterParameterGroupNotFoundFault
		if errors.As(err, &notFoundError) {
			slog.Warn("Parameter group does not exist.", "ParameterGroup", ParameterGroupName)
			err = nil
		} else {
			slog.Warn("Error getting parameter group", "ParameterGroup", ParameterGroupName, "Error", err)
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
		slog.Warn("Error getting parameters", "error", err)
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
			slog.Warn("Error getting parameters", "error", err)
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
			slog.Warn("Parameter group does not exist.", "parameterGroup", parameterGroupName)
			err = nil
		} else {
			slog.Warn("Error getting parameter group", "parametergroup", parameterGroupName, "Error", err)
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
		slog.Warn("Error getting parameters", "error", err)
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
			slog.Warn("Error getting parameters", "error", err)
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
		slog.Warn("error creating parameter group", "error", err)
		return output, err
	}
	return output, err
}

//CreateClusterParameterGroup creates a pararmeter group for a DB instance
func (instances *DbInstances) CreateClusterParameterGroup(p *types.DBClusterParameterGroup) (*rds.CreateDBClusterParameterGroupOutput, error) {
	input := rds.CreateDBClusterParameterGroupInput{
		DBParameterGroupFamily:      p.DBParameterGroupFamily,
		DBClusterParameterGroupName: p.DBClusterParameterGroupName,
		Description:                 p.Description,
	}
	output, err := instances.RdsClient.CreateDBClusterParameterGroup(context.TODO(), &input)
	if err != nil {
		slog.Warn("error creating parameter group ", "error", err)
		return output, err
	}
	return output, nil
}

//ModifyOptionGroup modifies the option group
func (instances *DbInstances) ModifyOptionGroup(OptionGroupName string, Include []types.OptionConfiguration) error {
	input := rds.ModifyOptionGroupInput{
		OptionGroupName:  &OptionGroupName,
		ApplyImmediately: true,
		OptionsToInclude: Include,
	}
	_, err := instances.RdsClient.ModifyOptionGroup(context.TODO(), &input)
	return err
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
			slog.Warn("error updating parameters", "Error", err)
		}
	}
	return nil
}

//ModifyClusterParameterGroup adds all the parameters to a db cluster parameter group
func (instances *DbInstances) ModifyClusterParameterGroup(pg string, parameters []types.Parameter) error {
	//batch this thing
	batchSize := 20
	batches := make([][]types.Parameter, 0, (len(parameters)+batchSize-1)/batchSize)

	for batchSize < len(parameters) {
		parameters, batches = parameters[batchSize:], append(batches, parameters[0:batchSize:batchSize])
	}
	batches = append(batches, parameters)

	for _, batch := range batches {
		_, err := instances.RdsClient.ModifyDBClusterParameterGroup(context.TODO(), &rds.ModifyDBClusterParameterGroupInput{
			DBClusterParameterGroupName: aws.String(pg),
			Parameters:                  batch,
		})
		if err != nil {
			slog.Warn("error updating parameters", "Error", err)
		}
	}
	return nil
}

//RestoreSnapshotCluster takes a snapshot turns it into a DB Cluster fun fact the cluster won't be ready from just this there will be no instances
func (instances *DbInstances) RestoreSnapshotCluster(input rds.RestoreDBClusterFromSnapshotInput) (*rds.RestoreDBClusterFromSnapshotOutput, error) {
	output, err := instances.RdsClient.RestoreDBClusterFromSnapshot(context.TODO(), &input)
	if err != nil {
		slog.Error("error creating snapshot cluster")
		return nil, err
	}
	return output, nil
}

// RestoreSnapshotInstance restores a single db instance from a snapshot
func (instances *DbInstances) RestoreSnapshotInstance(input rds.RestoreDBInstanceFromDBSnapshotInput) (*rds.RestoreDBInstanceFromDBSnapshotOutput, error) {

	output, err := instances.RdsClient.RestoreDBInstanceFromDBSnapshot(context.TODO(), &input)
	if err != nil {
		slog.Error("error creating instance from snapshot", "error", err)
		return nil, err
	}
	return output, nil
}

//RestoreInstanceForCluster our cluster has no instances by default it need's instances to be usable this makes them exist
func (instances *DbInstances) RestoreInstanceForCluster(input rds.CreateDBInstanceInput) (*rds.CreateDBInstanceOutput, error) {
	output, err := instances.RdsClient.CreateDBInstance(context.TODO(), &input)
	if err != nil {
		slog.Error("error creating instance", "error", err)
		return nil, err
	}
	return output, nil
}

//GetInstanceSnapshotARN get the arn for an instance snapshot
func (instances *DbInstances) GetInstanceSnapshotARN(name string, marker *string) (*string, error) {
	output, err := instances.RdsClient.DescribeDBSnapshots(context.TODO(), &rds.DescribeDBSnapshotsInput{
		// Marker: marker,
		DBSnapshotIdentifier: aws.String(name),
	})
	if err != nil {
		slog.Error("error with snapshots", "error", err)
		return nil, fmt.Errorf("error retreiving snapshot: %s", err)
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
