package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
)

type KmsClient interface {
	CreateKey(ctx context.Context, params *kms.CreateKeyInput, optFns ...func(*kms.Options)) (*kms.CreateKeyOutput, error)
}

type KmsOperations struct {
	Client KmsClient
}

type KmsConfig struct {
	Description *string
	Multiregion *bool
	Policy      *string
}

func (k KmsOperations) CreateKMSKey(cfg *KmsConfig) (*types.KeyMetadata, error) {
	// TODO handle multiregion keys
	// TODO handle key policies
	input := &kms.CreateKeyInput{}
	if cfg != nil {
		input = handleKmsConfig(*cfg)
	}
	output, err := k.Client.CreateKey(context.TODO(), input)
	if err != nil {
		log.Printf("Error creating KMS key %s", err)
		return nil, err
	}
	return output.KeyMetadata, nil
}

func handleKmsConfig(k KmsConfig) *kms.CreateKeyInput {
	input := kms.CreateKeyInput{}
	if k.Description != nil {
		input.Description = k.Description
	}
	if k.Multiregion != nil {
		input.MultiRegion = k.Multiregion
	}
	if k.Policy != nil {
		input.Policy = k.Policy
	}
	return &input
}
