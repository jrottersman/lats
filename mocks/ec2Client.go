package mock

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type MockEC2Client struct {
}

func (m MockEC2Client) CreateSecurityGroup(ctx context.Context, params *ec2.CreateSecurityGroupInput) (*ec2.CreateSecurityGroupOutput, error) {
	return &ec2.CreateSecurityGroupOutput{
		GroupId: aws.String("foobar"),
	}, nil
}
