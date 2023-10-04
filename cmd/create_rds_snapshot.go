package cmd

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/jrottersman/lats/aws"
	"github.com/jrottersman/lats/helpers"
	"github.com/jrottersman/lats/rdsState"
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
	cluster, err := dbi.GetCluster(dbName)
	if err != nil {
		log.Fatalf("error with step 1 get cluster %s", err)
	}
	if cluster == nil && err == nil {
		createSnapshotForInstance(dbi, sm, config.StateFileName)
	} else {
		createSnapshotForCluster(dbi, sm, cluster, config.StateFileName)
	}
}

func createSnapshotForCluster(dbi aws.DbInstances, sm state.StateManager, cluster *types.DBCluster, sfn string) {
	snapshot, err := dbi.CreateClusterSnapshot(dbName, snapshotName)
	if err != nil {
		log.Fatalf("error creating snapshot %s", err)
	}
	// create a stack
	store := state.RDSRestorationStore{
		Cluster:         cluster,
		ClusterSnapshot: snapshot,
	}
	stack, err := rdsState.GenerateRDSClusterStack(store, snapshotName, nil, dbi, ".state")
	if err != nil {
		log.Fatalf("error generating stack %s", err)
	}
	stackFn := fmt.Sprintf(".state/%s", *helpers.RandomStateFileName())
	err = stack.Write(stackFn)
	if err != nil {
		log.Fatalf("error writing stack %s", err)
	}
	sm.UpdateState(snapshotName, stackFn, "stack")
	sm.SyncState(sfn)
}

func createSnapshotForInstance(dbi aws.DbInstances, sm state.StateManager, sfn string) {
	db, err := dbi.GetInstance(dbName)
	if err != nil {
		log.Printf("didn't get instance %s", err)
	}
	snapshot, err := dbi.CreateSnapshot(dbName, snapshotName)
	if err != nil {
		log.Fatalf("error creating snapshot: %s", err)
	}

	store := state.RDSRestorationStore{
		Instance: db,
		Snapshot: snapshot,
	}
	pgs, err := aws.GetParameterGroups(store, dbi)
	if err != nil {
		fmt.Printf("error getting parameter groups %s", err)
	}

	stack, err := rdsState.GenerateRDSInstanceStack(store, snapshotName, nil, nil, pgs)
	if err != nil {
		log.Fatalf("error generating stack %s", err)
	}
	stackFn := fmt.Sprintf(".state/%s", *helpers.RandomStateFileName())
	err = stack.Write(stackFn)
	if err != nil {
		log.Fatalf("error writing stack %s", err)
	}
	sm.UpdateState(snapshotName, stackFn, "stack")
	sm.SyncState(sfn)
}

//GetState reads in our statefile and config for future processing
func GetState() (Config, state.StateManager) {
	config, err := readConfig(".latsConfig.json")
	if err != nil {
		log.Fatalf("Error reading config %s", err)
	}
	stateFileName := config.StateFileName
	sm, err := state.ReadState(stateFileName)
	if err != nil {
		log.Fatalf("Error reading state %s", err)
	}
	return config, sm
}
