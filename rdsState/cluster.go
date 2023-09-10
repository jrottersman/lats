package rdsState

import (
	"github.com/jrottersman/lats/aws"
	"github.com/jrottersman/lats/helpers"
	"github.com/jrottersman/lats/state"
)

func GenerateRDSClusterStack(r state.RDSRestorationStore, name string, fn *string, client aws.DbInstances, folder string) (*state.Stack, error) {
	if fn == nil {
		fn = helpers.RandomStateFileName()
	}
	objMap := make(map[int][]Objects)
	ClusterInput := state.GenerateRestoreDBClusterFromSnapshotInput(r)

	// This is the cluster
	bc := state.EncodeRestoreDBClusterFromSnapshotInput(ClusterInput)
	_, err := state.WriteOutput(*fn, bc)
	if err != nil {
		return nil, err
	}
	clusterObj := state.NewObject(*fn, 1, state.Cluster)
	var clusterObjects []state.Object
	clusterObjects = append(clustertObjects, clusterObj)
	objMap[1] = clusterObjects

	// TODO figure out how to handle the instances
	instanceObjects, err := ClusterInstancesToObjects(r.Cluster, client, folder, 2)
	if err != nil {
		return nil, err
	}
	objMap[2] = instanceObjects
	
	return state.Stack{
		Name: name,
		RestorationObjectName: state.Cluster,
		Objects: objMap

	}, nil
}
