package rdsState

import (
	"github.com/jrottersman/lats/aws"
	"github.com/jrottersman/lats/helpers"
	"github.com/jrottersman/lats/state"
)

func GenerateRDSClusterStack(r state.RDSRestorationStore, name string, fn *string, client aws.DbInstances) (*state.Stack, error) {
	if fn == nil {
		fn = helpers.RandomStateFileName()
	}

	ClusterInput := state.GenerateRestoreDBClusterFromSnapshotInput(r)

	// This is the cluster
	bc := state.EncodeRestoreDBClusterFromSnapshotInput(ClusterInput)
	_, err := state.WriteOutput(*fn, bc)
	if err != nil {
		return nil, err
	}
	clusterObj := state.NewObject(*fn, 1, state.Cluster)
	var firstObjects []state.Object
	firstObjects = append(firstObjects, clusterObj)

	// TODO figure out how to handle the instances
	// instanceObjects := ClusterInstancesToObjects(r.Cluster, client, )
	return nil, nil
}
