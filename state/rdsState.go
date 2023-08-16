package state

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"os"

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

func GetRDSDatabaseInstanceOutput(s StateManager, dbName string) (*types.DBInstance, error) {
	i := s.GetStateObject(dbName)
	dbi, ok := i.(types.DBInstance)
	if !ok {
		str := fmt.Sprintf("error decoding instance from interface %v", i)
		return nil, errors.New(str)
	}
	return &dbi, nil
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
