package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	CreateRDSSnapshotCmd = &cobra.Command{
		Use:     "CreateRDSSnapshot",
		Aliases: []string{"CreateSnapshot"},
		Short:   "Creates a snapshot for a given DB",
		Long:    "Creates a snapshot for an RDS or Aurora database",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Copies a snapshot to a new region")
		},
	}
)
