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

func CreateSnapshot() {
	//Get Config and state
	config, err := readConfig(".latsConfig.json")
	if err != nil {
		log.Fatalf("Error reading config %s", err)
	}
	stateFileName := config.StateFileName
	sm, err := state.ReadState(stateFileName)
	if err != nil {
		log.Fatalf("Error reading state %s", err)
	}
	dbi := aws.Init("us-east-1")
	cluster, err := dbi.GetCluster(dbName)
	if err != nil {
		log.Fatalf("error with step 1 get cluster %s", err)
	}
	if cluster == nil && err == nil {
		createSnapshotForInstance(dbi, sm)
	} else {
		createSnapshotForCluster(dbi, sm, cluster)
	}
}

func createSnapshotForCluster(dbi aws.DbInstances, sm state.StateManager, cluster *types.DBCluster) {
	snapshot, err := dbi.CreateClusterSnapshot(dbName, snapshotName)
	if err != nil {
		log.Fatalf("error creating snapshot %s", err)
	}
	// create a stack
	store := state.RDSRestorationStore{
		Cluster:         cluster,
		ClusterSnapshot: snapshot,
	}
	stack, err := rdsState.GenerateRDSClusterStack(store, dbName, nil, dbi, ".state")
	if err != nil {
		log.Printf("error generating stack %s", err)
	}
	stackFn := fmt.Sprintf(".state/%s", *helpers.RandomStateFileName())
	stack.Write(stackFn)
	sm.UpdateState(snapshotName, stackFn, "stack")
}

func createSnapshotForInstance(dbi aws.DbInstances, sm state.StateManager) {
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
	stack, err := state.GenerateRDSInstanceStack(store, *store.GetInstanceIdentifier(), helpers.RandomStateFileName())
	if err != nil {
		log.Printf("error generating stack %s", err)
	}
	stackFn := fmt.Sprintf(".state/%s", *helpers.RandomStateFileName())
	stack.Write(stackFn)
	sm.UpdateState(snapshotName, stackFn, "stack")
}
