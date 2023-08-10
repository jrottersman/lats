package state

import (
	"bytes"
	"encoding/gob"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/kms/types"
)

// EncodeKmsOutput encodes the output of KMS into a bytes.Buffer for writing
func EncodeKmsOutput(kmd *types.KeyMetadata) bytes.Buffer {
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)

	err := enc.Encode(kmd)
	if err != nil {
		log.Fatalf("Error encoding our database: %s", err)
	}
	return encoder
}

func DecodeKmsOutput(b bytes.Buffer) types.KeyMetadata {
	var kmsMetadata types.KeyMetadata
	dec := gob.NewDecoder(&b)
	err := dec.Decode(&kmsMetadata)
	if err != nil {
		log.Fatalf("Error decoding state for KMS Key: %s", err)
	}
	return kmsMetadata
}
