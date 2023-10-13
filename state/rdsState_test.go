package state

import (
	"bytes"
	"encoding/gob"
	"os"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
)

func TestEncodeOptionGroup(t *testing.T) {

	og := types.OptionGroup{
		EngineName: aws.String("foo"),
	}
	r := EncodeOptionGroup(&og)
	var result types.OptionGroup
	dec := gob.NewDecoder(&r)
	err := dec.Decode(&result)
	if err != nil {
		t.Errorf("decode error: %s", err)
	}
	if *result.EngineName != *og.EngineName {
		t.Errorf("got %s expected %s", *result.EngineName, *og.EngineName)
	}
}
func TestEncodeRestoreDBInstanceFromDBSnapshotInput(t *testing.T) {

	db := rds.RestoreDBInstanceFromDBSnapshotInput{
		DBInstanceIdentifier: aws.String("foo"),
	}
	r := EncodeRestoreDBInstanceFromDBSnapshotInput(&db)
	var result rds.RestoreDBInstanceFromDBSnapshotInput
	dec := gob.NewDecoder(&r)
	err := dec.Decode(&result)
	if err != nil {
		t.Errorf("decode error: %s", err)
	}
	if *result.DBInstanceIdentifier != *db.DBInstanceIdentifier {
		t.Errorf("got %s expected %s", *result.DBInstanceIdentifier, *db.DBInstanceIdentifier)
	}
}

func TestDecodeRestoreDBInstanceFromDBSnapshotInput(t *testing.T) {
	db := rds.RestoreDBInstanceFromDBSnapshotInput{
		DBInstanceIdentifier: aws.String("foo"),
	}
	r := EncodeRestoreDBInstanceFromDBSnapshotInput(&db)
	resp := DecodeRestoreDBInstanceFromDBSnapshotInput((r))
	if *resp.DBInstanceIdentifier != *db.DBInstanceIdentifier {
		t.Errorf("got %s expected %s", *resp.DBInstanceIdentifier, *db.DBInstanceIdentifier)
	}
}

func TestEncodeRestoreDBClusterFromSnapshotInput(t *testing.T) {

	db := rds.RestoreDBClusterFromSnapshotInput{
		DBClusterIdentifier: aws.String("foo"),
	}
	r := EncodeRestoreDBClusterFromSnapshotInput(&db)
	var result rds.RestoreDBClusterFromSnapshotInput
	dec := gob.NewDecoder(&r)
	err := dec.Decode(&result)
	if err != nil {
		t.Errorf("decode error: %s", err)
	}
	if *result.DBClusterIdentifier != *db.DBClusterIdentifier {
		t.Errorf("got %s expected %s", *result.DBClusterIdentifier, *db.DBClusterIdentifier)
	}
}

func TestDecodeRestoreDBClusterFromSnapshotInput(t *testing.T) {
	db := rds.RestoreDBClusterFromSnapshotInput{
		DBClusterIdentifier: aws.String("foo"),
	}
	r := EncodeRestoreDBClusterFromSnapshotInput(&db)
	resp := DecodeRestoreDBClusterFromSnapshotInput((r))
	if *resp.DBClusterIdentifier != *db.DBClusterIdentifier {
		t.Errorf("got %s expected %s", *resp.DBClusterIdentifier, *db.DBClusterIdentifier)
	}
}

func TestEncodeRDSDBOutput(t *testing.T) {

	db := types.DBInstance{
		AllocatedStorage:      1000,
		BackupRetentionPeriod: 30,
	}
	r := EncodeRDSDatabaseOutput(&db)
	var result types.DBInstance
	dec := gob.NewDecoder(&r)
	err := dec.Decode(&result)
	if err != nil {
		t.Errorf("decode error: %s", err)
	}
	if result.AllocatedStorage != db.AllocatedStorage {
		t.Errorf("got %d expected %d", result.AllocatedStorage, db.AllocatedStorage)
	}
}

