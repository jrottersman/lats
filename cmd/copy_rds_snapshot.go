package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	CopyRDSSnapshotCmd = &cobra.Command{
		Use:     "CopyRDSSnapshot",
		Aliases: []string{"CopySnapshot"},
		Short:   "Copies a snapshot for a given DB",
		Long:    "Copies a snapshot for an RDS or Aurora database into a new region",
		Run: func(cmd *cobra.Command, args []string) {
			createSnapshot()
		},
	}
)

func createSnapshot() {
	fmt.Println("TODO implement me")
}
