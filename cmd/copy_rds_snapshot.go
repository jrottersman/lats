package cmd

import (
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/jrottersman/lats/aws"
	"github.com/jrottersman/lats/helpers"
	"github.com/jrottersman/lats/stack"
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
	CopyRDSSnapshotCmd.Flags().StringVarP(&copySnapshotName, "new-snapshot", "c", "", "Name of the snapshot copy we are creating")
	CopyRDSSnapshotCmd.Flags().StringVarP(&originalSnapshotName, "snapshot", "s", "", "Snapshot we want to copy")
	CopyRDSSnapshotCmd.Flags().StringVarP(&configFile, "config-file", "f", "", "Config file for the snapshot that we want to parse")
}

func copySnapshot() {
	config, err := readConfig(".latsConfig.json")
	if err != nil {
		slog.Error("Error reading config", "error", err)
	}
	stateFileName := config.StateFileName
	sm, err := state.ReadState(stateFileName)
	if err != nil {
		slog.Error("Error reading state", "error", err)
	}

	if configFile != "" {
		viper.SetConfigFile(configFile)
		viper.AddConfigPath(".")
		err := viper.ReadInConfig() // Find and read the config file
		if err != nil {             // Handle errors reading the config file
			slog.Error("Viper Error reading config", "error", err)
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
		slog.Info("creating KMS key")
		kmsKey = createKMSKey(config, sm)
	}

	// Get RDS Client
	dbi := aws.Init(config.BackupRegion)
	dbi2 := aws.Init(config.MainRegion)

	// Copy Snapshot
	origStack, err := FindStack(sm, originalSnapshotName)
	if err != nil {
		slog.Error("Error finding stack", "error", err)
	}

	if origStack.RestorationObjectName == stack.Cluster {
		slog.Info("copying cluster snapshot")
		arn, err := dbi2.GetSnapshotARN(originalSnapshotName, true)
		if err != nil {
			slog.Error("Couldn't find snapshot ", "snapshot", originalSnapshotName)
		}
		_, err = dbi.CopyClusterSnaphot(*arn, copySnapshotName, config.MainRegion, kmsKey)
		if err != nil {
			slog.Error("Couldn't copy snapshot ", "error", err)
		}
	}
	if origStack.RestorationObjectName == stack.LoneInstance {
		slog.Info("copying instance snapshot")
		iarn, err := dbi2.GetSnapshotARN(originalSnapshotName, false)
		if err != nil {
			slog.Error("Couldn't find snapshot ", "snapshot", originalSnapshotName)
		}
		_, err = dbi.CopySnapshot(*iarn, copySnapshotName, config.MainRegion, kmsKey)
		if err != nil {
			slog.Error("Couldn't copy snapshot ", "error", err)
		}
	}
	stack := NewStack(*origStack, copySnapshotName)

	fn := fmt.Sprintf(".state/%s", *helpers.RandomStateFileName())
	err = stack.Write(fn)
	if err != nil {
		slog.Error("Couldn't write stack", "error", err)
	}

	sm.UpdateState(stack.Name, fn, "stack")
	sm.SyncState(stateFileName)
}

func createKMSKey(config Config, sm state.StateManager) string {
	var kmsStruct *types.KeyMetadata
	c := aws.InitKms(config.BackupRegion)
	kmsStruct, err := c.CreateKMSKey(nil)
	if err != nil {
		slog.Error("failed creating KMS key", "error", err)
	}
	return *kmsStruct.KeyId

}

//FindStack get's a stack for creating our new stack when we copy the snapshot
func FindStack(sm state.StateManager, snapshot string) (*stack.Stack, error) {
	sm.Mu.Lock()
	defer sm.Mu.Unlock()
	if len(sm.StateLocations) == 0 {
		slog.Error("There are no state locations create a snapshot with Lats before attempting to copy it")
	}
	slog.Info("state locations", "locations", sm.StateLocations)
	for _, v := range sm.StateLocations {
		slog.Info("Looping through objects this one is ", "name", v.Object)
		if v.ObjectType != "stack" {
			continue
		}
		stack, err := stack.ReadStack(v.FileLocation)
		if err != nil {
			slog.Error("error reading stack", "error", err)
			return nil, err
		}
		if stack == nil {
			slog.Error("Read stack returned a nil value")
		}
		slog.Info("stack name is", "name", stack.Name, "snapshot", snapshot)
		if stack.Name == snapshot {
			return stack, nil
		}
	}
	slog.Error("Returning a nil stack")
	return nil, nil
}

// NewStack generates the new stack that we are going to use
func NewStack(oldStack stack.Stack, name string) *stack.Stack {
	objs := make(map[int][]stack.Object)
	for k, v := range oldStack.Objects {
		objs[k] = []stack.Object{}
		for _, i := range v {
			obj := i.ReadObject()
			switch i.ObjType {
			case stack.LoneInstance:
				slog.Info("Generating lone instance object")
				s := getLoneInstanceObject(obj, name, k)
				objs[k] = append(objs[k], s)
			case stack.Cluster:
				s := getClusterObject(obj, name, k)
				objs[k] = append(objs[k], s)
			case stack.Instance:
				s := getInstanceObject(obj, name, k)
				objs[k] = append(objs[k], s)
			case stack.DBClusterParameterGroup:
				objs[k] = append(objs[k], i)
			case stack.DBParameterGroup:
				objs[k] = append(objs[k], i)
			}
		}
	}
	return &stack.Stack{
		Name:                  fmt.Sprintf("%s", name),
		RestorationObjectName: oldStack.RestorationObjectName,
	}
}

func getLoneInstanceObject(obj interface{}, name string, order int) stack.Object {
	obj2 := obj.(*rds.RestoreDBInstanceFromDBSnapshotInput)
	insID := fmt.Sprintf("%s-instance", name)
	obj2.DBInstanceIdentifier = &insID
	obj2.DBSnapshotIdentifier = &copySnapshotName
	obj2.AvailabilityZone = nil
	obj2.DBParameterGroupName = nil
	obj2.DBSubnetGroupName = nil
	obj2.VpcSecurityGroupIds = nil
	obj2.OptionGroupName = nil
	b := state.EncodeRestoreDBInstanceFromDBSnapshotInput(obj2)
	fn := fmt.Sprintf(".state/%s", *helpers.RandomStateFileName())
	_, err := state.WriteOutput(fn, b)
	if err != nil {
		slog.Error("Error writing ouptut", "error", err)
	}
	s := stack.Object{
		FileName: fn,
		Order:    order,
		ObjType:  stack.LoneInstance,
	}
	return s
}

func getClusterObject(obj interface{}, name string, order int) stack.Object {
	obj2 := obj.(*rds.RestoreDBClusterFromSnapshotInput)
	clsID := fmt.Sprintf("%s", name)
	obj2.DBClusterIdentifier = &clsID
	obj2.SnapshotIdentifier = &copySnapshotName
	obj2.AvailabilityZones = nil
	obj2.DBClusterParameterGroupName = nil
	obj2.DBSubnetGroupName = nil
	obj2.KmsKeyId = nil
	obj2.VpcSecurityGroupIds = nil
	b := state.EncodeRestoreDBClusterFromSnapshotInput(obj2)
	fn := fmt.Sprintf(".state/%s", *helpers.RandomStateFileName())
	_, err := state.WriteOutput(fn, b)
	if err != nil {
		slog.Error("Error writing output", "Error", err)
	}
	return stack.Object{
		FileName: fn,
		Order:    order,
		ObjType:  stack.Cluster,
	}
}

func getInstanceObject(obj interface{}, ending string, order int) stack.Object {
	obj2 := obj.(rds.CreateDBInstanceInput)
	insID := fmt.Sprintf("%s-%s", *obj2.DBInstanceIdentifier, ending)
	clusterID := fmt.Sprintf("%s-%s", *obj2.DBClusterIdentifier, ending)
	obj2.DBInstanceIdentifier = &insID
	obj2.DBClusterIdentifier = &clusterID
	obj2.AvailabilityZone = nil
	obj2.DBParameterGroupName = nil
	obj2.DBSubnetGroupName = nil
	b := state.EncodeCreateDBInstanceInput(&obj2)
	fn := fmt.Sprintf(".state/%s", *helpers.RandomStateFileName())
	_, err := state.WriteOutput(fn, b)
	if err != nil {
		slog.Error("Error writing output", "Error", err)
	}
	return stack.Object{
		FileName: fn,
		Order:    order,
		ObjType:  stack.Instance,
	}
}
