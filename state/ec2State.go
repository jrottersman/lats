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
	GroupID       *string
	GroupName     *string
	FromPort      *int32
	ToPort        *int32
	IPProtocol    *string
	PrefixIdsList []string
	IPRanges      []types.IpRange
}

// SecurityGroupNeeds is a function that takes a security group and get's the parts we need out more for thought then anything
func SecurityGroupNeeds(sg SecurityGroupOutput) []SGRuleStorage {
	var sgRules []SGRuleStorage
	for _, v := range sg.SecurityGroups {
		gid := v.GroupId
		gname := v.GroupName
		for _, z := range v.IpPermissions {
			// What is needed for ipv4 and SG rules
			fromPort := z.FromPort
			toPort := z.ToPort
			IpProtocol := z.IpProtocol
			prefixes := []string{}
			for _, prefix := range z.PrefixListIds {
				prefixes = append(prefixes, *prefix.PrefixListId)
			}
			IpRanges := z.IpRanges
			sgRules = append(sgRules, SGRuleStorage{GroupID: gid, GroupName: gname, FromPort: fromPort, ToPort: toPort, IPProtocol: IpProtocol, IPRanges: IpRanges, PrefixIdsList: prefixes})
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

func DecodeSGRulesStorage(b bytes.Buffer) []SGRuleStorage {
	var sg []SGRuleStorage
	dec := gob.NewDecoder(&b)
	err := dec.Decode(&sg)
	if err != nil {
		slog.Error("Error decoding state for Security Group rules", "error", err)
	}
	return sg
}

func EncodeVpc(vpc types.Vpc) bytes.Buffer {
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)

	err := enc.Encode(vpc)
	if err != nil {
		slog.Error("Error encoding our VPC", "error", err)
	}
	return encoder
}

func DecodeVpc(b bytes.Buffer) types.Vpc {
	var vpc types.Vpc
	dec := gob.NewDecoder(&b)
	err := dec.Decode(&vpc)
	if err != nil {
		slog.Error("Error decoding state for VPC", "error", err)
	}
	return vpc
}

func EncodeSubnets(subnets []types.Subnet) bytes.Buffer {
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)

	err := enc.Encode(subnets)
	if err != nil {
		slog.Error("Error encoding our Subnets", "error", err)
	}
	return encoder
}

func DecodeSubnets(b bytes.Buffer) []types.Subnet {
	var subnets []types.Subnet
	dec := gob.NewDecoder(&b)
	err := dec.Decode(&subnets)
	if err != nil {
		slog.Error("Error decoding state for Subnets", "error", err)
	}
	return subnets
}

func EncodeInternetGateways(igws []types.InternetGateway) bytes.Buffer {
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)

	err := enc.Encode(igws)
	if err != nil {
		slog.Error("Error encoding our Internet Gateways", "error", err)
	}
	return encoder
}

func DecodeInternetGateways(b bytes.Buffer) []types.InternetGateway {
	var igws []types.InternetGateway
	dec := gob.NewDecoder(&b)
	err := dec.Decode(&igws)
	if err != nil {
		slog.Error("Error decoding state for Internet Gateways", "error", err)
	}
	return igws
}

func EncodeRouteTables(rts []types.RouteTable) bytes.Buffer {
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)

	err := enc.Encode(rts)
	if err != nil {
		slog.Error("Error encoding our Route Tables", "error", err)
	}
	return encoder
}

func SgRuleStorageToIpPermission(sg SGRuleStorage) types.IpPermission {
	var ip types.IpPermission
	ip.FromPort = sg.FromPort
	ip.ToPort = sg.ToPort
	ip.IpProtocol = sg.IPProtocol
	ip.IpRanges = sg.IPRanges
	ip.PrefixListIds = prefixlistGenerator(sg.PrefixIdsList)
	return ip
}

func prefixlistGenerator(pl []string) []types.PrefixListId {
	var plid []types.PrefixListId
	for _, v := range pl {
		plid = append(plid, types.PrefixListId{PrefixListId: &v})
	}
	return plid
}

func SgRuleStoragesToIpPermissions(s []SGRuleStorage) []types.IpPermission {
	var ips []types.IpPermission
	for _, v := range s {
		ips = append(ips, SgRuleStorageToIpPermission(v))
	}
	return ips
}
