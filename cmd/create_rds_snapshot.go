package cmd

import (
	"fmt"
	"log"

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
			run()
		},
	}
)

func init() {
	CreateRDSSnapshotCmd.Flags().StringVarP(&dbName, "database-name", "d", "", "Database name we want to create the snapshot for")
}

func run() {
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
	// Create an RDS client
	dbi := aws.Init(config.MainRegion)

	// Get Instance
	i, err := dbi.GetInstance(dbName)
	if err != nil {
		log.Fatalf("Error %s trying to retrive instance %s\n", err, dbName)
	}

	// Update state with instance
	b := state.EncodeRDSDatabaseOutput(i)
	f1 := helpers.RandomStateFileName()
	_, err = state.WriteOutput(*f1, b)
	if err != nil {
		log.Fatalf("failed to write state file: %s\n", err)
	}
	sm.UpdateState(dbName, *f1, "RDSSnapshot")

	// Copy Snapshot
	snapName := helpers.SnapshotName(dbName)
	snap, err := dbi.CreateSnapshot(dbName, snapName)
	if err != nil {
		log.Fatalf("failed to create snapshot %s\n", err)
	}

	// Update State with snapshot
	f2 := helpers.RandomStateFileName()
	b2 := state.EncodeRDSSnapshotOutput(snap)
	_, err = state.WriteOutput(*f2, b2)
	if err != nil {
		log.Fatalf("failed to write state file: %s\n", err)
	}
	sm.UpdateState(*snap.DBSnapshotIdentifier, *f2, state.SnapshotType)
	sm.SyncState(stateFileName)
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
		CreateSnasphotForInstance()
	}
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
	stackFn := fmt.Sprintf(".state/%s", helpers.RandomStateFileName())
	stack.Write(stackFn)
	sm.UpdateState(snapshotName, stackFn, "stack")
}

func CreateSnasphotForInstance() {
	log.Printf("implement me")
}
