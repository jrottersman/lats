package mock

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// EC2Client is a client for ec2 that's a mock
type EC2Client struct {
}

// CreateSecurityGroup mock security group creaton
func (m EC2Client) CreateSecurityGroup(ctx context.Context, params *ec2.CreateSecurityGroupInput, optFns ...func(*ec2.Options)) (*ec2.CreateSecurityGroupOutput, error) {
	if params.Description == nil {
		return nil, fmt.Errorf("nil description error this is fake")
	}
	return &ec2.CreateSecurityGroupOutput{
		GroupId: aws.String("foobar"),
	}, nil
}

// AuthorizeSecurityGroupEgress mock authorize security group egress
func (m EC2Client) AuthorizeSecurityGroupEgress(ctx context.Context, params *ec2.AuthorizeSecurityGroupEgressInput, optFns ...func(*ec2.Options)) (*ec2.AuthorizeSecurityGroupEgressOutput, error) {
	boo := true
	return &ec2.AuthorizeSecurityGroupEgressOutput{Return: &boo}, nil
}

// DescribeSecurityGroups mock descirbe security groups
func (m EC2Client) DescribeSecurityGroups(ctx context.Context, params *ec2.DescribeSecurityGroupsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeSecurityGroupsOutput, error) {
	return &ec2.DescribeSecurityGroupsOutput{
		SecurityGroups: []types.SecurityGroup{{GroupId: aws.String("foobar")}},
	}, nil
}

// AuthorizeSecurityGroupIngress another mock
func (m EC2Client) AuthorizeSecurityGroupIngress(ctx context.Context, params *ec2.AuthorizeSecurityGroupIngressInput, optFns ...func(*ec2.Options)) (*ec2.AuthorizeSecurityGroupIngressOutput, error) {
	boo := true
	return &ec2.AuthorizeSecurityGroupIngressOutput{Return: &boo}, nil
}

func (m EC2Client) DescribeSubnets(ctx context.Context, params *ec2.DescribeSubnetsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeSubnetsOutput, error) {
	return &ec2.DescribeSubnetsOutput{}, nil
}
