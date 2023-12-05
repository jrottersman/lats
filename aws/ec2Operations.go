package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type Ec2Client interface {
	CreateSecurityGroup(ctx context.Context, params *ec2.CreateSecurityGroupInput) (*ec2.CreateSecurityGroupOutput, error)
}
