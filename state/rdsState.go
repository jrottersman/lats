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

// EncodeRDSClusterOutput takes a cluster snapshot and creates bytes
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

//GetRDSSnapshotOutput reads a snapshot at this point I think we can delete this
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

// CreateInstanceInput creates an instance to prep for creating our Cluster
func CreateDbInstanceInput(i *types.DBInstance, ci *string) *rds.CreateDBInstanceInput {
	return &rds.CreateDBInstanceInput{
		DBInstanceIdentifier: i.DBInstanceIdentifier,
		DBName:               i.DBName,
		Engine:               i.Engine,
		EngineVersion:        i.EngineVersion,
		DBInstanceClass:      i.DBInstanceClass,
		StorageType:          i.StorageType,
		DBClusterIdentifier:  ci,
		AllocatedStorage:     &i.AllocatedStorage,
	}
}

//EncodeCreateDBInstanceInput bytes buffer for create
func EncodeCreateDBInstanceInput(c *rds.CreateDBInstanceInput) bytes.Buffer {
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)
	err := enc.Encode(&c)
	if err != nil {
		log.Fatalf("Error encoding our database: %s", err)
	}
	return encoder
}

//DecodeCreateDBInstanceInput creates the instance from our bytes buffer when we want to replay
func DecodeCreateDBInstanceInput(b bytes.Buffer) *rds.CreateDBInstanceInput {
	var dbCluster rds.CreateDBInstanceInput
	dec := gob.NewDecoder(&b)
	err := dec.Decode(&dbCluster)
	if err != nil {
		log.Fatalf("Error decoding state for RDS Cluster: %s", err)
	}
	return &dbCluster
}

//EncodeParameterGroup takes a parameter group and turns it into a bytes.Buffer
func EncodeParameterGroup(p *types.DBParameterGroup) bytes.Buffer {
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)
	err := enc.Encode(&p)
	if err != nil {
		log.Fatalf("Error encoding our parameter group: %s", err)
	}
	return encoder
}

//DecodeParameterGroup creates the instance from our bytes buffer when we want to replay
func DecodeParameterGroup(b bytes.Buffer) *types.DBParameterGroup {
	var pg types.DBParameterGroup
	dec := gob.NewDecoder(&b)
	err := dec.Decode(&pg)
	if err != nil {
		log.Fatalf("Error decoding state for parameter group: %s", err)
	}
	return &pg
}

//EncodeClusterParameterGroup takes a parameter group and turns it into a bytes.Buffer
func EncodeClusterParameterGroup(p *types.DBClusterParameterGroup) bytes.Buffer {
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)
	err := enc.Encode(&p)
	if err != nil {
		log.Fatalf("Error encoding our parameter group: %s", err)
	}
	return encoder
}

//DecodeClusterParameterGroup creates the instance from our bytes buffer when we want to replay
func DecodeClusterParameterGroup(b bytes.Buffer) *types.DBClusterParameterGroup {
	var pg types.DBClusterParameterGroup
	dec := gob.NewDecoder(&b)
	err := dec.Decode(&pg)
	if err != nil {
		log.Fatalf("Error decoding state for cluster parameter group: %s", err)
	}
	return &pg
}

//EncodeParameters encodes the parameters for a paramter group and turns it into a bytes.Buffer
func EncodeParameters(p *[]types.Parameter) bytes.Buffer {
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)
	err := enc.Encode(&p)
	if err != nil {
		log.Fatalf("Error encoding our parameters %s", err)
	}
	return encoder
}

func DecodeParameters(b bytes.Buffer) *[]types.Parameter {
	var lp []types.Parameter
	dec := gob.NewDecoder(&b)
	err := dec.Decode(&lp)
	if err != nil {
		log.Fatalf("Error decoding parameters %s", err)
	}
	return &lp
}
