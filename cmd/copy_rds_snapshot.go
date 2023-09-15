package cmd

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/aws/aws-sdk-go-v2/service/rds"
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
	//CopyRDSSnapshotCmd creates the copy snapshot command.
	CopyRDSSnapshotCmd = &cobra.Command{
		Use:     "CopyRDSSnapshot",
		Aliases: []string{"CopySnapshot"},
		Short:   "Copies a snapshot for a given DB",
		Long:    "Copies a snapshot for an RDS or Aurora database into a new region",
		Run: func(cmd *cobra.Command, args []string) {
			copySnapshot()
		},
	}
)

func init() {
	CopyRDSSnapshotCmd.Flags().StringVarP(&kmsKey, "kms-key", "k", "", "KMS key to use for the snapshot optional")
	CopyRDSSnapshotCmd.Flags().StringVarP(&copySnapshotName, "copy-snapshot", "c", "", "Name of the snapshot copy we are creating")
	CopyRDSSnapshotCmd.Flags().StringVarP(&originalSnapshotName, "snapshot", "s", "", "Snapshot we want to copy")
	CopyRDSSnapshotCmd.Flags().StringVarP(&configFile, "config-file", "f", "", "Config file for the snapshot that we want to parse")
}

func copySnapshot() {
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
			fmt.Printf("Error reading config file: %s", err)
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
	sm.UpdateState(*snap.DBSnapshotIdentifier, *f2, state.SnapshotType)
	sm.SyncState(stateFileName)
}

func createKMSKey(config Config, sm state.StateManager) string {
	var kmsStruct *types.KeyMetadata
	c := aws.InitKms(config.BackupRegion)
	kmsStruct, err := c.CreateKMSKey(nil)
	if err != nil {
		log.Fatalf("failed creating KMS key %s", err)
	}
	return *kmsStruct.KeyId

}

//FindStack get's a stack for creating our new stack when we copy the snapshot
func FindStack(sm state.StateManager, snapshot string) *state.Stack {
	sm.Mu.Lock()
	defer sm.Mu.Unlock()

	for _, v := range sm.StateLocations {
		if v.ObjectType != "stack" {
			continue
		}
		stack, err := state.ReadStack(v.FileLocation)
		if err != nil {
			log.Printf("error reading stack %s", err)
		}
		if stack.Name == snapshot {
			return stack
		}
	}
	return nil
}

// NewStack generates the new stack that we are going touse
func NewStack(oldStack state.Stack, ending string) *state.Stack {
	objs := make(map[int][]state.Object)
	for k, v := range oldStack.Objects {
		objs[k] = []state.Object{}
		for _, i := range v {
			obj := i.ReadObject()
			switch i.ObjType {
			case state.LoneInstance:
				s := getLoneInstanceObject(obj, ending, k)
				objs[k] = append(objs[k], s)
			}
		}
	}
	return &state.Stack{
		Name:                  fmt.Sprintf("%s-%s", oldStack.Name, ending),
		RestorationObjectName: oldStack.RestorationObjectName,
	}
}

func getLoneInstanceObject(obj interface{}, ending string, order int) state.Object {
	obj2 := obj.(rds.RestoreDBInstanceFromDBSnapshotInput)
	insID := fmt.Sprintf("%s-%s", *obj2.DBInstanceIdentifier, ending)
	obj2.DBInstanceIdentifier = &insID
	obj2.DBSnapshotIdentifier = &copySnapshotName
	b := state.EncodeRestoreDBInstanceFromDBSnapshotInput(&obj2)
	fn := helpers.RandomStateFileName()
	_, err := state.WriteOutput(*fn, b)
	if err != nil {

	}
	s := state.Object{
		FileName: *fn,
		Order:    order,
		ObjType:  state.LoneInstance,
	}
	return s
}
