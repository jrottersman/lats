package state

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/jrottersman/lats/helpers"
)

// EncodeRDSDatabaseOutput converts a dbInstace to an array of bytes in preperation for wrtiing it to disk
func EncodeRDSDatabaseOutput(db *types.DBInstance) bytes.Buffer {
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)

	err := enc.Encode(db)
	if err != nil {
		log.Fatalf("Error encoding our database: %s", err)
	}
	return encoder
}

// DecodeRDSClusterOutput takes a bytes buffer and returns it to a DbCluster type in preperation of restoring the database
func DecodeRDSClusterOutput(b bytes.Buffer) types.DBCluster {
	var dbCluster types.DBCluster
	dec := gob.NewDecoder(&b)
	err := dec.Decode(&dbCluster)
	if err != nil {
		log.Fatalf("Error decoding state for RDS Cluster: %s", err)
	}
	return dbCluster
}

func EncodeRDSClusterOutput(db *types.DBCluster) bytes.Buffer {
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)

	err := enc.Encode(db)
	if err != nil {
		log.Fatalf("Error encoding our database: %s", err)
	}
	return encoder
}

// DecodeRDSDatabaseOutput takes a bytes buffer and returns it to a DbInstance type in preperation of restoring the database
func DecodeRDSDatabaseOutput(b bytes.Buffer) types.DBInstance {
	var dbInstance types.DBInstance
	dec := gob.NewDecoder(&b)
	err := dec.Decode(&dbInstance)
	if err != nil {
		log.Fatalf("Error decoding state for RDS Instance: %s", err)
	}
	return dbInstance
}

// EncodeRDSSnapshotOutput converts a DbSnapshot struct to an array of bytes in preperation for wrtiing it to disk
func EncodeRDSSnapshotOutput(snapshot *types.DBSnapshot) bytes.Buffer {
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)

	err := enc.Encode(snapshot)
	if err != nil {
		log.Fatalf("Error encoding our snapshot: %s", err)
	}
	return encoder
}

func GetRDSSnapshotOutput(s StateManager, snap string) (*types.DBSnapshot, error) {
	i := s.GetStateObject(snap)
	snapshot, ok := i.(types.DBSnapshot)
	if !ok {
		str := fmt.Sprintf("error decoding snapshot from interface %v", i)
		return nil, errors.New(str)
	}
	return &snapshot, nil
}

func EncodeRDSClusterSnapshotOutput(snapshot *types.DBClusterSnapshot) bytes.Buffer {
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)

	err := enc.Encode(snapshot)
	if err != nil {
		log.Fatalf("Error encoding our snapshot: %s", err)
	}
	return encoder
}

func GenerateRestoreDBInstanceFromDBSnapshotInput(r RDSRestorationStore) *rds.RestoreDBInstanceFromDBSnapshotInput {
	return &rds.RestoreDBInstanceFromDBSnapshotInput{
		DBInstanceClass:             r.GetInstanceClass(),
		DBInstanceIdentifier:        r.GetInstanceIdentifier(),
		AutoMinorVersionUpgrade:     r.GetAutoMinorVersionUpgrade(),
		AllocatedStorage:            r.GetAllocatedStorage(),
		BackupTarget:                r.GetBackupTarget(),
		DBSnapshotIdentifier:        r.GetSnapshotIdentifier(),
		DeletionProtection:          r.GetDeleteProtection(),
		EnableCloudwatchLogsExports: r.GetEnabledCloudwatchLogsExports(),
	}
}

func GenerateRestoreDBInstanceFromDBClusterSnapshotInput(r RDSRestorationStore) *rds.RestoreDBInstanceFromDBSnapshotInput {
	return &rds.RestoreDBInstanceFromDBSnapshotInput{
		DBInstanceClass:             r.GetInstanceClass(),
		DBInstanceIdentifier:        r.GetInstanceIdentifier(),
		AutoMinorVersionUpgrade:     r.GetAutoMinorVersionUpgrade(),
		AllocatedStorage:            r.GetAllocatedStorage(),
		BackupTarget:                r.GetBackupTarget(),
		DBClusterSnapshotIdentifier: r.GetClusterSnapshotIdentifier(),
		DeletionProtection:          r.GetDeleteProtection(),
		EnableCloudwatchLogsExports: r.GetEnabledCloudwatchLogsExports(),
	}
}

func EncodeRestoreDBInstanceFromDBSnapshotInput(r *rds.RestoreDBInstanceFromDBSnapshotInput) bytes.Buffer {
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)

	err := enc.Encode(r)
	if err != nil {
		log.Fatalf("Error encoding our snapshot: %s", err)
	}
	return encoder
}

func DecodeRestoreDBInstanceFromDBSnapshotInput(b bytes.Buffer) *rds.RestoreDBInstanceFromDBSnapshotInput {
	var Restore rds.RestoreDBInstanceFromDBSnapshotInput
	dec := gob.NewDecoder(&b)
	err := dec.Decode(&Restore)
	if err != nil {
		log.Fatalf("Error decoding state for RestoreDBInstance struct: %s", err)
	}
	return &Restore
}