func TestDecodeRDSDBOutput(t *testing.T) {
	db := types.DBInstance{
		AllocatedStorage:      1000,
		BackupRetentionPeriod: 30,
	}
	r := EncodeRDSDatabaseOutput(&db)
	resp := DecodeRDSDatabaseOutput(r)
	if resp.AllocatedStorage != db.AllocatedStorage {
		t.Errorf("Expected %d, got %d", resp.AllocatedStorage, db.AllocatedStorage)
	}

}

func TestEncodeRDSClusterOutput(t *testing.T) {
	var storage int32 = 1000
	var retention int32 = 30

	db := types.DBCluster{
		AllocatedStorage:      &storage,
		BackupRetentionPeriod: &retention,
	}
	r := EncodeRDSClusterOutput(&db)
	var result types.DBCluster
	dec := gob.NewDecoder(&r)
	err := dec.Decode(&result)
	if err != nil {
		t.Errorf("decode error: %s", err)
	}
	if *result.AllocatedStorage != storage {
		t.Errorf("got %d expected %d", result.AllocatedStorage, storage)
	}
}

func TestDecodeRDSClusterOutput(t *testing.T) {
	var storage int32 = 1000
	var retention int32 = 30

	db := types.DBCluster{
		AllocatedStorage:      &storage,
		BackupRetentionPeriod: &retention,
	}
	r := EncodeRDSClusterOutput(&db)
	resp := DecodeRDSClusterOutput(r)
	if *resp.AllocatedStorage != storage {
		t.Errorf("Expected %d, got %d", *resp.AllocatedStorage, storage)
	}

}

func TestEncodeRDSSnapshotOutput(t *testing.T) {

	snap := types.DBSnapshot{
		AllocatedStorage: 1000,
		Encrypted:        true,
		PercentProgress:  100,
	}
	r := EncodeRDSSnapshotOutput(&snap)
	var result types.DBSnapshot
	dec := gob.NewDecoder(&r)
	err := dec.Decode(&result)
	if err != nil {
		t.Errorf("decode error: %s", err)
	}
	if result.AllocatedStorage != snap.AllocatedStorage {
		t.Errorf("got %d expected %d", result.AllocatedStorage, snap.AllocatedStorage)
	}
}

func TestDecodeRDSSnapshotOutput(t *testing.T) {
	snap := types.DBSnapshot{
		AllocatedStorage: 1000,
		Encrypted:        true,
		PercentProgress:  100,
	}
	r := EncodeRDSSnapshotOutput(&snap)
	resp := DecodeRDSSnapshotOutput(r)
	if resp.AllocatedStorage != snap.AllocatedStorage {
		t.Errorf("Expected %d, got %d", resp.AllocatedStorage, snap.AllocatedStorage)
	}

}

func TestEncodeRDSClusterSnapshotOutput(t *testing.T) {

	snap := types.DBClusterSnapshot{
		AllocatedStorage: 1000,
		PercentProgress:  100,
	}
	r := EncodeRDSClusterSnapshotOutput(&snap)
	var result types.DBClusterSnapshot
	dec := gob.NewDecoder(&r)
	err := dec.Decode(&result)
	if err != nil {
		t.Errorf("decode error: %s", err)
	}
	if result.AllocatedStorage != snap.AllocatedStorage {
		t.Errorf("got %d expected %d", result.AllocatedStorage, snap.AllocatedStorage)
	}
}

func TestDecodeRDSClusterSnapshotOutput(t *testing.T) {
	snap := types.DBClusterSnapshot{
		AllocatedStorage: 1000,
		PercentProgress:  100,
	}
	r := EncodeRDSClusterSnapshotOutput(&snap)
	resp := DecodeRDSClusterSnapshotOutput(r)
	if resp.AllocatedStorage != snap.AllocatedStorage {
		t.Errorf("Expected %d, got %d", resp.AllocatedStorage, snap.AllocatedStorage)
	}

}

