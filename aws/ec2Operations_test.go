package aws

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	mock "github.com/jrottersman/lats/mocks"
)

func TestEC2Instances_CreateSG(t *testing.T) {
	type fields struct {
		Client Ec2Client
	}
	type args struct {
		description *string
		groupName   *string
		vpcID       *string
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
			fields: fields{Client: mock.MockEC2Client{}},
			args: args{
				description: aws.String("foo"),
				groupName:   aws.String("bar"),
				vpcID:       aws.String("baz"),
			},
			want: &ec2.CreateSecurityGroupOutput{
				GroupId: aws.String("foobar"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &EC2Instances{
				Client: tt.fields.Client,
			}
			got, err := c.CreateSG(tt.args.description, tt.args.groupName, tt.args.vpcID, nil)
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

func TestEC2Instances_SGEgress(t *testing.T) {
	type fields struct {
		Client Ec2Client
	}
	type args struct {
		s SGInput
	}
	tr := true
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ec2.AuthorizeSecurityGroupEgressOutput
		wantErr bool
	}{
		{name: "test",
			fields: fields{Client: mock.MockEC2Client{}},
			args:   args{s: SGInput{}},
			want: &ec2.AuthorizeSecurityGroupEgressOutput{
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
			got, err := c.SGEgress(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("EC2Instances.SGEgress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EC2Instances.SGEgress() = %v, want %v", got, tt.want)
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
				Client: mock.MockEC2Client{},
			},
			args:    args{sgName: "foo"},
			want:    &ec2.DescribeSecurityGroupsOutput{},
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
		s SGInput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ec2.AuthorizeSecurityGroupIngressOutput
		wantErr bool
	}{
		{name: "test",
			fields: fields{Client: mock.MockEC2Client{}},
			args:   args{s: SGInput{}},
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
			got, err := c.SGIngress(tt.args.s)
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
