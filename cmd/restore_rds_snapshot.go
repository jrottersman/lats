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

func RestoreSnapshot(instances *aws.DbInstances, stateKV state.StateManager, snapshotName string) error {
	RestorationBuilder, err := state.RDSRestorationStoreBuilder(stateKV, snapshotName)
	if err != nil {
		fmt.Printf("error getting restoration store %s", err)
		return err
	}
	SnapshotInput := state.GenerateRestoreDBInstanceFromDBClusterSnapshotInput(*RestorationBuilder)
	_, nil := instances.RestoreSnapshotInstance(*SnapshotInput)
	return nil
}
