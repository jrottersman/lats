package state

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	ktypes "github.com/aws/aws-sdk-go-v2/service/kms/types"
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
	var s []StateKV
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
	var s []StateKV
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
	var sf []StateKV
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
	var s []StateKV
	kv := StateKV{
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

func TestGetStateObjectClusterSnapshot(t *testing.T) {
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
		mu,
		s,
	}

	result := sm.GetStateObject("foo")
	if result == nil {
		t.Errorf("result is nil")
	}
	res, ok := result.(types.DBClusterSnapshot)
	if !ok {
		t.Errorf("issue converting struct")
	}
	if *res.DBClusterIdentifier != "foobar" {
		t.Errorf("got %s expected foobar", *res.DBClusterIdentifier)
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
	var s []StateKV
	kv := StateKV{
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

func TestGetStateObjectKMS(t *testing.T) {
	filename := "/tmp/foo"
	key := ktypes.KeyMetadata{
		KeyId: aws.String("foo"),
	}

	defer os.Remove(filename)
	r := EncodeKmsOutput(&key)
	_, err := WriteOutput(filename, r)
	if err != nil {
		t.Errorf("error writing file %s", err)
	}

	var mu sync.Mutex
	var s []StateKV
	kv := StateKV{
		Object:       "foo",
		FileLocation: filename,
		ObjectType:   KMSKeyType,
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
	res, ok := result.(ktypes.KeyMetadata)
	if !ok {
		t.Errorf("issue converting struct")
	}
	if *res.KeyId != "foo" {
		t.Errorf("got %s expected foo", *res.KeyId)
	}
}

func TestCreateStackLookUp(t *testing.T) {
	filename := "/tmp/foo"
	defer os.Remove(filename)

	s := Stack{
		Name: "bar",
	}

	resp := CreateStackLookUp(s, filename)
	if resp.File != filename {
		t.Errorf("expected %s got %s", filename, resp.File)
	}
}

func TestAppendStack(t *testing.T) {
	sf := StackFiles{
		Stacks: []StackLookup{},
	}
	sl := StackLookup{
		Name: "bar",
		File: "/tmp/bar",
	}
	sf.AppendStackLookup(sl)
	if len(sf.Stacks) != 1 {
		t.Errorf("length should be 1 is %d", len(sf.Stacks))
	}
}

func TestGetStack(t *testing.T) {
	filename := "/tmp/bar"
	stackName := "bar"
	defer os.Remove(filename)

	sf := StackFiles{
		Stacks: []StackLookup{},
	}
	sl := StackLookup{
		Name: stackName,
		File: filename,
	}
	sf.AppendStackLookup(sl)

	st := Stack{
		Name: stackName,
	}
	err := st.Write(filename)
	if err != nil {
		t.Errorf("err should be nil it's %s", err)
	}
	s, err := sf.GetStack(filename)
	if err != nil {
		t.Errorf("stack error should be nil it's %s", err)
	}
	if s.Name != stackName {
		t.Errorf("expected %s got %s", stackName, s.Name)
	}

}
