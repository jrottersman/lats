package cmd

import (
	"github.com/jrottersman/lats/helpers"
	"github.com/spf13/cobra"
)

type Config struct {
	MainRegion   string
	BackupRegion string
}

var (
	//Used for flags
	mainRegion   string
	backupRegion string

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
	if mainRegion != "" {
		return mainRegion
	}

	mainRegionPromptContent := helpers.PromptContent{
		"Please provide an AWS region.",
		"What is the AWS region your database is running in?",
	}
	mainRegion = helpers.PromptGetInput(mainRegionPromptContent)
	return mainRegion
}

func getBackupRegion() string {
	if backupRegion != "" {
		return backupRegion
	}
	backupRegionPromptContent := helpers.PromptContent{
		"Please provide an AWS region.",
		"What is the AWS region your database is running in?",
	}
	backupRegion := helpers.PromptGetInput(backupRegionPromptContent)
	return backupRegion
}

func init() {
	initCmd.Flags().StringVarP(&mainRegion, "main-region", "", "", "AWS Region the application is running in")
	initCmd.Flags().StringVarP(&backupRegion, "backup-region", "", "", "AWS region we want backup the application to")
}

func newConfig(mainRegion string, backupRegion string) Config {
	return Config{
		MainRegion:   mainRegion,
		BackupRegion: backupRegion,
	}
}
