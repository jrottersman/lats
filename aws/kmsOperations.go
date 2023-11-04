package aws

import (
	"context"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
)

// KmsClient type for mocks
type KmsClient interface {
	CreateKey(ctx context.Context, params *kms.CreateKeyInput, optFns ...func(*kms.Options)) (*kms.CreateKeyOutput, error)
}

// KmsOperations struct with the KmsClient
type KmsOperations struct {
	Client KmsClient
}

// KmsConfig descripes our KmsConfig and they way it works
type KmsConfig struct {
	Description *string
	Multiregion *bool
	Policy      *string
}

// CreateKMSKey creates a new KMS key for multiregion deploys
func (k KmsOperations) CreateKMSKey(cfg *KmsConfig) (*types.KeyMetadata, error) {
	input := &kms.CreateKeyInput{}
	if cfg != nil {
		input = handleKmsConfig(*cfg)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	output, err := k.Client.CreateKey(ctx, input)
	if err != nil {
		slog.Warn("Error creating KMS key", "error", err)
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
