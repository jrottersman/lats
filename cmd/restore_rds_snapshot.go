package cmd

import (
	"fmt"
	"log/slog"

	"github.com/jrottersman/lats/aws"
	"github.com/jrottersman/lats/stack"
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
			_, sm := GetState()
			RestoreSnapshot(sm, restoreSnapshotName)
		},
	}
)

func init() {
	RestoreRDSSnapshotCmd.Flags().StringVarP(&restoreSnapshotName, "snapshot-name", "s", "", "name of the snapshot we want to restore: choose one of snapshotName or db name")
	RestoreRDSSnapshotCmd.Flags().StringVarP(&restoreDbName, "database-name", "d", "", "name of the database we want to restore the snapshot for")
	RestoreRDSSnapshotCmd.Flags().StringVarP(&region, "region", "r", "", "AWS region we are restoring in")
}

//RestoreSnapshot is the function that restores a snapshot
func RestoreSnapshot(stateKV state.StateManager, restoreSnapshotName string) error {
	slog.Info("Starting restore snapshot procedure")
	dbi := aws.Init(region)
	slog.Info("finding the stack")
	SnapshotStack, err := FindStack(stateKV, restoreSnapshotName)
	if err != nil {
		slog.Error("Error finding stack", "error", err)
	}
	slog.Info("Stack is", "stack", SnapshotStack)
	if SnapshotStack.RestorationObjectName == stack.Cluster {
		slog.Info("Restoring a cluster")
		return dbi.CreateClusterFromStack(SnapshotStack)
	} else if SnapshotStack.RestorationObjectName == stack.LoneInstance {
		slog.Info("Restoring an Instance")
		return dbi.CreateInstanceFromStack(SnapshotStack)
	}
	slog.Error("Invalid type of stack for restoring an object", "StackType", SnapshotStack.RestorationObjectName)
	return fmt.Errorf("Error invalid type of stack to restore a snapshot")
}
