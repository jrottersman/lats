package cmd

import (
	"github.com/jrottersman/lats/helpers"
	"github.com/spf13/cobra"
)

var (
	initCmd = &cobra.Command{
		Use:     "init",
		Aliases: []string{"initialize", "initialise", "config"},
		Short:   "Initalizes lats and configures it for creating backups",
		Long:    "Initalize (lats init) will setup lats with the correct regions and let you choose where you want to store state",
		Run: func(cmd *cobra.Command, args []string) {
			getMainRegion()
		},
	}
)

func getMainRegion() string {
	mainRegionPromptContent := helpers.PromptContent{
		"Please provide an AWS region.",
		"What is the AWS region your database is running in?",
	}
	mainRegion := helpers.PromptGetInput(mainRegionPromptContent)
	return mainRegion
}
