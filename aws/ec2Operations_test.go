package aws

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	mock "github.com/jrottersman/lats/mocks"
)

func TestEC2Instances_CreateSG(t *testing.T) {
	type fields struct {
		Client Ec2Client
	}
	type args struct {
		i CreateSGInput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ec2.CreateSecurityGroupOutput
		wantErr bool
	}{
		{
			name:   "pass",
			fields: fields{Client: mock.EC2Client{}},
			args: args{
				CreateSGInput{
					description: aws.String("foo"),
					groupName:   aws.String("bar"),
					vpcID:       aws.String("baz"),
					groupID:     nil,
				},
			},
			want: &ec2.CreateSecurityGroupOutput{
				GroupId: aws.String("foobar"),
			},
			wantErr: false,
		},
		{
			name:   "fail",
			fields: fields{Client: mock.EC2Client{}},
			args: args{
				CreateSGInput{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "describe",
			fields: fields{Client: mock.EC2Client{}},
			args: args{
				CreateSGInput{
					description: aws.String("foo"),
					groupID:     aws.String("foobar"),
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &EC2Instances{
				Client: tt.fields.Client,
			}
			got, err := c.CreateSG(tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("EC2Instances.CreateSG() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EC2Instances.CreateSG() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEC2Instances_DescribeSG(t *testing.T) {
	type fields struct {
		Client Ec2Client
	}
	type args struct {
		sgName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ec2.DescribeSecurityGroupsOutput
		wantErr bool
	}{
		{
			name: "pass",
			fields: fields{
				Client: mock.EC2Client{},
			},
			args: args{sgName: "foo"},
			want: &ec2.DescribeSecurityGroupsOutput{
				SecurityGroups: []types.SecurityGroup{{GroupId: aws.String("foobar")}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &EC2Instances{
				Client: tt.fields.Client,
			}
			got, err := c.DescribeSG(tt.args.sgName)
			if (err != nil) != tt.wantErr {
				t.Errorf("EC2Instances.DescribeSG() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EC2Instances.DescribeSG() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEC2Instances_SGIngress(t *testing.T) {
	tr := true
	type fields struct {
		Client Ec2Client
	}
	type args struct {
		n string
		s []PassedIPs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ec2.AuthorizeSecurityGroupIngressOutput
		wantErr bool
	}{
		{name: "test",
			fields: fields{Client: mock.EC2Client{}},
			args:   args{n: "foo", s: []PassedIPs{}},
			want: &ec2.AuthorizeSecurityGroupIngressOutput{
				Return: &tr,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &EC2Instances{
				Client: tt.fields.Client,
			}
			got, err := c.SGIngress(tt.args.n, tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("EC2Instances.SGIngress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EC2Instances.SGIngress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPassedIPs_CreateSgInput(t *testing.T) {
	type fields struct {
		Port        int
		Permissions string
		Description string
		Type        string
	}
	type args struct {
		SGID *string
	}
	sgid := "sg-1234"
	ipp := types.IpPermission{
		FromPort:   aws.Int32(80),
		ToPort:     aws.Int32(80),
		IpProtocol: aws.String("tcp"),
		IpRanges:   []types.IpRange{{CidrIp: aws.String("10.0.0.4/22")}},
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   SGInput
	}{
		{
			name: "pass",
			fields: fields{
				Port:        80,
				Permissions: "10.0.0.4/22",
				Description: "foo",
				Type:        "tcp",
			},
			args: args{SGID: &sgid},
			want: SGInput{
				SGId: &sgid,
				IPPermissions: []types.IpPermission{
					ipp,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PassedIPs{
				Port:        tt.fields.Port,
				Permissions: tt.fields.Permissions,
				Description: tt.fields.Description,
				Type:        tt.fields.Type,
			}
			if got := p.CreateSgInput(tt.args.SGID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PassedIPs.CreateSgInput() = %v, want %v", got.IPPermissions, tt.want.IPPermissions)
			}
		})
	}
}

func TestEC2Instances_SGEgress(t *testing.T) {
	type fields struct {
		Client Ec2Client
	}
	type args struct {
		n string
		s []PassedIPs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ec2.AuthorizeSecurityGroupEgressOutput
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "test",
			fields: fields{Client: mock.EC2Client{}},
			args:   args{n: "foo", s: []PassedIPs{}},
			want: &ec2.AuthorizeSecurityGroupEgressOutput{
				Return: aws.Bool(true),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &EC2Instances{
				Client: tt.fields.Client,
			}
			got, err := c.SGEgress(tt.args.n, tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("EC2Instances.SGEngress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EC2Instances.SGEngress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEC2Instances_GetSubnet(t *testing.T) {
	type fields struct {
		Client Ec2Client
	}
	type args struct {
		subnetID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ec2.DescribeSubnetsOutput
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &EC2Instances{
				Client: tt.fields.Client,
			}
			got, err := c.GetSubnet(tt.args.subnetID)
			if (err != nil) != tt.wantErr {
				t.Errorf("EC2Instances.GetSubnet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EC2Instances.GetSubnet() = %v, want %v", got, tt.want)
			}
		})
	}
}