func GenerateRestoreDBClusterFromSnapshotInput(r RDSRestorationStore) *rds.RestoreDBClusterFromSnapshotInput {
	return &rds.RestoreDBClusterFromSnapshotInput{
		DBClusterIdentifier: r.GetDBClusterIdentifier(),
		Engine:              r.GetClusterEngine(),
		SnapshotIdentifier:  r.GetClusterSnapshotIdentifier(),
	}
}

func EncodeRestoreDBClusterFromSnapshotInput(r *rds.RestoreDBClusterFromSnapshotInput) bytes.Buffer {
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)

	err := enc.Encode(r)
	if err != nil {
		log.Fatalf("Error encoding our snapshot: %s", err)
	}
	return encoder
}

func DecodeRestoreDBClusterFromSnapshotInput(b bytes.Buffer) *rds.RestoreDBClusterFromSnapshotInput {
	var Restore rds.RestoreDBClusterFromSnapshotInput
	dec := gob.NewDecoder(&b)
	err := dec.Decode(&Restore)
	if err != nil {
		log.Fatalf("Error decoding state for RestoreDBCluster struct: %s", err)
	}
	return &Restore
}

func DecodeRDSClusterSnapshotOutput(b bytes.Buffer) types.DBClusterSnapshot {
	var dbSnapshot types.DBClusterSnapshot
	dec := gob.NewDecoder(&b)
	err := dec.Decode(&dbSnapshot)
	if err != nil {
		log.Fatalf("Error decoding state for cluster snapshot: %s", err)
	}
	return dbSnapshot
}

func GetRDSClusterSnapshotOutput(s StateManager, snap string) (*types.DBClusterSnapshot, error) {
	i := s.GetStateObject(snap)
	snapshot, ok := i.(types.DBClusterSnapshot)
	if !ok {
		str := fmt.Sprintf("error decoding cluster snapshot from interface %v", i)
		return nil, errors.New(str)
	}
	return &snapshot, nil
}

func GetRDSDatabaseInstanceOutput(s StateManager, dbName string) (*types.DBInstance, error) {
	i := s.GetStateObject(dbName)
	dbi, ok := i.(types.DBInstance)
	if !ok {
		str := fmt.Sprintf("error decoding instance from interface %v", i)
		return nil, errors.New(str)
	}
	return &dbi, nil
}

func GetRDSDatabaseClusterOutput(s StateManager, dbName string) (*types.DBCluster, error) {
	i := s.GetStateObject(dbName)
	dbi, ok := i.(types.DBCluster)
	if !ok {
		str := fmt.Sprintf("error decoding cluster from interface %v", i)
		return nil, errors.New(str)
	}
	return &dbi, nil
}

func GetClusterInstances(s StateManager, clusterIdentifier string) ([]*types.DBInstance, error) {
	for _, v := range s.StateLocations {
		if v.ObjectType == RdsInstanceType {
			return nil, nil
		}
	}
	return nil, nil
}

// DecodeRDSSnapshhotOutput takes a bytes buffer and returns it to a DbSnapshot type in preperation of restoring the database
func DecodeRDSSnapshotOutput(b bytes.Buffer) types.DBSnapshot {
	var dbSnapshot types.DBSnapshot
	dec := gob.NewDecoder(&b)
	err := dec.Decode(&dbSnapshot)
	if err != nil {
		log.Fatalf("Error decoding state for snapshot: %s", err)
	}
	return dbSnapshot
}

func WriteOutput(filename string, b bytes.Buffer) (int64, error) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Error creating file: %s", err)
	}
	defer f.Close()
	n, err := b.WriteTo(f)
	if err != nil {
		log.Fatalf("error writing to file %s", err)
	}
	return n, err
}

func GenerateRDSClusterStack(r RDSRestorationStore, name string, fn *string) (*Stack, error) {
	if fn == nil {
		fn = helpers.RandomStateFileName()
	}

	ClusterInput := GenerateRestoreDBClusterFromSnapshotInput(r)

	// This is the cluster
	bc := EncodeRestoreDBClusterFromSnapshotInput(ClusterInput)
	_, err := WriteOutput(*fn, bc)
	if err != nil {
		return nil, err
	}
	clusterObj := NewObject(*fn, 1, Cluster)
	var firstObjects []Object
	firstObjects = append(firstObjects, clusterObj)

	// TODO figure out how to handle the instances
	return nil, nil
}

func ClusterInstancesToObjects(t *types.DBCluster) ([]Object, error) {
	return nil, nil
}

func GenerateRDSInstanceStack(r RDSRestorationStore, name string, fn *string) (*Stack, error) {
	if fn == nil {
		fn = helpers.RandomStateFileName()
	}

	DBInput := GenerateRestoreDBInstanceFromDBSnapshotInput(r)

	b := EncodeRestoreDBInstanceFromDBSnapshotInput(DBInput)
	_, err := WriteOutput(*fn, b)
	if err != nil {
		return nil, err
	}

	obj := NewObject(*fn, 1, LoneInstance) // 1 is the order currently we just have the instance so this is 1 we will have to update it once we are handling parameter groups

	var objects []Object
	objects = append(objects, obj)

	m := make(map[int][]Object)
	m[1] = objects

	return &Stack{
		Name:                  name,
		RestorationObjectName: LoneInstance,
		Objects:               m,
	}, nil
}