func TestWriteOutput(t *testing.T) {
	type foo struct {
		A string
	}
	s := foo{
		"bar",
	}
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)

	err := enc.Encode(s)
	if err != nil {
		t.Errorf("Error encoding our test: %s", err)
	}

	var expected int64
	expected = 33
	filename := "/tmp/foo.gob"
	defer os.Remove(filename)

	n, err := WriteOutput(filename, encoder)
	if err != nil {
		t.Errorf("got an error: %s", err)
	}
	if n != expected {
		t.Errorf("got: %d expected %d", n, expected)
	}
}

func TestGetRDSSnapshotOutput(t *testing.T) {
	filename := "/tmp/foo"
	snap := types.DBSnapshot{
		AllocatedStorage:     1000,
		Encrypted:            true,
		PercentProgress:      100,
		DBInstanceIdentifier: aws.String("foobar"),
	}

	defer os.Remove(filename)
	r := EncodeRDSSnapshotOutput(&snap)
	_, err := WriteOutput(filename, r)
	if err != nil {
		t.Errorf("error writing file %s", err)
	}

	var mu sync.Mutex
	var s []StateKV
	kv := StateKV{
		Object:       "foo",
		FileLocation: filename,
		ObjectType:   "RDSSnapshot",
	}
	s = append(s, kv)
	sm := StateManager{
		&mu,
		s,
	}
	newSnap, err := GetRDSSnapshotOutput(sm, "foo")
	if err != nil {
		t.Errorf("got error: %s", err)
	}

	if *&newSnap.AllocatedStorage != 1000 {
		t.Errorf("expected %d got 1000", *&newSnap.AllocatedStorage)
	}
}

func TestGetRDSClusterSnapshotOutput(t *testing.T) {
	filename := "/tmp/foo"
	snap := types.DBClusterSnapshot{
		AllocatedStorage:    1000,
		PercentProgress:     100,
		DBClusterIdentifier: aws.String("foobar"),
	}

	defer os.Remove(filename)
	r := EncodeRDSClusterSnapshotOutput(&snap)
	_, err := WriteOutput(filename, r)
	if err != nil {
		t.Errorf("error writing file %s", err)
	}

	var mu sync.Mutex
	var s []StateKV
	kv := StateKV{
		Object:       "foo",
		FileLocation: filename,
		ObjectType:   ClusterSnapshotType,
	}
	s = append(s, kv)
	sm := StateManager{
		&mu,
		s,
	}
	newSnap, err := GetRDSClusterSnapshotOutput(sm, "foo")
	if err != nil {
		t.Errorf("got error: %s", err)
	}

	if *&newSnap.AllocatedStorage != 1000 {
		t.Errorf("expected %d got 1000", *&newSnap.AllocatedStorage)
	}
}

func TestGetRDSInstanceOutput(t *testing.T) {
	filename := "/tmp/foo"
	dbi := types.DBInstance{
		AllocatedStorage:     1000,
		DBInstanceIdentifier: aws.String("foobar"),
	}

	defer os.Remove(filename)
	r := EncodeRDSDatabaseOutput(&dbi)
	_, err := WriteOutput(filename, r)
	if err != nil {
		t.Errorf("error writing file %s", err)
	}

	var mu sync.Mutex
	var s []StateKV
	kv := StateKV{
		Object:       "foo",
		FileLocation: filename,
		ObjectType:   RdsInstanceType,
	}
	s = append(s, kv)
	sm := StateManager{
		&mu,
		s,
	}
	newDbi, err := GetRDSDatabaseInstanceOutput(sm, "foo")
	if err != nil {
		t.Errorf("got error: %s", err)
	}

	if *&newDbi.AllocatedStorage != 1000 {
		t.Errorf("expected %d got 1000", *&newDbi.AllocatedStorage)
	}
}

