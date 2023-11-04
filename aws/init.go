package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

// Init creates an RDS Client
func Init(region string) DbInstances {
	cfg := createConfig(region)
	client := getRDSClient(cfg)
	return DbInstances{
		RdsClient: client,
	}
}

// InitKms creates a KMS client
func InitKms(region string) KmsOperations {
	cfg := createConfig(region)
	client := getKMSClient(cfg)
	return KmsOperations{
		Client: client,
	}
}

func getRDSClient(cfg aws.Config) *rds.Client {
	return rds.NewFromConfig(cfg)
}

func getKMSClient(cfg aws.Config) *kms.Client {
	return kms.NewFromConfig(cfg)
}

func createConfig(region string) aws.Config {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		panic("configuration error, " + err.Error())
	}
	return cfg
}
