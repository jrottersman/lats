package state

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
)

func TestInitState(t *testing.T) {
	filename := "/tmp/.stateConf.json"
	defer os.Remove(filename)

	err := InitState(filename)
	if err != nil {
		t.Errorf("error creating state file %s", err)
	}
}

func TestReadState(t *testing.T) {
	filename := "/tmp/.stateConf.json"
	defer os.Remove(filename)

	err := InitState(filename)
	if err != nil {
		t.Errorf("error creating state file %s", err)
	}
	s, err := ReadState(filename)
	if err != nil {
		t.Errorf("error creating state in memory: %s", err)
	}
	if len(s.StateLocations) > 0 {
		t.Errorf("expected length to be 0")
	}

}

func TestUpdateState(t *testing.T) {
	var mu sync.Mutex
	var s []stateKV
	sm := StateManager{
		mu,
		s,
	}
	filename := "/tmp/foo"
	obj := "boo"
	ot := "bar"
	sm.UpdateState(obj, filename, ot)
	if sm.StateLocations[0].Object != "boo" {
		t.Errorf("got %s expected %s\n", sm.StateLocations[0].Object, obj)
	}
}

func TestSyncState(t *testing.T) {
	var mu sync.Mutex
	var s []stateKV
	sm := StateManager{
		mu,
		s,
	}

	filename := "/tmp/foo"
	obj := "boo"
	ot := "baz"
	defer os.Remove(filename)
	sm.UpdateState(obj, filename, ot)
	sm.SyncState(filename)
	f, err := os.ReadFile(filename)
	if err != nil {
		t.Errorf("couldn't read file got %s", err)
	}
	var sf []stateKV
	err = json.Unmarshal(f, &sf)
	if err != nil {
		fmt.Printf("Couldn't unmarshall the json %s", err)
	}
	if sf[0].Object != obj {
		t.Errorf("got %s expected %s", sf[0].Object, obj)
	}

}

func TestGetStateObject(t *testing.T) {
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

	result := sm.GetStateObject("foo")
	if result == nil {
		t.Errorf("result is nil")
	}
	res, ok := result.(types.DBSnapshot)
	if !ok {
		t.Errorf("issue converting struct")
	}
	if *res.DBInstanceIdentifier != "foobar" {
		t.Errorf("got %s expected foobar", *res.DBInstanceIdentifier)
	}
}

func TestGetStateObjectInstance(t *testing.T) {
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
	var s []stateKV
	kv := stateKV{
		Object:       "foo",
		FileLocation: filename,
		ObjectType:   RdsInstanceType,
	}
	s = append(s, kv)
	sm := StateManager{
		mu,
		s,
	}

	result := sm.GetStateObject("foo")
	if result == nil {
		t.Errorf("result is nil")
	}
	res, ok := result.(types.DBInstance)
	if !ok {
		t.Errorf("issue converting struct")
	}
	if *res.DBInstanceIdentifier != "foobar" {
		t.Errorf("got %s expected foobar", *res.DBInstanceIdentifier)
	}
}
