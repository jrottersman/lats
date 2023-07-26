package cmd

import "github.com/spf13/cobra"


var (
	initCmd = &cobra.Command{
		Use:     "init",
		Aliases: []string{"initialize", "initialise", "config"},
		Short:   "Initalizes lats and configures it for creating backups",
		Long:    "Initalize (lats init) will setup lats with the correct regions and let you choose where you want to store state",
	}
)
