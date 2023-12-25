package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/jrottersman/lats/aws"
	"github.com/jrottersman/lats/helpers"
	"github.com/jrottersman/lats/rdsstate"
	"github.com/jrottersman/lats/state"
	"github.com/spf13/cobra"
)

var (
	// Variables used for flags
	dbName       string
	snapshotName string
	//CreateRDSSnapshotCmd is the args for creating RDS snapshot call
	CreateRDSSnapshotCmd = &cobra.Command{
		Use:     "CreateRDSSnapshot",
		Aliases: []string{"CreateSnapshot"},
		Short:   "Creates a snapshot for a given DB",
		Long:    "Creates a snapshot for an RDS or Aurora database",
		Run: func(cmd *cobra.Command, args []string) {
			CreateSnapshot()
		},
	}
)

func init() {
	CreateRDSSnapshotCmd.Flags().StringVarP(&dbName, "database-name", "d", "", "Database name we want to create the snapshot for")
	CreateRDSSnapshotCmd.Flags().StringVarP(&snapshotName, "snapshot-name", "s", "", "Snapshot name that we want to create our snapshot with")
}

//CreateSnapshot generates a snapshot in AWS
func CreateSnapshot() {
	//Get Config and state
	config, sm := GetState()
	dbi := aws.Init(config.MainRegion)
	ec2 := aws.InitEc2(config.MainRegion)
	cluster, err := dbi.GetCluster(dbName)
	if err != nil {
		slog.Info("not a cluster with step 1 get cluster ", "error", err)
	}
	if cluster == nil {
		c := CreateInstanceSnapshotInput{
			dbi: dbi,
			ec2: ec2,
			sm:  sm,
			sfn: config.StateFileName,
		}
		createSnapshotForInstance(c)
	} else {
		c := CreateClusterSnapshotInput{
			dbi:     dbi,
			ec2:     ec2,
			sm:      sm,
			cluster: cluster,
			sfn:     config.StateFileName,
		}
		createSnapshotForCluster(c)
	}
}

func createSnapshotForCluster(c CreateClusterSnapshotInput) {
	slog.Info("creating snapshot for cluster")
	snapshot, err := c.dbi.CreateClusterSnapshot(dbName, snapshotName)
	if err != nil {
		slog.Error("error creating snapshot", "error", err)
		os.Exit(1)
	}
	// create a stack
	store := state.RDSRestorationStore{
		Cluster:         c.cluster,
		ClusterSnapshot: snapshot,
	}
	// Get Security groups to add to the stack
	var sgOutput state.SecurityGroupOutput
	sgs := c.cluster.VpcSecurityGroups
	if len(sgs) != 0 {
		out, err := getSGs(c.ec2, sgs)
		if err != nil {
			slog.Error("can not get security groups", "error", err)
		}
		var groups []ec2types.SecurityGroup
		for _, v := range out {
			groups = append(groups, v.SecurityGroups...)
		}
		sgOutput = state.SecurityGroupOutput{SecurityGroups: groups}
	}
	input := rdsstate.ClusterStackInput{
		R:              store,
		StackName:      snapshotName,
		Client:         c.dbi,
		SecurityGroups: &sgOutput,
		Folder:         ".state",
	}
	slog.Info("generating the stack")
	stack, err := rdsstate.GenerateRDSClusterStack(input)
	if err != nil {
		slog.Error("error generating stack ", "error", err)
		os.Exit(1)
	}
	counter := 0
	for {
		status, err := c.dbi.GetClusterSnapshotStatus(snapshotName)
		if err != nil {
			slog.Error("error getting status", "error", err)
		}
		if *status == "available" {
			break
		}
		if counter == 10 {
			break
		}
		slog.Info("snapshot creation in progess", "Status", *status)
		counter++
		time.Sleep(30 * time.Second)
	}
	stackFn := fmt.Sprintf(".state/%s", *helpers.RandomStateFileName())
	slog.Info("Writing the stack")
	err = stack.Write(stackFn)
	if err != nil {
		slog.Error("error writing stack ", "error", err)
		os.Exit(1)
	}
	c.sm.UpdateState(snapshotName, stackFn, "stack")
	c.sm.SyncState(c.sfn)
	slog.Info("Snapshot created")
}

func createSnapshotForInstance(c CreateInstanceSnapshotInput) {
	slog.Info("starting create snapshot for instance")
	db, err := c.dbi.GetInstance(dbName)
	if err != nil {
		slog.Warn("didn't get instance", "problem", err)
	}
	sgs := db.VpcSecurityGroups
	slog.Debug("creating snapshot")
	snapshot, err := c.dbi.CreateSnapshot(dbName, snapshotName)
	if err != nil {
		slog.Error("error creating snapshot: ", "error", err)
	}

	store := state.RDSRestorationStore{
		Instance: db,
		Snapshot: snapshot,
	}
	slog.Debug("getting parameter groups")
	pgs, err := aws.GetParameterGroups(store, c.dbi)
	if err != nil {
		slog.Warn("error getting parameter groups", "error", err)
	}
	stackInput := rdsstate.InstanceStackInputs{
		R:               store,
		StackName:       snapshotName,
		ParameterGroups: pgs,
		SecurityGroups:  sgs,
	}
	slog.Debug("generating stack")
	stack, err := rdsstate.GenerateRDSInstanceStack(stackInput)
	if err != nil {
		slog.Warn("error generating stack", "error", err)
	}
	stackFn := fmt.Sprintf(".state/%s", *helpers.RandomStateFileName())
	err = stack.Write(stackFn)
	if err != nil {
		slog.Warn("error writing stack", "error", err)
	}
	counter := 0
	for {
		status, err := c.dbi.GetInstanceSnapshotStatus(snapshotName)
		if err != nil {
			slog.Error("error getting status", "error", err)
		}
		if *status == "available" {
			break
		}
		if counter == 10 {
			break
		}
		slog.Info("snapshot creation in progess", "Status", *status)
		counter++
		time.Sleep(30 * time.Second)
	}
	c.sm.UpdateState(snapshotName, stackFn, "stack")
	c.sm.SyncState(c.sfn)
}

//GetState reads in our statefile and config for future processing
func GetState() (Config, state.StateManager) {
	slog.Debug("getting config")
	config, err := readConfig(".latsConfig.json")
	if err != nil {
		slog.Warn("Error reading config", "error", err)
	}
	slog.Debug("Getting state")
	stateFileName := config.StateFileName
	sm, err := state.ReadState(stateFileName)
	if err != nil {
		slog.Warn("Error reading state", "error", err)
	}
	return config, sm
}

func getSGs(ec2 aws.EC2Instances, sgs []types.VpcSecurityGroupMembership) ([]state.SecurityGroupOutput, error) {
	var sgOut []state.SecurityGroupOutput
	for _, sg := range sgs {
		id := sg.VpcSecurityGroupId
		sgO, err := ec2.DescribeSG(*id)
		if err != nil {
			slog.Warn("describing SG", "error", err)
			return nil, err
		}
		sgs := sgO.SecurityGroups
		sgOut = append(sgOut, state.SecurityGroupOutput{
			SecurityGroups: sgs,
		})

	}
	return sgOut, nil
}
