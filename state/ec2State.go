package state

import (
	"bytes"
	"encoding/gob"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

//SecurityGroupOutput is the wrapper for our security groups not sure we need this but for now it's here
type SecurityGroupOutput struct {
	SecurityGroups []types.SecurityGroup
}

//EncodeSecurityGroups encodes a security group to bytes
func EncodeSecurityGroups(sg SecurityGroupOutput) bytes.Buffer {
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)

	err := enc.Encode(sg)
	if err != nil {
		slog.Error("Error encoding our database", "error", err)
	}
	return encoder
}

//DecodeSecurityGroups takes bytes and gives us a securitygroupoutput for resotration
func DecodeSecurityGroups(b bytes.Buffer) SecurityGroupOutput {
	var securityGroups SecurityGroupOutput
	dec := gob.NewDecoder(&b)
	err := dec.Decode(&securityGroups)
	if err != nil {
		slog.Error("Error decoding state for Security Groups", "error", err)
	}
	return securityGroups
}
