package state

import (
	"bytes"
	"encoding/gob"
	"os"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
)

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
	var s []stateKV
	kv := stateKV{
		Object:       "foo",
		FileLocation: filename,
		ObjectType:   "RDSSnapshot",
	}
	s = append(s, kv)
	sm := StateManager{
		mu,
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
