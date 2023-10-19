package state

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/service/kms/types"
)

// EncodeKmsOutput encodes the output of KMS into a bytes.Buffer for writing
func EncodeKmsOutput(kmd *types.KeyMetadata) bytes.Buffer {
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)

	err := enc.Encode(kmd)
	if err != nil {
		slog.Error("Error encoding our database:", "Error", err)
	}
	return encoder
}

// DecodeKmsOutput takes bytes and turns them into KeyMetadata
func DecodeKmsOutput(b bytes.Buffer) types.KeyMetadata {
	var kmsMetadata types.KeyMetadata
	dec := gob.NewDecoder(&b)
	err := dec.Decode(&kmsMetadata)
	if err != nil {
		slog.Error("Error decoding state for KMS Key:", "error", err)
	}
	return kmsMetadata
}

// GetKmsOutput takes a keyid and get's the key we can probably delete this one
func GetKmsOutput(s *StateManager, KeyID string) (*types.KeyMetadata, error) {
	i := s.GetStateObject(KeyID)
	key, ok := i.(types.KeyMetadata)
	if !ok {
		str := fmt.Sprintf("error decoding KMS key from interface %v", i)
		return nil, errors.New(str)
	}
	return &key, nil
}
