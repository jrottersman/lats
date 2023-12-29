package aws

import (
	"context"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type Ec2Client interface {
	CreateSecurityGroup(ctx context.Context, params *ec2.CreateSecurityGroupInput, optFns ...func(*ec2.Options)) (*ec2.CreateSecurityGroupOutput, error)
	AuthorizeSecurityGroupEgress(ctx context.Context, params *ec2.AuthorizeSecurityGroupEgressInput, optFns ...func(*ec2.Options)) (*ec2.AuthorizeSecurityGroupEgressOutput, error)
	DescribeSecurityGroups(ctx context.Context, params *ec2.DescribeSecurityGroupsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeSecurityGroupsOutput, error)
	AuthorizeSecurityGroupIngress(ctx context.Context, params *ec2.AuthorizeSecurityGroupIngressInput, optFns ...func(*ec2.Options)) (*ec2.AuthorizeSecurityGroupIngressOutput, error)
}

type SGInput struct {
	SGId          *string
	IpPermissions []types.IpPermission
}

type EC2Instances struct {
	Client Ec2Client
}

func (c *EC2Instances) CreateSG(description *string, groupName *string, vpcID *string, groupId *string) (*ec2.CreateSecurityGroupOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	params := ec2.CreateSecurityGroupInput{
		Description: description,
		GroupName:   groupName,
		VpcId:       vpcID,
	}

	describe := *&ec2.DescribeSecurityGroupsInput{
		GroupIds: []string{*groupId},
	}

	groups, err := c.Client.DescribeSecurityGroups(ctx, &describe)
	if err != nil {
		return nil, err
	}
	if len(groups.SecurityGroups) > 0 {
		slog.Info("Security group alread exists skipping creation")
		return nil, nil
	}

	output, err := c.Client.CreateSecurityGroup(ctx, &params)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (c *EC2Instances) SGEgress(s SGInput) (*ec2.AuthorizeSecurityGroupEgressOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	params := ec2.AuthorizeSecurityGroupEgressInput{
		GroupId:       s.SGId,
		IpPermissions: s.IpPermissions,
	}

	output, err := c.Client.AuthorizeSecurityGroupEgress(ctx, &params)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (c *EC2Instances) SGIngress(s SGInput) (*ec2.AuthorizeSecurityGroupIngressOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	params := ec2.AuthorizeSecurityGroupIngressInput{
		GroupId:       s.SGId,
		IpPermissions: s.IpPermissions,
	}

	output, err := c.Client.AuthorizeSecurityGroupIngress(ctx, &params)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (c *EC2Instances) DescribeSG(sgIds string) (*ec2.DescribeSecurityGroupsOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	params := ec2.DescribeSecurityGroupsInput{
		GroupIds: []string{sgIds},
	}

	output, err := c.Client.DescribeSecurityGroups(ctx, &params)
	if err != nil {
		return nil, err
	}
	return output, nil
}
