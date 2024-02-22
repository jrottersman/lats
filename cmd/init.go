package cmd

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/jrottersman/lats/helpers"
	"github.com/jrottersman/lats/state"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Config tells us our regions and where the state file is
type Config struct {
	MainRegion    string `json:"mainRegion"`
	BackupRegion  string `json:"backupRegion"`
	StateFileName string `json:"stateFileName"`
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
			slog.Info("Initalizing lats")
			viper.SetConfigName(".latsConfig")
			viper.SetConfigType("json")
			viper.AddConfigPath(".")

			if err := viper.ReadInConfig(); err != nil {
				if _, ok := err.(viper.ConfigFileNotFoundError); ok {
					genConfigFile()
				} else {
					slog.Error("Error parsing config file ", "error", err)
				}
			}

			// Config file found and successfully parsed
			state.InitState(".confState.json")
			os.Mkdir(".state", os.ModePerm)

		},
	}
)

func genConfigFile() {
	slog.Info("Generating config file")
	c := genConfig(getMainRegion, getBackupRegion)
	writeConfig(c, ".latsConfig.json")
	state.InitState(".confState.json")
	slog.Info("creating .state directory")
	os.Mkdir(".state", os.ModePerm)
}
func getMainRegion() string {
	if mainRegion != "" {
		return mainRegion
	}

	mainRegionPromptContent := helpers.PromptContent{
		ErrorMsg: "Please provide an AWS region.",
		Label:    "What is the AWS region your database is running in?",
	}
	mainRegion = helpers.PromptGetInput(mainRegionPromptContent)
	return mainRegion
}

func getBackupRegion() string {
	if backupRegion != "" {
		return backupRegion
	}
	backupRegionPromptContent := helpers.PromptContent{
		ErrorMsg: "Please provide an AWS region.",
		Label:    "What is the AWS region your backup should be in?",
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
		MainRegion:    mainRegion,
		BackupRegion:  backupRegion,
		StateFileName: ".confState.json",
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
		slog.Warn("Error writing config: ", "error", err)
		return err
	}
	err = os.WriteFile(filename, conf, 0644)
	if err != nil {
		slog.Warn("Error writing config", "error", err)
		return err
	}
	return nil
}

func readConfig(filename string) (Config, error) {
	confFile, err := os.ReadFile(filename)
	if err != nil {
		slog.Info("Error reading config,", "err", err)
	}
	conf := Config{}
	err = json.Unmarshal(confFile, &conf)
	if err != nil {
		slog.Warn("Error unmarshalling json", "error", err)
	}
	return conf, err
}
