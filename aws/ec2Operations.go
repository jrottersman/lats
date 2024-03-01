package aws

import (
	"context"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/jrottersman/lats/state"
)

// Ec2Client allows mocking of the ec2 client
type Ec2Client interface {
	CreateSecurityGroup(ctx context.Context, params *ec2.CreateSecurityGroupInput, optFns ...func(*ec2.Options)) (*ec2.CreateSecurityGroupOutput, error)
	AuthorizeSecurityGroupEgress(ctx context.Context, params *ec2.AuthorizeSecurityGroupEgressInput, optFns ...func(*ec2.Options)) (*ec2.AuthorizeSecurityGroupEgressOutput, error)
	DescribeSecurityGroups(ctx context.Context, params *ec2.DescribeSecurityGroupsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeSecurityGroupsOutput, error)
	AuthorizeSecurityGroupIngress(ctx context.Context, params *ec2.AuthorizeSecurityGroupIngressInput, optFns ...func(*ec2.Options)) (*ec2.AuthorizeSecurityGroupIngressOutput, error)
	DescribeSubnets(ctx context.Context, params *ec2.DescribeSubnetsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeSubnetsOutput, error)
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
	Rules         []state.SGRuleStorage
}

// PassedIPs allows us to update our sg we need to transform this to an SGInput
type PassedIPs struct {
	Port        int
	Permissions string
	Protocol    string
	Description string
	Type        string
}

// CreateSgInput this only hanldes TCP and IPv4 right now this is a stub while I think of how to do it better
func (p PassedIPs) CreateSgInput(SGID *string) SGInput {
	port := int32(p.Port)
	cidr := types.IpRange{
		CidrIp: &p.Permissions,
	}
	cidrs := []types.IpRange{}
	cidrs = append(cidrs, cidr)
	ipPerms := []types.IpPermission{}
	ipPerm := types.IpPermission{
		FromPort:   &port,
		ToPort:     &port,
		IpRanges:   cidrs,
		IpProtocol: &p.Type,
	}
	ipPerms = append(ipPerms, ipPerm)
	return SGInput{
		SGId:          SGID,
		IPPermissions: ipPerms,
	}
}

// EC2Instances is the struct to hold our ec2 client
type EC2Instances struct {
	Client Ec2Client
}

// CreateSG creates a new security group
func (c *EC2Instances) CreateSG(i CreateSGInput) (*ec2.CreateSecurityGroupOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	if i.groupID != nil {
		describe := ec2.DescribeSecurityGroupsInput{
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

	params := ec2.CreateSecurityGroupInput{
		Description: i.description,
		GroupName:   i.groupName,
		VpcId:       i.vpcID,
	}

	output, err := c.Client.CreateSecurityGroup(ctx, &params)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// SGIngress updates a security group with ingress ips
func (c *EC2Instances) SGIngress(sgname string, s []PassedIPs) (*ec2.AuthorizeSecurityGroupIngressOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	IPPermissions := []types.IpPermission{}
	if len(s) > 0 {

		for _, v := range s {

			port := int32(v.Port)
			permissions := types.IpRange{
				CidrIp: &v.Permissions,
			}

			perm := types.IpPermission{
				FromPort:   &port,
				ToPort:     &port,
				IpProtocol: &v.Protocol,
				IpRanges:   []types.IpRange{permissions},
			}
			IPPermissions = append(IPPermissions, perm)
		}
	}

	params := ec2.AuthorizeSecurityGroupIngressInput{
		GroupId:       &sgname,
		IpPermissions: IPPermissions,
	}

	output, err := c.Client.AuthorizeSecurityGroupIngress(ctx, &params)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (c *EC2Instances) SGEgress(sgname string, s []PassedIPs) (*ec2.AuthorizeSecurityGroupEgressOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	IPPermissions := []types.IpPermission{}
	if len(s) > 0 {

		for _, v := range s {

			port := int32(v.Port)
			permissions := types.IpRange{
				CidrIp: &v.Permissions,
			}

			perm := types.IpPermission{
				FromPort:   &port,
				ToPort:     &port,
				IpProtocol: &v.Protocol,
				IpRanges:   []types.IpRange{permissions},
			}
			IPPermissions = append(IPPermissions, perm)
		}
	}

	params := ec2.AuthorizeSecurityGroupEgressInput{
		GroupId:       &sgname,
		IpPermissions: IPPermissions,
	}

	output, err := c.Client.AuthorizeSecurityGroupEgress(ctx, &params)
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

func (c *EC2Instances) GetSubnet(subnetID string) (*ec2.DescribeSubnetsOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	params := ec2.DescribeSubnetsInput{
		SubnetIds: []string{subnetID},
	}

	output, err := c.Client.DescribeSubnets(ctx, &params)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (c *EC2Instances) GetSGs(sgIds []string) (*ec2.DescribeSecurityGroupsOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	params := ec2.DescribeSecurityGroupsInput{
		GroupIds: sgIds,
	}

	output, err := c.Client.DescribeSecurityGroups(ctx, &params)
	if err != nil {
		return nil, err
	}
	return output, nil
}
