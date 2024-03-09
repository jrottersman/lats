package cmd

import (
	"fmt"
	"log/slog"

	"github.com/jrottersman/lats/aws"
	"github.com/jrottersman/lats/stack"
	"github.com/jrottersman/lats/state"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	//Variables used for flags
	restoreSnapshotName string
	restoreDbName       string
	region              string
	dbSubnetGroupName   string
	vpcID               string
	subnets             []string
	addresses           []string
	ports               []int
	ruleTypes           []string
	protocols           []string
	restConfigFile      string

	//RestoreRDSSnapshotCmd restores an RDS snapshot
	RestoreRDSSnapshotCmd = &cobra.Command{
		Use:     "restoreRDSSnapshot",
		Aliases: []string{"RestoreSnapshot"},
		Short:   "Restores an RDS snapshot",
		Long:    "Restores an RDS snapshot",
		Run: func(cmd *cobra.Command, args []string) {
			_, sm := GetState()
			RestoreSnapshot(sm, restoreSnapshotName)
		},
	}
)

func init() {
	RestoreRDSSnapshotCmd.Flags().StringVarP(&restoreSnapshotName, "snapshot-name", "s", "", "name of the snapshot we want to restore: choose one of snapshotName or db name")
	RestoreRDSSnapshotCmd.Flags().StringVarP(&restoreDbName, "database-name", "d", "", "name of the database we want to restore the snapshot for")
	RestoreRDSSnapshotCmd.Flags().StringVarP(&region, "region", "r", "", "AWS region we are restoring in")
	RestoreRDSSnapshotCmd.Flags().StringVarP(&dbSubnetGroupName, "subnet-group", "g", "", "DB subnet group we are restoring the snapshot to")
	RestoreRDSSnapshotCmd.Flags().StringVarP(&vpcID, "vpc-id", "v", "", "VPC Id we are restoring the db to")
	RestoreRDSSnapshotCmd.Flags().StringArrayVar(&subnets, "subnets", []string{}, "Subnets that we want to create a subnet group in")
	RestoreRDSSnapshotCmd.Flags().StringArrayVar(&addresses, "addresses", []string{}, "Addresses that we want to update our security group with")
	RestoreRDSSnapshotCmd.Flags().IntSliceVar(&ports, "ports", []int{}, "Ports that we want to update our security group with")
	RestoreRDSSnapshotCmd.Flags().StringArrayVar(&ruleTypes, "rule-types", []string{}, "Rule types that we want to update our security group with")
	RestoreRDSSnapshotCmd.Flags().StringArrayVar(&protocols, "protocols", []string{}, "Protocols that we want to update our security group with")
	RestoreRDSSnapshotCmd.Flags().StringVarP(&restConfigFile, "config-file", "f", "", "Config file for the snapshot that we want to parse")
}