func TestGetRDSClusterOutput(t *testing.T) {
	var storage int32 = 1000
	filename := "/tmp/foo"
	dbi := types.DBCluster{
		AllocatedStorage:    &storage,
		DBClusterIdentifier: aws.String("foobar"),
	}

	defer os.Remove(filename)
	r := EncodeRDSClusterOutput(&dbi)
	_, err := WriteOutput(filename, r)
	if err != nil {
		t.Errorf("error writing file %s", err)
	}

	var mu sync.Mutex
	var s []StateKV
	kv := StateKV{
		Object:       "foo",
		FileLocation: filename,
		ObjectType:   RdsClusterType,
	}
	s = append(s, kv)
	sm := StateManager{
		&mu,
		s,
	}
	newDbi, err := GetRDSDatabaseClusterOutput(sm, "foo")
	if err != nil {
		t.Errorf("got error: %s", err)
	}

	if *newDbi.AllocatedStorage != storage {
		t.Errorf("expected %d got 1000", *&newDbi.AllocatedStorage)
	}
}

func TestGenerateRestoreDBInstanceFromDBSnapshotInput(t *testing.T) {
	rStore := RDSRestorationStore{
		Snapshot: &types.DBSnapshot{DBSnapshotIdentifier: aws.String("boo")},
		Instance: &types.DBInstance{},
	}

	resp := GenerateRestoreDBInstanceFromDBSnapshotInput(rStore)
	if *resp.DBSnapshotIdentifier != "boo" {
		t.Errorf("expected boo got %s", *resp.DBSnapshotIdentifier)
	}
}

func TestGenerateRestoreDBInstanceFromDBClusterSnapshotInput(t *testing.T) {
	rStore := RDSRestorationStore{
		ClusterSnapshot: &types.DBClusterSnapshot{DBClusterSnapshotIdentifier: aws.String("boo")},
		Instance:        &types.DBInstance{},
	}

	resp := GenerateRestoreDBInstanceFromDBClusterSnapshotInput(rStore)
	if *resp.DBClusterSnapshotIdentifier != "boo" {
		t.Errorf("expected boo got %v", *resp.DBClusterSnapshotIdentifier)
	}
}

func Test_CreateInstanceInput(t *testing.T) {
	i := types.DBInstance{
		DBInstanceIdentifier: aws.String("foo"),
		DBName:               aws.String("bar"),
		Engine:               aws.String("mysql"),
		EngineVersion:        aws.String("8"),
		DBInstanceClass:      aws.String("baz"),
		StorageType:          aws.String("bat"),
		AllocatedStorage:     1000,
	}
	result := CreateDbInstanceInput(&i, nil)
	if *result.AllocatedStorage != 1000 {
		t.Errorf("got %d expected 1000", *result.AllocatedStorage)
	}
}

func Test_EncodeCreateDBInstanceInput(t *testing.T) {

	db := rds.CreateDBInstanceInput{
		DBInstanceIdentifier: aws.String("foo"),
	}
	r := EncodeCreateDBInstanceInput(&db)
	var result rds.CreateDBInstanceInput
	dec := gob.NewDecoder(&r)
	err := dec.Decode(&result)
	if err != nil {
		t.Errorf("decode error: %s", err)
	}
	if *result.DBInstanceIdentifier != *db.DBInstanceIdentifier {
		t.Errorf("got %s expected %s", *result.DBInstanceIdentifier, *db.DBInstanceIdentifier)
	}
}

func Test_DecodeCreateDBInstanceInput(t *testing.T) {

	db := rds.CreateDBInstanceInput{
		DBInstanceIdentifier: aws.String("foo"),
	}
	r := EncodeCreateDBInstanceInput(&db)
	resp := DecodeCreateDBInstanceInput((r))
	if *resp.DBInstanceIdentifier != *db.DBInstanceIdentifier {
		t.Errorf("got %s expected %s", *resp.DBInstanceIdentifier, *db.DBInstanceIdentifier)
	}
}
