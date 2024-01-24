package state

import (
	"bytes"
	"encoding/gob"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// SecurityGroupOutput is the wrapper for our security groups not sure we need this but for now it's here
type SecurityGroupOutput struct {
	SecurityGroups []types.SecurityGroup
}

// EncodeSecurityGroups encodes a security group to bytes
func EncodeSecurityGroups(sg SecurityGroupOutput) bytes.Buffer {
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)

	err := enc.Encode(sg)
	if err != nil {
		slog.Error("Error encoding our database", "error", err)
	}
	return encoder
}

// DecodeSecurityGroups takes bytes and gives us a securitygroupoutput for resotration
func DecodeSecurityGroups(b bytes.Buffer) SecurityGroupOutput {
	var securityGroups SecurityGroupOutput
	dec := gob.NewDecoder(&b)
	err := dec.Decode(&securityGroups)
	if err != nil {
		slog.Error("Error decoding state for Security Groups", "error", err)
	}
	return securityGroups
}

type SGRuleStorage struct {
	GroupID    *string
	GroupName  *string
	FromPort   *int32
	ToPort     *int32
	IPProtocol *string
	IPRanges   []types.IpRange
}

// SecurityGroupNeeds is a function that takes a security group and get's the parts we need out more for thought then anything
func SecurityGroupNeeds(sg SecurityGroupOutput) []SGRuleStorage {
	var sgRules []SGRuleStorage
	for _, v := range sg.SecurityGroups {
		gid := v.GroupId
		gname := v.GroupName
		for _, z := range v.IpPermissions {
			// What is needed for ipv4
			fromPort := z.FromPort
			toPort := z.ToPort
			IpProtocol := z.IpProtocol
			IpRanges := z.IpRanges
			sgRules = append(sgRules, SGRuleStorage{GroupID: gid, GroupName: gname, FromPort: fromPort, ToPort: toPort, IPProtocol: IpProtocol, IPRanges: IpRanges})
		}
	}
	return sgRules
}

func EncodeSGRulesStorage(sg []SGRuleStorage) bytes.Buffer {
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)

	err := enc.Encode(sg)
	if err != nil {
		slog.Error("Error encoding our Security Group rules", "error", err)
	}
	return encoder
}
