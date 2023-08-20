package state

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/jrottersman/lats/helpers"
)

type RDSRestorationStore struct {
	Snapshot *types.DBSnapshot
	Instance *types.DBInstance
	Cluster  *types.DBCluster
}

func RDSRestorationStoreBuilder(sm StateManager, snapshotName string) (*RDSRestorationStore, error) {
	snap, err := GetRDSSnapshotOutput(sm, snapshotName)
	if err != nil {
		fmt.Printf("got error getting snapshot %s", err)
		return nil, err
	}
	dbi := snap.DBInstanceIdentifier
	db, err := GetRDSDatabaseInstanceOutput(sm, *dbi)
	if err != nil {
		fmt.Printf("error getting database %s", err)
	}

	cID := helpers.GetClusterId(*db)
	if cID == nil {
		return &RDSRestorationStore{
			Snapshot: snap,
			Instance: db,
		}, nil
	}
	cluster, err := GetRDSDatabaseClusterOutput(sm, *cID)
	if err != nil {
		fmt.Printf("error getting cluster %s", err)
	}
	return &RDSRestorationStore{
		Snapshot: snap,
		Instance: db,
		Cluster:  cluster,
	}, nil

}
