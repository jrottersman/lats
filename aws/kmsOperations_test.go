package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
)

type mockKMSClient struct{}

func (m mockKMSClient) CreateKey(ctx context.Context, params *kms.CreateKeyInput, optFns ...func(*kms.Options)) (*kms.CreateKeyOutput, error) {
	kid := "foobar"
	r := kms.CreateKeyOutput{
		KeyMetadata: &types.KeyMetadata{
			KeyId: aws.String(kid),
		},
	}
	return &r, nil
}

func TestCreateKMSKey(t *testing.T) {
	c := mockKMSClient{}
	kmsOp := KmsOperations{
		Client: c,
	}
	results, err := kmsOp.CreateKMSKey(nil)
	if err != nil {
		t.Errorf("error calling CreateKMSKey %s", err)
	}
	if *results.KeyId != "foobar" {
		t.Errorf("got %s expected foobar", *results.KeyId)
	}
}

func TestHandleKmsConfig(t *testing.T) {
	desc := "foobarbaz"
	policy := "iampolicy"
	mr := false
	k := KmsConfig{
		Description: &desc,
		Policy:      &policy,
		Multiregion: &mr,
	}
	keyInput := handleKmsConfig(k)

	if *keyInput.Description != desc {
		t.Errorf("expected %s got %s", desc, *keyInput.Description)
	}
	if *keyInput.MultiRegion != false {
		t.Errorf("multiregion should be false")
	}
	if *keyInput.Policy != policy {
		t.Errorf("policy should be %s got %s", policy, *keyInput.Policy)
	}
}
