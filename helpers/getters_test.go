package helpers

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
)

func TestGetClusterId(t *testing.T) {
	ni := types.DBInstance{}
	nr := GetClusterId(&ni)
	if nr != nil {
		t.Errorf("this should be nil")
	}

	i := types.DBInstance{
		DBClusterIdentifier: aws.String("foobar"),
	}
	r := GetClusterId(&i)
	if *r != "foobar" {
		t.Errorf("got %s expected foobar", *r)
	}

}
