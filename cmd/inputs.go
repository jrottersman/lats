package cmd

import (
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/jrottersman/lats/aws"
	"github.com/jrottersman/lats/state"
)

//CreateClusterSnapshotInput input for create snapshot for cluster
type CreateClusterSnapshotInput struct {
	dbi     aws.DbInstances
	ec2     aws.EC2Instances
	sm      state.StateManager
	cluster *types.DBCluster
	sfn     string
}

//CreateInstanceSnapshotInput input for create snapshot for instance
type CreateInstanceSnapshotInput struct {
	dbi aws.DbInstances
	ec2 aws.EC2Instances
	sm  state.StateManager
	sfn string
}

type SecurityGroupOutput struct {
	SecurityGroups []ec2types.SecurityGroup
}
