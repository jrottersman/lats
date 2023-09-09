package rds

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/jrottersman/lats/aws"
	"github.com/jrottersman/lats/state"
)

// ClusterInstancesToObjects makes a list of instances as objects for our stack
func ClusterInstancesToObjects(t *types.DBCluster, c aws.DbInstances) ([]state.Object, error) {
	// Cluster is empty
	if len(t.DBClusterMembers) == 0 {
		return nil, nil
	}
	for _, v := range t.DBClusterMembers {
		// GEtInstance here that will get us an error an instanc type we then need to generate an object from this instance
		inst, err := c.GetInstance(*v.DBInstanceIdentifier)
		if err != nil {
			fmt.Printf("error %s getting instance %s", err, *v.DBInstanceIdentifier)
		}
		state.CreateInstanceInput(inst, t.DBClusterIdentifier)
	}
	return nil, nil
}
