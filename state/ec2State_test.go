package state

import (
	"encoding/gob"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func TestEncodeSecurityGroups(t *testing.T) {

	sg := types.SecurityGroup{
		Description: aws.String("foo"),
	}
	sg2 := SecurityGroupOutput{
		SecurityGroups: []types.SecurityGroup{sg},
	}
	r := EncodeSecurityGroups(sg2)
	var result SecurityGroupOutput
	dec := gob.NewDecoder(&r)
	err := dec.Decode(&result)
	if err != nil {
		t.Errorf("decode error: %s", err)
	}
	if len(result.SecurityGroups) != len(sg2.SecurityGroups) {
		t.Errorf("got %d expected %d", len(result.SecurityGroups), len(sg2.SecurityGroups))
	}
}

func TestDecodeSecurityGroups(t *testing.T) {

	sg := types.SecurityGroup{
		Description: aws.String("foo"),
	}
	sg2 := SecurityGroupOutput{
		SecurityGroups: []types.SecurityGroup{sg},
	}
	r := EncodeSecurityGroups(sg2)
	result := DecodeSecurityGroups(r)
	if len(result.SecurityGroups) != len(sg2.SecurityGroups) {
		t.Errorf("got %d expected %d", len(result.SecurityGroups), len(sg2.SecurityGroups))
	}
}

func TestSecurityGroupNeeds(t *testing.T) {
	type args struct {
		sg SecurityGroupOutput
	}
	tests := []struct {
		name string
		args args
		want []SGRuleStorage
	}{
		{name: "test", args: args{sg: SecurityGroupOutput{SecurityGroups: []types.SecurityGroup{{GroupId: aws.String("foobar"), IpPermissions: []types.IpPermission{{IpProtocol: aws.String("tcp")}}}}}}, want: []SGRuleStorage{{GroupID: aws.String("foobar"), GroupName: nil, FromPort: nil, ToPort: nil, IPProtocol: aws.String("tcp"), IPRanges: nil, PrefixIdsList: []string{}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SecurityGroupNeeds(tt.args.sg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SecurityGroupNeeds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodeSGRuleStorage(t *testing.T) {
	var sgs []SGRuleStorage
	sgr := SGRuleStorage{
		GroupID: aws.String("foobar"),
	}
	sgs = append(sgs, sgr)
	r := EncodeSGRulesStorage(sgs)
	var result []SGRuleStorage
	dec := gob.NewDecoder(&r)
	err := dec.Decode(&result)
	if err != nil {
		t.Errorf("decode error: %s", err)
	}
	if len(result) != len(sgs) {
		t.Errorf("got %d expected %d", len(result), len(sgs))
	}
}

func TestDecodeSGRuleStorage(t *testing.T) {
	var sgs []SGRuleStorage
	sgr := SGRuleStorage{
		GroupID: aws.String("foobar"),
	}
	sgs = append(sgs, sgr)
	r := EncodeSGRulesStorage(sgs)
	result := DecodeSGRulesStorage(r)
	if len(result) != len(sgs) {
		t.Errorf("got %d expected %d", len(result), len(sgs))
	}
}

func TestEncodeVpc(t *testing.T) {
	vpc := types.Vpc{
		CidrBlock: aws.String("foo"),
	}
	r := EncodeVpc(vpc)
	var result types.Vpc
	dec := gob.NewDecoder(&r)
	err := dec.Decode(&result)
	if err != nil {
		t.Errorf("decode error: %s", err)
	}
	if *result.CidrBlock != *vpc.CidrBlock {
		t.Errorf("got %s expected %s", *result.CidrBlock, *vpc.CidrBlock)
	}
}

func TestDecodeVpc(t *testing.T) {
	vpc := types.Vpc{
		CidrBlock: aws.String("foo"),
	}
	r := EncodeVpc(vpc)
	result := DecodeVpc(r)
	if *result.CidrBlock != *vpc.CidrBlock {
		t.Errorf("got %s expected %s", *result.CidrBlock, *vpc.CidrBlock)
	}
}

func TestEncodeSubnets(t *testing.T) {
	sn := types.Subnet{
		CidrBlock: aws.String("foo"),
	}
	subnets := []types.Subnet{sn}
	r := EncodeSubnets(subnets)
	var result []types.Subnet
	dec := gob.NewDecoder(&r)
	err := dec.Decode(&result)
	if err != nil {
		t.Errorf("decode error: %s", err)
	}
	if len(result) != len(subnets) {
		t.Errorf("got %d expected %d", len(result), len(subnets))
	}
}

func TestDecodeSubnets(t *testing.T) {
	sn := types.Subnet{
		CidrBlock: aws.String("foo"),
	}
	subnets := []types.Subnet{sn}
	r := EncodeSubnets(subnets)
	result := DecodeSubnets(r)
	if len(result) != len(subnets) {
		t.Errorf("got %d expected %d", len(result), len(subnets))
	}
}

func TestEncodeInternetGateways(t *testing.T) {
	ig := types.InternetGateway{
		InternetGatewayId: aws.String("foo"),
	}
	igs := []types.InternetGateway{ig}
	r := EncodeInternetGateways(igs)
	var result []types.InternetGateway
	dec := gob.NewDecoder(&r)
	err := dec.Decode(&result)
	if err != nil {
		t.Errorf("decode error: %s", err)
	}
	if len(result) != len(igs) {
		t.Errorf("got %d expected %d", len(result), len(igs))
	}
}

func TestDecodeInternetGateways(t *testing.T) {
	ig := types.InternetGateway{
		InternetGatewayId: aws.String("foo"),
	}
	igs := []types.InternetGateway{ig}
	r := EncodeInternetGateways(igs)
	result := DecodeInternetGateways(r)
	if len(result) != len(igs) {
		t.Errorf("got %d expected %d", len(result), len(igs))
	}
}

func TestSgRuleStorageToIpPermission(t *testing.T) {
	type args struct {
		sg SGRuleStorage
	}
	tests := []struct {
		name string
		args args
		want types.IpPermission
	}{
		{name: "test", args: args{sg: SGRuleStorage{GroupID: aws.String("foobar"), FromPort: aws.Int32(8000), ToPort: aws.Int32(8000), IPProtocol: aws.String("tcp")}}, want: types.IpPermission{FromPort: aws.Int32(8000), ToPort: aws.Int32(8000), IpProtocol: aws.String("tcp"), IpRanges: nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SgRuleStorageToIpPermission(tt.args.sg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SgRuleStorageToIpPermission() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSgRuleStoragesToIpPermissions(t *testing.T) {
	type args struct {
		s []SGRuleStorage
	}
	tests := []struct {
		name string
		args args
		want []types.IpPermission
	}{
		{name: "test", args: args{s: []SGRuleStorage{{GroupID: aws.String("foobar"), FromPort: aws.Int32(8000), ToPort: aws.Int32(8000), IPProtocol: aws.String("tcp")}}}, want: []types.IpPermission{{FromPort: aws.Int32(8000), ToPort: aws.Int32(8000), IpProtocol: aws.String("tcp"), IpRanges: nil}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SgRuleStoragesToIpPermissions(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SgRuleStoragesToIpPermissions() = %v, want %v", got, tt.want)
			}
		})
	}
}
