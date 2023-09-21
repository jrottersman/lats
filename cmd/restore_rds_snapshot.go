package cmd

import (
	"fmt"

	"github.com/jrottersman/lats/aws"
	"github.com/jrottersman/lats/state"
	"github.com/spf13/cobra"
)

var (
	//Variables used for flags
	restoreSnapshotName string
	restoreDbName       string

	//RestoreRDSSnapshotCmd restores an RDS snapshot
	RestoreRDSSnapshotCmd = &cobra.Command{
		Use:     "restoreRDSSnapshot",
		Aliases: []string{"RestoreSnapshot"},
		Short:   "Restores an RDS snapshot",
		Long:    "Restores an RDS snapshot",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Implement me ")
		},
	}
)

func init() {
	RestoreRDSSnapshotCmd.Flags().StringVarP(&restoreSnapshotName, "snapshot-name", "s", "", "name of the snapshot we want to restore: choose one of snapshotName or db name")
	RestoreRDSSnapshotCmd.Flags().StringVarP(&restoreDbName, "database-name", "d", "", "name of the database we want to restore the snapshot for")
}

//RestoreSnapshot is the function that restores a snapshot
func RestoreSnapshot(instances *aws.DbInstances, stateKV state.StateManager, snapshotName string) error {
	SnapshotStack := FindStack(stateKV, snapshotName)
	if SnapshotStack.RestorationObjectName == state.Cluster {
		fmt.Printf("restoring a cluster")
	} else if SnapshotStack.RestorationObjectName == state.LoneInstance {
		fmt.Printf("restoring an instance")
	}
	return nil
}
