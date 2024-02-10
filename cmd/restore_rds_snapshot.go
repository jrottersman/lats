package cmd

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

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
	ingress             []string
	egress              []string
	addresses           []string
	ports               []int
	ruleTypes           []string
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
	RestoreRDSSnapshotCmd.Flags().StringVarP(&dbSubnetGroupName, "vpc-id", "v", "", "VPC Id we are restoring the db to")
	RestoreRDSSnapshotCmd.Flags().StringArrayVar(&subnets, "subnets", []string{}, "Subnets that we want to create a subnet group in")
	RestoreRDSSnapshotCmd.Flags().StringArrayVar(&ingress, "ingress", []string{}, "Ingress rules that we want to update our security group with")
	RestoreRDSSnapshotCmd.Flags().StringArrayVar(&egress, "egress", []string{}, "Egress rules that we want to update our security group with")
	RestoreRDSSnapshotCmd.Flags().StringVarP(&restConfigFile, "config-file", "f", "", "Config file for the snapshot that we want to parse")
}

// RestoreSnapshot is the function that restores a snapshot
func RestoreSnapshot(stateKV state.StateManager, restoreSnapshotName string) error {
	slog.Info("Starting restore snapshot procedure")
	dbi := aws.Init(region)

	slog.Info("Check if config file is passed in")
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
		subnets = viper.("subnets")
		securityGroups := []sgRulesStruct{}
		err = viper.UnmarshalKey("securitygroups", &securityGroups)
		if err != nil {
			slog.Error("Error unmarshalling security groups", "error", err)
		}
	}
	slog.Info("finding the stack")
	SnapshotStack, err := FindStack(stateKV, restoreSnapshotName)
	if err != nil {
		slog.Error("Error finding stack", "error", err)
	}
	slog.Info("Stack is", "stack", SnapshotStack)

	// Creating subnet group
	if dbSubnetGroupName == "" {
		slog.Info("creating a subnet group")
		name := fmt.Sprintf("%s-subnets", restoreDbName)
		desc := fmt.Sprintf("%s-subnets created by lats for restoring database", restoreDbName)
		sg, err := dbi.CreateDBSubnetGroup(name, desc, subnets)
		if err != nil {
			slog.Error("problem creating subnet group", "error", err)
		}
		dbSubnetGroupName = *sg.DBSubnetGroup.DBSubnetGroupName
	}

	//Security Groups
	var ingressRules []aws.PassedIPs
	if len(ingress) != 0 {
		ingressRules = sgRuleConvert(ingress)
	}
	var egressRules []aws.PassedIPs
	if len(egress) != 0 {
		egressRules = sgRuleConvert(egress)
	}
	slog.Info("starting restore", "type", SnapshotStack.RestorationObjectName)
	if SnapshotStack.RestorationObjectName == stack.Cluster {
		slog.Info("Restoring a cluster")
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
		slog.Info("Restoring an Instance")
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
	return fmt.Errorf("Error invalid type of stack to restore a snapshot")
}

func sgRuleConvert(rules []string) []aws.PassedIPs {
	l := []aws.PassedIPs{}
	for _, v := range rules {
		res := strings.Split(v, "-")
		if len(res) != 2 {
			slog.Error("length of our split should be 2", "is", len(res))
			return nil
		}
		perms := res[0]
		port, err := strconv.Atoi(res[1])
		if err != nil {
			slog.Error("couldn't convert string to int", "error", err)
			return nil
		}
		s := aws.PassedIPs{
			Permissions: perms,
			Port:        port,
		}
		l = append(l, s)

	}
	return l
}