// RestoreSnapshot is the function that restores a snapshot
func RestoreSnapshot(stateKV state.StateManager, restoreSnapshotName string) error {
	slog.Info("Starting restore snapshot procedure")
	dbi := aws.Init(region)
	ec2 := aws.InitEc2(region)

	slog.Info("Check if config file is passed in")
	var ingressRules []aws.PassedIPs
	var egressRules []aws.PassedIPs
	if restConfigFile != "" {
		slog.Info("Config file was passed in", "configFile", restConfigFile)
		viper.SetConfigFile(restConfigFile)
		viper.AddConfigPath(".")
		err := viper.ReadInConfig() // Find and read the config file
		if err != nil {             // Handle errors reading the config file
			slog.Error("Viper Error reading config", "error", err)
		}
		restoreSnapshotName = viper.Get("snapshot").(string)
		restoreDbName = viper.Get("database").(string)
		region = viper.Get("region").(string)
		dbSubnetGroupName = viper.Get("dbsubnetgroupname").(string)
		vpcID = viper.Get("vpcid").(string)
		subnets = viper.GetStringSlice("subnets")
		securityGroups := []sgRulesStruct{}
		err = viper.UnmarshalKey("securitygroups", &securityGroups)
		if err != nil {
			slog.Error("error unmarshalling security groups", "error", err)
		}

		for _, v := range securityGroups {
			slog.Info("Security Group", "sg", v)
			pi := aws.PassedIPs{
				Port:        v.port,
				Type:        v.ruleType,
				Protocol:    v.protocol,
				Permissions: v.source,
			}
			slog.Info("PassedIPs", "pi", pi)
			if v.ruleType == "ingress" {
				ingressRules = append(ingressRules, pi)
			} else if v.ruleType == "egress" {
				egressRules = append(egressRules, pi)
			}
		}
	}
	if ports != nil && addresses != nil && ruleTypes != nil && protocols != nil {
		if len(ports) == len(addresses) && len(addresses) == len(ruleTypes) && len(ruleTypes) == len(protocols) {
			for i := range ports {
				pi := aws.PassedIPs{
					Port:        ports[i],
					Type:        ruleTypes[i],
					Protocol:    protocols[i],
					Permissions: addresses[i],
				}
				if ruleTypes[i] == "ingress" {
					ingressRules = append(ingressRules, pi)
				} else if ruleTypes[i] == "egress" {
					egressRules = append(egressRules, pi)
				}
			}
		} else {
			slog.Error("error in parsing the security group rules")
		}
	}
	slog.Info("finding the stack")
	SnapshotStack, err := FindStack(stateKV, restoreSnapshotName)
	if err != nil {
		slog.Error("Error finding stack", "error", err)
	}
	slog.Info("Stack is", "stack", SnapshotStack)

	if SnapshotStack.RestorationObjectName == stack.Cluster && len(subnets) > 2 {
		slog.Error("subnet creation will fail for a cluster less then two azs", "subnets", subnets)
		return fmt.Errorf("error subnet creation will fail for a cluster less then two azs")
	}

	// Creating subnet group
	slog.Info("Db subnet group name", "dbSubnetGroupName", dbSubnetGroupName)
	if dbSubnetGroupName == "" {
		slog.Info("creating a subnet group")
		name := fmt.Sprintf("%s-subnets", restoreDbName)
		desc := fmt.Sprintf("%s-subnets created by lats for restoring database", restoreDbName)
		check, err := ec2.GetSubnets(subnets)
		if err != nil {
			slog.Error("error getting subnets", "error", err)
		}
		if len(check.Subnets) != len(subnets) { // TODO Check if all subnets exist this isn't quite good enough cause of the next token so we need to improve it
			slog.Error("error in subnets", "subnets", subnets)
			return fmt.Errorf("error in subnets")
		}
		slog.Info("Creating subnet group", "name", name, "description", desc, "subnets", subnets)
		sg, err := dbi.CreateDBSubnetGroup(name, desc, subnets)
		if err != nil {
			slog.Error("problem creating subnet group", "error", err)
		}
		dbSubnetGroupName = *sg.DBSubnetGroup.DBSubnetGroupName
	}

	slog.Info("starting restore", "type", SnapshotStack.RestorationObjectName)
	if SnapshotStack.RestorationObjectName == stack.Cluster {
		slog.Info("Restoring a cluster with inputs", "restoreDbName", "dbSubnetGroupName", "vpcID", restoreDbName, dbSubnetGroupName, vpcID)
		c := aws.CreateClusterFromStackInput{
			S:             SnapshotStack,
			ClusterName:   &restoreDbName,
			DBSubnetGroup: &dbSubnetGroupName,
			VpcID:         &vpcID,
			Ingress:       ingressRules,
			Egress:        egressRules,
		}
		return dbi.CreateClusterFromStack(c)
	} else if SnapshotStack.RestorationObjectName == stack.LoneInstance {
		slog.Info("Restoring an Instance with inputs", "restoreDbName", "dbSubnetGroupName", "vpcID", restoreDbName, dbSubnetGroupName, vpcID)
		c := aws.CreateInstanceFromStackInput{
			Stack:         SnapshotStack,
			DBName:        &restoreDbName,
			DBSubnetGroup: &dbSubnetGroupName,
			VpcID:         &vpcID,
			Ingress:       ingressRules,
			Egress:        egressRules,
		}
		return dbi.CreateInstanceFromStack(c)
	}

	slog.Error("Invalid type of stack for restoring an object", "StackType", SnapshotStack.RestorationObjectName)
	return fmt.Errorf("error invalid type of stack to restore a snapshot")
}
