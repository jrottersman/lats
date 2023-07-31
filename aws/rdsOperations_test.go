package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
)

type mockRDSClient struct{}

func (m mockRDSClient) DescribeDBInstances(ctx context.Context, input *rds.DescribeDBInstancesInput) (*rds.DescribeDBInstancesOutput, error) {
	f := "foo"
	r := &rds.DescribeDBInstancesOutput{
		DBInstances: []types.DBInstance{types.DBInstance{DBInstanceIdentifier: &f}},
	}
	return r, nil
}

func TestGetInstance(t *testing.T) {
	expected := "foo"
	c := mockRDSClient{}
	dbi := DbInstances{
		RdsClient: c,
	}
	resp, err := dbi.GetInstance("foo")
	if err != nil {
		t.Errorf("got the following error %v", err)
	}
	if *resp.DBInstanceIdentifier != expected {
		t.Errorf("got %s expected %s", *resp.DBInstanceIdentifier, expected)
	}
}