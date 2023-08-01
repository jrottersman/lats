package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jrottersman/lats/helpers"
	"github.com/spf13/cobra"
)

type Config struct {
	MainRegion   string `json:"mainRegion"`
	BackupRegion string `json:"backupRegion"`
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
			c := genConfig(getMainRegion, getBackupRegion)
			writeConfig(c, ".latsConfig.json")
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
		"What is the AWS region your backup should be in?",
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

func genConfig(mr, br func() string) Config {
	mainRegion := mr()
	backupRegion := br()
	return newConfig(mainRegion, backupRegion)
}

func writeConfig(c Config, filename string) error {
	conf, err := json.Marshal(c)
	if err != nil {
		fmt.Errorf("Error writing config: %s", err)
		return err
	}
	err = os.WriteFile(filename, conf, 0644)
	return nil
}

func readConfig(filename string) (Config, error) {
	// LOOK INTO IO PACKAGE there is probably something that fixes this mess
	confFile, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading config: %s", err)
	}
	conf := Config{}
	err = json.Unmarshal(confFile, &conf)
	if err != nil {
		fmt.Printf("Error unmarshalling json: %s\n", err)
	}
	return conf, err
}
