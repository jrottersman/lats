package state

import (
	"encoding/gob"
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
