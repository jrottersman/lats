package state

import (
	"encoding/gob"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
)

func TestEncodeKmsOutput(t *testing.T) {
	kmd := types.KeyMetadata{
		KeyId: aws.String("foo"),
	}
	b := EncodeKmsOutput(&kmd)
	var result types.KeyMetadata
	dec := gob.NewDecoder(&b)
	err := dec.Decode(&result)
	if err != nil {
		t.Errorf("decode error: %s", err)
	}
	if *result.KeyId != *kmd.KeyId {
		t.Errorf("got %s expected %s", *result.KeyId, *kmd.KeyId)
	}
}

func TestDecodeKmsOutput(t *testing.T) {
	kmd := types.KeyMetadata{
		KeyId: aws.String("foo"),
	}
	b := EncodeKmsOutput(&kmd)
	result := DecodeKmsOutput(b)
	if *result.KeyId != *kmd.KeyID {
		t.Errorf("got %s expected %s", *result.KeyId, *kmd.KeyId)
	}
}
