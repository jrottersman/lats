package state

import (
	"bytes"
	"log"
	"encoding/gob"

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