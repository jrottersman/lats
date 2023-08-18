package state

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/rds/types"
)

type RDSRestorationStore struct {
	Snapshot *types.DBSnapshot
	Instance *types.DBInstance
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
	return &RDSRestorationStore{
		Snapshot: snap,
		Instance: db,
	}, nil
}
