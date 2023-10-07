package rdsState

import (
	"fmt"

	"github.com/jrottersman/lats/aws"
	"github.com/jrottersman/lats/helpers"
	"github.com/jrottersman/lats/state"
)

type ClusterStackInput struct {
	R                 state.RDSRestorationStore
	StackName         string
	Filename          string
	Client            aws.DbInstances
	Folder            string
	ParameterFileName string
	ParameterGroups   []aws.ParameterGroup
}

//GenerateRDSClusterStack creates a stack to restore a cluster and it's instances.
func GenerateRDSClusterStack(c ClusterStackInput) (*state.Stack, error) {
	if c.Filename == "" {
		c.Filename = *helpers.RandomStateFileName()
	}
	if c.ParameterFileName == "" {
		c.ParameterFileName = *helpers.RandomStateFileName()
	}
	objMap := make(map[int][]state.Object)
	bp := encodeParameterGroups(c.ParameterGroups)
	_, err := state.WriteOutput(c.ParameterFileName, bp)
	if err != nil {
		return nil, fmt.Errorf("error writing parameters %s", err)
	}
	parameterObj := state.NewObject(c.ParameterFileName, 1, state.DBClusterParameterGroup)
	var paramObjects []state.Object
	paramObjects = append(paramObjects, parameterObj)
	objMap[1] = paramObjects

	ClusterInput := state.GenerateRestoreDBClusterFromSnapshotInput(c.R)

	// This is the cluster
	bc := state.EncodeRestoreDBClusterFromSnapshotInput(ClusterInput)
	_, err = state.WriteOutput(c.Filename, bc)
	if err != nil {
		return nil, err
	}
	clusterObj := state.NewObject(c.Filename, 1, state.Cluster)
	var clusterObjects []state.Object
	clusterObjects = append(clusterObjects, clusterObj)
	objMap[2] = clusterObjects

	instanceObjects, err := ClusterInstancesToObjects(c.R.Cluster, c.Client, c.Folder, 2)
	if err != nil {
		return nil, err
	}
	objMap[3] = instanceObjects

	return &state.Stack{
		Name:                  c.StackName,
		RestorationObjectName: state.Cluster,
		Objects:               objMap,
	}, nil
}
