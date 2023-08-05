package state

import (
	"encoding/gob"
	"testing"

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
