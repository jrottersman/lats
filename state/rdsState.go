package state

import (
	"bytes"
	"log"
	"encoding/gob"

	"github.com/aws/aws-sdk-go-v2/service/rds/types"
)

func EncodeRDSDatabaseOutput(db *types.DBInstance) bytes.Buffer {
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)

	err := enc.Encode(db)
	if err != nil {
		log.Printf("Error encoding our database: %s", err)
	}
	return encoder
}