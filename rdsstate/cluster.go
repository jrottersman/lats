package rdsstate

import (
	"fmt"

	"github.com/jrottersman/lats/aws"
	"github.com/jrottersman/lats/helpers"
	"github.com/jrottersman/lats/pgstate"
	"github.com/jrottersman/lats/stack"
	"github.com/jrottersman/lats/state"

	"github.com/aws/aws-sdk-go-v2/service/rds/types"
)

//ClusterStackInput is the input for a ClusterStack
type ClusterStackInput struct {
	R                     state.RDSRestorationStore
	StackName             string
	Filename              string
	Client                aws.DbInstances
	Folder                string
	ParameterFileName     string
	ParameterGroups       []pgstate.ParameterGroup
	OptionGroupFileName   string
	OptionGroup           *types.OptionGroup
	SecurityGroups        *state.SecurityGroupOutput
	SecurityGroupFileName string
}

//GenerateRDSClusterStack creates a stack to restore a cluster and it's instances.
func GenerateRDSClusterStack(c ClusterStackInput) (*stack.Stack, error) {
	if c.Filename == "" {
		c.Filename = fmt.Sprintf(".state/%s", *helpers.RandomStateFileName())
	}
	if c.ParameterFileName == "" {
		c.ParameterFileName = fmt.Sprintf(".state/%s", *helpers.RandomStateFileName())
	}
	if c.OptionGroupFileName == "" {
		c.OptionGroupFileName = fmt.Sprintf(".state/%s", *helpers.RandomStateFileName())
	}
	if c.SecurityGroupFileName == "" {
		c.SecurityGroupFileName = fmt.Sprintf(".state/%s", *helpers.RandomStateFileName())
	}

	objMap := make(map[int][]stack.Object)
	bp := pgstate.EncodeParameterGroups(c.ParameterGroups)
	_, err := state.WriteOutput(c.ParameterFileName, bp)
	if err != nil {
		return nil, fmt.Errorf("error writing parameters %s", err)
	}
	parameterObj := stack.NewObject(c.ParameterFileName, 1, stack.DBClusterParameterGroup)
	var paramObjects []stack.Object
	paramObjects = append(paramObjects, parameterObj)

	if c.OptionGroup != nil {
		b := state.EncodeOptionGroup(c.OptionGroup)
		_, err := state.WriteOutput(c.OptionGroupFileName, b)
		if err != nil {
			return nil, fmt.Errorf("Error writing option Group %s", err)
		}
		optionObj := stack.NewObject(c.OptionGroupFileName, 1, stack.OptionGroup)
		paramObjects = append(paramObjects, optionObj)
	}

	objMap[1] = paramObjects

	ClusterInput := state.GenerateRestoreDBClusterFromSnapshotInput(c.R)

	// This is the cluster
	bc := state.EncodeRestoreDBClusterFromSnapshotInput(ClusterInput)
	_, err = helpers.WriteOutput(c.Filename, bc)
	if err != nil {
		return nil, err
	}
	clusterObj := stack.NewObject(c.Filename, 2, stack.Cluster)
	var clusterObjects []stack.Object
	clusterObjects = append(clusterObjects, clusterObj)
	objMap[2] = clusterObjects

	instanceObjects, err := ClusterInstancesToObjects(c.R.Cluster, c.Client, c.Folder, 3)
	if err != nil {
		return nil, err
	}
	objMap[3] = instanceObjects

	return &stack.Stack{
		Name:                  c.StackName,
		RestorationObjectName: stack.Cluster,
		Objects:               objMap,
	}, nil
}
