package cmd

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/jrottersman/lats/aws"
	"github.com/jrottersman/lats/helpers"
	"github.com/jrottersman/lats/state"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	//Variables used for flags
	kmsKey               string
	originalSnapshotName string
	copySnapshotName     string
	configFile           string

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

func init() {
	CopyRDSSnapshotCmd.Flags().StringVarP(&kmsKey, "kms-key", "k", "", "KMS key to use for the snapshot optional")
	CopyRDSSnapshotCmd.Flags().StringVarP(&copySnapshotName, "copy-snapshot", "c", "", "Name of the snapshot copy we are creating")
	CopyRDSSnapshotCmd.Flags().StringVarP(&originalSnapshotName, "snapshot", "s", "", "Snapshot we want to copy")
	CopyRDSSnapshotCmd.Flags().StringVarP(&configFile, "config-file", "f", "", "Config file for the snapshot that we want to parse")
}

func createSnapshot() {
	config, err := readConfig(".latsConfig.json")
	if err != nil {
		log.Fatalf("Error reading config %s", err)
	}
	stateFileName := config.StateFileName
	sm, err := state.ReadState(stateFileName)
	if err != nil {
		log.Fatalf("Error reading state %s", err)
	}

	if configFile != "" {
		viper.SetConfigFile(configFile)
		viper.AddConfigPath(".")
		err := viper.ReadInConfig() // Find and read the config file
		if err != nil {             // Handle errors reading the config file
			fmt.Errorf("Error reading config file: %w", err)
		}
		getKey := fmt.Sprintf("%s", viper.Get("kmsKey"))
		kmsKey = getKey

		origSnap := fmt.Sprintf("%s", viper.Get("OriginalSnapshotName"))
		originalSnapshotName = origSnap

		copySnap := fmt.Sprintf("%s", viper.Get("copySnapshotName"))
		copySnapshotName = copySnap
	}

	// Create KMS key
	if kmsKey == "" {
		kmsKey = createKMSKey(config, sm)
	}
	// Get RDS Client
	dbi := aws.Init(config.BackupRegion)

	// Copy Snapshot
	snap, err := dbi.CopySnapshot(originalSnapshotName, copySnapshotName, config.MainRegion, kmsKey)
	if err != nil {
		log.Fatalf("Error copying snapshot %s", err)
	}
	f2 := helpers.RandomStateFileName()
	b2 := state.EncodeRDSSnapshotOutput(snap)
	_, err = state.WriteOutput(*f2, b2)
	if err != nil {
		log.Fatalf("failed to write state file: %s\n", err)
	}
	sm.UpdateState(*snap.DBSnapshotIdentifier, *f2)
	sm.SyncState(stateFileName)
}

func createKMSKey(config Config, sm state.StateManager) string {
	var kmsStruct *types.KeyMetadata
	c := aws.InitKms(config.BackupRegion)
	kmsStruct, err := c.CreateKMSKey(nil)
	if err != nil {
		log.Fatalf("failed creating KMS key %s", err)
	}
	kf := helpers.RandomStateFileName()
	b := state.EncodeKmsOutput(kmsStruct)
	_, err = state.WriteOutput(*kf, b)
	if err != nil {
		log.Printf("Issues writing state %s", err)
	}
	sm.UpdateState(*kmsStruct.KeyId, *kf)
	return *kmsStruct.KeyId

}
