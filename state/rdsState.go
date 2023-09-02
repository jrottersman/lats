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
	// id := instance.DBInstanceClass //have this gettewrr
	// allocatedStorage := snapshot.AllocatedStorage // have this getter
	// AutoMinorVersionUpgrade := instance.AutoMinorVersionUpgrade
	// backupTarget := instance.BackupTarget
	// instanceClass := instance.DBInstanceClass
	// dbSnapshotId := snapshot.DBSnapshotIdentifier
	// deleteProtection := instance.DeletionProtection
	// cloudwatchLogs := instance.EnabledCloudwatchLogsExports
	return &rds.RestoreDBInstanceFromDBSnapshotInput{}
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

func CheckSnasphotType(snapshot interface{}) {
	switch snapshot.(type) {
	case types.DBClusterSnapshot:
		GenerateRDSClusterStack()
	case types.DBSnapshot:
		GenerateRDSInstanceStack()
	}
}

func GenerateRDSClusterStack() {
	fmt.Println("TODO: implement cluster type")
}

func GenerateRDSInstanceStack() {
	fmt.Println("TODO: Implement instance type")
}
