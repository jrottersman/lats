package rdsState

import (
	"github.com/jrottersman/lats/helpers"
	"github.com/jrottersman/lats/state"
)

//GenerateRDSInstaceStack creates a stack for restoration for an RDS instance
func GenerateRDSInstanceStack(r state.RDSRestorationStore, name string, fn *string) (*state.Stack, error) {
	if fn == nil {
		fn = helpers.RandomStateFileName()
	}

	DBInput := state.GenerateRestoreDBInstanceFromDBSnapshotInput(r)

	b := state.EncodeRestoreDBInstanceFromDBSnapshotInput(DBInput)
	_, err := state.WriteOutput(*fn, b)
	if err != nil {
		return nil, err
	}

	instanceObj := state.NewObject(*fn, 1, state.LoneInstance) // 1 is the order currently we just have the instance so this is 1 we will have to update it once we are handling parameter groups

	var instanceObjects []state.Object
	instanceObjects = append(instanceObjects, instanceObj)

	m := make(map[int][]state.Object)
	m[1] = instanceObjects

	return &state.Stack{
		Name:                  name,
		RestorationObjectName: state.LoneInstance,
		Objects:               m,
	}, nil
}
