package rdsState

import (
	"github.com/jrottersman/lats/aws"
	"github.com/jrottersman/lats/helpers"
	"github.com/jrottersman/lats/state"
)

type ClusterStackInput struct {
	R         state.RDSRestorationStore
	StackName string
	Filename  string
	Client    aws.DbInstances
	Folder    string
}

//GenerateRDSClusterStack creates a stack to restore a cluster and it's instances.
func GenerateRDSClusterStack(c ClusterStackInput) (*state.Stack, error) {
	if c.Filename == "" {
		c.Filename = *helpers.RandomStateFileName()
	}
	objMap := make(map[int][]state.Object)
	ClusterInput := state.GenerateRestoreDBClusterFromSnapshotInput(c.R)

	// This is the cluster
	bc := state.EncodeRestoreDBClusterFromSnapshotInput(ClusterInput)
	_, err := state.WriteOutput(c.Filename, bc)
	if err != nil {
		return nil, err
	}
	clusterObj := state.NewObject(c.Filename, 1, state.Cluster)
	var clusterObjects []state.Object
	clusterObjects = append(clusterObjects, clusterObj)
	objMap[1] = clusterObjects

	instanceObjects, err := ClusterInstancesToObjects(c.R.Cluster, c.Client, c.Folder, 2)
	if err != nil {
		return nil, err
	}
	objMap[2] = instanceObjects

	return &state.Stack{
		Name:                  c.StackName,
		RestorationObjectName: state.Cluster,
		Objects:               objMap,
	}, nil
}
