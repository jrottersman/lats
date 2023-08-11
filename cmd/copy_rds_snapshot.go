package cmd

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/jrottersman/lats/aws"
	"github.com/jrottersman/lats/helpers"
	"github.com/jrottersman/lats/state"
	"github.com/spf13/cobra"
)

var (
	//Variables used for flags
	kmsKey               string
	originalSnapshotName string
	copySnapshotName     string

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

	fmt.Printf("TODO implement me, %v so this passes", sm)
}

func createKMSKey(config Config, sm state.StateManager) string {
	var kmsStruct *types.KeyMetadata
	c := aws.InitKms(config.BackupRegion)
	kmsStruct, err := c.CreateKMSKey()
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
