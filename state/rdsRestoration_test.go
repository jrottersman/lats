package state

import (
	"os"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
)

func TestRDSRestorationStoreBuilder(t *testing.T) {
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

	filename2 := "/tmp/foo2"
	dbi := types.DBInstance{
		AllocatedStorage:     1000,
		DBInstanceIdentifier: aws.String("foobar"),
	}

	defer os.Remove(filename2)
	r2 := EncodeRDSDatabaseOutput(&dbi)
	_, err = WriteOutput(filename2, r2)
	if err != nil {
		t.Errorf("error writing file %s", err)
	}
	kv2 := StateKV{
		Object:       "foobar",
		FileLocation: filename2,
		ObjectType:   RdsInstanceType,
	}
	s = append(s, kv2)
	sm := StateManager{
		mu,
		s,
	}

	resp, err := RDSRestorationStoreBuilder(sm, "foo")
	if err != nil {
		t.Errorf("error writing file %s", err)
	}

	if resp.Snapshot.AllocatedStorage != 1000 {
		t.Errorf("oops WTF wanted 1000 got %d", resp.Snapshot.AllocatedStorage)
	}
}

func TestGetNilStoreInstanceIdentifier(t *testing.T) {
	NilStore := RDSRestorationStore{}
	s := NilStore.GetInstanceIdentifier()
	if s != nil {
		t.Errorf("s should be nil it is %v", s)
	}
}

func TestGetNilInstanceIdentifier(t *testing.T) {
	store := RDSRestorationStore{
		Instance: &types.DBInstance{},
	}
	s := store.GetInstanceIdentifier()
	if s != nil {
		t.Errorf("s should be nil it is %v", s)
	}
}

func TestGetValueInstanceIdentifier(t *testing.T) {
	i := types.DBInstance{
		DBInstanceIdentifier: aws.String("foo"),
	}
	store := RDSRestorationStore{
		Instance: &i,
	}
	s := store.GetInstanceIdentifier()

	if *s != "foo" {
		t.Errorf("expected foo got %s", *s)
	}
}

func TestGetNilStoreInstanceClass(t *testing.T) {
	NilStore := RDSRestorationStore{}
	s := NilStore.GetInstanceClass()
	if s != nil {
		t.Errorf("s should be nil it is %v", s)
	}
}

func TestGetNilInstanceClass(t *testing.T) {
	store := RDSRestorationStore{
		Instance: &types.DBInstance{},
	}
	s := store.GetInstanceClass()
	if s != nil {
		t.Errorf("s should be nil it is %v", s)
	}
}

func TestGetValueInstanceClass(t *testing.T) {
	i := types.DBInstance{
		DBInstanceClass: aws.String("t3.micro"),
	}
	store := RDSRestorationStore{
		Instance: &i,
	}
	s := store.GetInstanceClass()

	if *s != "t3.micro" {
		t.Errorf("expected t3.micro got %s", *s)
	}
}
