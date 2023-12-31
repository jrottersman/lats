package aws

import (
	"context"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// Ec2Client allows mocking of the ec2 client
type Ec2Client interface {
	CreateSecurityGroup(ctx context.Context, params *ec2.CreateSecurityGroupInput, optFns ...func(*ec2.Options)) (*ec2.CreateSecurityGroupOutput, error)
	AuthorizeSecurityGroupEgress(ctx context.Context, params *ec2.AuthorizeSecurityGroupEgressInput, optFns ...func(*ec2.Options)) (*ec2.AuthorizeSecurityGroupEgressOutput, error)
	DescribeSecurityGroups(ctx context.Context, params *ec2.DescribeSecurityGroupsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeSecurityGroupsOutput, error)
	AuthorizeSecurityGroupIngress(ctx context.Context, params *ec2.AuthorizeSecurityGroupIngressInput, optFns ...func(*ec2.Options)) (*ec2.AuthorizeSecurityGroupIngressOutput, error)
}

// CreateSGInput input for the create SG function
type CreateSGInput struct {
	description *string
	groupName   *string
	vpcID       *string
	groupID     *string
}

// SGInput input for updating security group
type SGInput struct {
	SGId          *string
	IPPermissions []types.IpPermission
}

// PassedIPs allows us to update our sg we need to transform this to an SGInput
type PassedIPs struct {
	Port        int
	Permissions string
	Description string
}

// EC2Instances is the struct to hold our ec2 client
type EC2Instances struct {
	Client Ec2Client
}

// CreateSG creates a new security group
func (c *EC2Instances) CreateSG(i CreateSGInput) (*ec2.CreateSecurityGroupOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	params := ec2.CreateSecurityGroupInput{
		Description: i.description,
		GroupName:   i.groupName,
		VpcId:       i.vpcID,
	}

	if i.groupID != nil {
		describe := *&ec2.DescribeSecurityGroupsInput{
			GroupIds: []string{*i.groupID},
		}

		groups, err := c.Client.DescribeSecurityGroups(ctx, &describe)
		if err != nil {
			return nil, err
		}
		if len(groups.SecurityGroups) > 0 {
			slog.Info("Security group alread exists skipping creation")
			return nil, nil
		}
	}

	output, err := c.Client.CreateSecurityGroup(ctx, &params)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// SGEgress updates a security group with egress info
func (c *EC2Instances) SGEgress(s SGInput) (*ec2.AuthorizeSecurityGroupEgressOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	params := ec2.AuthorizeSecurityGroupEgressInput{
		GroupId:       s.SGId,
		IpPermissions: s.IPPermissions,
	}

	output, err := c.Client.AuthorizeSecurityGroupEgress(ctx, &params)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// SGIngress updates a security group with ingress ips
func (c *EC2Instances) SGIngress(s SGInput) (*ec2.AuthorizeSecurityGroupIngressOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	params := ec2.AuthorizeSecurityGroupIngressInput{
		GroupId:       s.SGId,
		IpPermissions: s.IPPermissions,
	}

	output, err := c.Client.AuthorizeSecurityGroupIngress(ctx, &params)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// DescribeSG describes a security group
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
