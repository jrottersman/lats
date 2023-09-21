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
	region              string

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
	RestoreRDSSnapshotCmd.Flags().StringVarP(&region, "region", "r", "", "AWS region we are restoring in")
}

//RestoreSnapshot is the function that restores a snapshot
func RestoreSnapshot(instances *aws.DbInstances, stateKV state.StateManager, snapshotName string) error {
	dbi := aws.Init(region)
	SnapshotStack := FindStack(stateKV, snapshotName)
	if SnapshotStack.RestorationObjectName == state.Cluster {
		return dbi.CreateClusterFromStack(SnapshotStack)
		// TODO finish create Cluster from stack
	} else if SnapshotStack.RestorationObjectName == state.LoneInstance {
		fmt.Printf("restoring an instance")
		//TODO create CreateInstanceFromStack
	}
	return fmt.Errorf("Error invalid type of stack to restore a snapshot")
}
