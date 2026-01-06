package state

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
)

func TestCreateDBClusterInput(t *testing.T) {
	c := types.DBCluster{
		DBClusterIdentifier: aws.String("foo"),
	}
	result := CreateDBClusterInput(&c)
	if *result.DBClusterIdentifier != *c.DBClusterIdentifier {
		t.Errorf("got %s expected %s", *result.DBClusterIdentifier, *c.DBClusterIdentifier)
	}
	if result.PreSignedUrl != nil {
		t.Errorf("expected PreSignedUrl to be nil, got %v", result.PreSignedUrl)
	}
}
