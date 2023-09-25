package state

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/jrottersman/lats/helpers"
)

// RDSRestorationStore stores all of our things we could possibly need to restore a db or cluster
type RDSRestorationStore struct {
	Snapshot        *types.DBSnapshot
	Instance        *types.DBInstance
	Cluster         *types.DBCluster
	ClusterSnapshot *types.DBClusterSnapshot
}

// GetInstanceIdentifier get's the instance identifier from the restoration store
func (r RDSRestorationStore) GetInstanceIdentifier() *string {
	if r.Instance == nil {
		return nil
	}
	if r.Instance.DBInstanceIdentifier == nil {
		return nil
	}
	return r.Instance.DBInstanceIdentifier
}

//GetInstanceClass getter for instance class
func (r RDSRestorationStore) GetInstanceClass() *string {
	if r.Instance == nil {
		return nil
	}
	if r.Instance.DBInstanceClass == nil {
		return nil
	}
	return r.Instance.DBInstanceClass
}

//GetAllocatedStorage getter for allocated storage
func (r RDSRestorationStore) GetAllocatedStorage() *int32 {
	if r.Snapshot == nil {
		return nil
	}
	if r.Snapshot.AllocatedStorage == 0 {
		return nil
	}
	return &r.Snapshot.AllocatedStorage
}

//GetAutoMinorVersionUpgrade yet another getter
func (r RDSRestorationStore) GetAutoMinorVersionUpgrade() *bool {
	if r.Instance == nil {
		return nil
	}
	return &r.Instance.AutoMinorVersionUpgrade
}

//GetBackupTarget yet another getter
func (r RDSRestorationStore) GetBackupTarget() *string {
	if r.Instance == nil {
		return nil
	}
	return r.Instance.BackupTarget
}

//GetDeleteProtection yet another getter
func (r RDSRestorationStore) GetDeleteProtection() *bool {
	if r.Instance == nil {
		return nil
	}
	return &r.Instance.DeletionProtection
}

//GetSnapshotIdentifier yet another getter
func (r RDSRestorationStore) GetSnapshotIdentifier() *string {
	if r.Snapshot == nil {
		return nil
	}
	if r.Snapshot.DBSnapshotIdentifier == nil {
		return nil
	}
	return r.Snapshot.DBSnapshotIdentifier
}

//GetEnabledCloudwatchLogsExports yet another getter
func (r RDSRestorationStore) GetEnabledCloudwatchLogsExports() []string {
	if r.Instance == nil {
		return []string{}
	}
	if len(r.Instance.EnabledCloudwatchLogsExports) == 0 {
		return []string{}
	}
	return r.Instance.EnabledCloudwatchLogsExports
}

//GetClusterSnapshotIdentifier yet another getter
func (r RDSRestorationStore) GetClusterSnapshotIdentifier() *string {
	if r.ClusterSnapshot == nil {
		return nil
	}
	if r.ClusterSnapshot.DBClusterSnapshotIdentifier == nil {
		return nil
	}
	return r.ClusterSnapshot.DBClusterSnapshotIdentifier
}

//GetDBClusterIdentifier yet another getter
func (r RDSRestorationStore) GetDBClusterIdentifier() *string {
	if r.ClusterSnapshot == nil {
		return nil
	}
	if r.ClusterSnapshot.DBClusterIdentifier == nil {
		return nil
	}
	return r.ClusterSnapshot.DBClusterIdentifier
}

//GetClusterEngine yet another getter
func (r RDSRestorationStore) GetClusterEngine() *string {
	if r.ClusterSnapshot == nil {
		return nil
	}
	if r.ClusterSnapshot.Engine == nil {
		return nil
	}
	return r.ClusterSnapshot.Engine
}

//GetDBClusterInstanceClass yet another getter
func (r RDSRestorationStore) GetDBClusterInstanceClass() *string {
	if r.Cluster == nil {
		return nil
	}
	if r.Cluster.DBClusterInstanceClass == nil {
		return nil
	}
	return r.Cluster.DBClusterInstanceClass
}

//GetKmsKey yet another getter
func (r RDSRestorationStore) GetKmsKey() *string {
	if r.ClusterSnapshot == nil {
		return nil
	}
	if r.ClusterSnapshot.KmsKeyId == nil {
		return nil
	}
	return r.ClusterSnapshot.KmsKeyId
}

//GetClusterAZs yet another getter
func (r RDSRestorationStore) GetClusterAZs() *[]string {
	if r.ClusterSnapshot == nil {
		return nil
	}
	if r.ClusterSnapshot.AvailabilityZones == nil {
		return nil
	}
	return &r.ClusterSnapshot.AvailabilityZones
}

//GetDBClusterMembers yet another getter
func (r RDSRestorationStore) GetDBClusterMembers() *[]types.DBClusterMember {
	if r.Cluster == nil {
		return nil
	}
	if len(r.Cluster.DBClusterMembers) == 0 {
		return nil
	}
	return &r.Cluster.DBClusterMembers
}

//RDSRestorationStoreBuilder builds an RDSRestorationStore
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
