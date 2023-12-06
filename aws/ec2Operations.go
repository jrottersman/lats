package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type Ec2Client interface {
	CreateSecurityGroup(ctx context.Context, params *ec2.CreateSecurityGroupInput) (*ec2.CreateSecurityGroupOutput, error)
}

type EC2Instances struct {
	Client Ec2Client
}

func (c *EC2Instances) CreateSG(description *string, groupName *string, vpcID *string) (*ec2.CreateSecurityGroupOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	params := ec2.CreateSecurityGroupInput{
		Description: description,
		GroupName:   groupName,
		VpcId:       vpcID,
	}

	output, err := c.Client.CreateSecurityGroup(ctx, &params)
	if err != nil {
		return nil, err
	}
	return output, nil
}
