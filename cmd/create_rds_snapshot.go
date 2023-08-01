package cmd

import (
	"fmt"

	"github.com/jrottersman/lats/aws"
	"github.com/spf13/cobra"
)

var (
	// Variables used for flags
	dbName string

	CreateRDSSnapshotCmd = &cobra.Command{
		Use:     "CreateRDSSnapshot",
		Aliases: []string{"CreateSnapshot"},
		Short:   "Creates a snapshot for a given DB",
		Long:    "Creates a snapshot for an RDS or Aurora database",
		Run: func(cmd *cobra.Command, args []string) {
			conf, err := readConfig(".latsConfig.json")
			if err != nil {
				fmt.Printf("Error is %s\n", err)
			}
			dbi := aws.Init(conf.MainRegion) // TODO read this from config
			dbi.GetInstance(dbName)
		},
	}
)

func init() {
	CreateRDSSnapshotCmd.Flags().StringVarP(&dbName, "database-name", "d", "", "Database name we want to create the snapshot for")
}
