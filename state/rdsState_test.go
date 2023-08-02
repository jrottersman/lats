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
