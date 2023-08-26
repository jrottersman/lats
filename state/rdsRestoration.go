package state

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/jrottersman/lats/helpers"
)

type RDSRestorationStore struct {
	Snapshot        *types.DBSnapshot
	Instance        *types.DBInstance
	Cluster         *types.DBCluster
	ClusterSnapshot *types.DBClusterSnapshot
}

func (r RDSRestorationStore) GetInstanceIdentifier() *string {
	if r.Instance == nil {
		return nil
	}
	if r.Instance.DBInstanceIdentifier == nil {
		return nil
	}
	return r.Instance.DBInstanceIdentifier
}

func (r RDSRestorationStore) GetInstanceClass() *string {
	if r.Instance == nil {
		return nil
	}
	if r.Instance.DBInstanceClass == nil {
		return nil
	}
	return r.Instance.DBInstanceClass
}

func (r RDSRestorationStore) GetAllocatedStorage() *int32 {
	if r.Snapshot == nil {
		return nil
	}
	if r.Snapshot.AllocatedStorage == 0 {
		return nil
	}
	return &r.Snapshot.AllocatedStorage
}

func (r RDSRestorationStore) GetSnapshotIdentifier() *string {
	if r.Snapshot == nil {
		return nil
	}
	if r.Snapshot.DBSnapshotIdentifier == nil {
		return nil
	}
	return r.Snapshot.DBSnapshotIdentifier
}

func (r RDSRestorationStore) GetClusterSnapshotIdentifier() *string {
	if r.ClusterSnapshot == nil {
		return nil
	}
	if r.ClusterSnapshot.DBClusterSnapshotIdentifier == nil {
		return nil
	}
	return r.ClusterSnapshot.DBClusterSnapshotIdentifier
}

func (r RDSRestorationStore) GetDBClusterIdentifier() *string {
	if r.ClusterSnapshot == nil {
		return nil
	}
	if r.ClusterSnapshot.DBClusterIdentifier == nil {
		return nil
	}
	return r.ClusterSnapshot.DBClusterIdentifier
}

func (r RDSRestorationStore) GetClusterEngine() *string {
	if r.ClusterSnapshot == nil {
		return nil
	}
	if r.ClusterSnapshot.Engine == nil {
		return nil
	}
	return r.ClusterSnapshot.Engine
}

func (r RDSRestorationStore) GetKmsKey() *string {
	if r.ClusterSnapshot == nil {
		return nil
	}
	if r.ClusterSnapshot.KmsKeyId == nil {
		return nil
	}
	return r.ClusterSnapshot.KmsKeyId
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

	cID := helpers.GetClusterId(db)
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
