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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SecurityGroupNeeds(tt.args.sg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SecurityGroupNeeds() = %v, want %v", got, tt.want)
			}
		})
	}
}
