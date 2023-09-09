package state

import (
	"encoding/gob"
	"os"
	"sync"
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
	if *result.KeyId != *kmd.KeyId {
		t.Errorf("got %s expected %s", *result.KeyId, *kmd.KeyId)
	}
}

func TestGetKmsKeyOutput(t *testing.T) {
	filename := "/tmp/foo"
	kmd := types.KeyMetadata{
		KeyId: aws.String("foo"),
	}

	defer os.Remove(filename)
	r := EncodeKmsOutput(&kmd)
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
	newKmd, err := GetKmsOutput(&sm, "foo")
	if err != nil {
		t.Errorf("got error: %s", err)
	}

	if *newKmd.KeyId != "foo" {
		t.Errorf("got %s expected foo", *newKmd.KeyId)
	}
}
