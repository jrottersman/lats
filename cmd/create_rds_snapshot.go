package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"time"

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
		createSnapshotForInstance(dbi, ec2, sm, config.StateFileName)
	} else {
		createSnapshotForCluster(dbi, ec2, sm, cluster, config.StateFileName)
	}
}

func createSnapshotForCluster(dbi aws.DbInstances, ec2 aws.EC2Instances, sm state.StateManager, cluster *types.DBCluster, sfn string) {
	slog.Info("creating snapshot for cluster")
	snapshot, err := dbi.CreateClusterSnapshot(dbName, snapshotName)
	if err != nil {
		slog.Error("error creating snapshot", "error", err)
		os.Exit(1)
	}
	// create a stack
	store := state.RDSRestorationStore{
		Cluster:         cluster,
		ClusterSnapshot: snapshot,
	}
	input := rdsstate.ClusterStackInput{
		R:         store,
		StackName: snapshotName,
		Client:    dbi,
		Folder:    ".state",
	}
	slog.Info("generating the stack")
	stack, err := rdsstate.GenerateRDSClusterStack(input)
	if err != nil {
		slog.Error("error generating stack ", "error", err)
		os.Exit(1)
	}
	counter := 0
	for {
		status, err := dbi.GetClusterSnapshotStatus(snapshotName)
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
	sm.UpdateState(snapshotName, stackFn, "stack")
	sm.SyncState(sfn)
	slog.Info("Snapshot created")
}

func createSnapshotForInstance(dbi aws.DbInstances, ec2 aws.EC2Instances, sm state.StateManager, sfn string) {
	slog.Info("starting create snapshot for instance")
	db, err := dbi.GetInstance(dbName)
	if err != nil {
		slog.Warn("didn't get instance", "problem", err)
	}
	sgs := db.VpcSecurityGroups
	slog.Debug("creating snapshot")
	snapshot, err := dbi.CreateSnapshot(dbName, snapshotName)
	if err != nil {
		slog.Error("error creating snapshot: ", "error", err)
	}

	store := state.RDSRestorationStore{
		Instance: db,
		Snapshot: snapshot,
	}
	slog.Debug("getting parameter groups")
	pgs, err := aws.GetParameterGroups(store, dbi)
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
		status, err := dbi.GetInstanceSnapshotStatus(snapshotName)
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
	sm.UpdateState(snapshotName, stackFn, "stack")
	sm.SyncState(sfn)
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
