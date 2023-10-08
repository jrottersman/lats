package rdsstate

import (
	"fmt"

	"github.com/jrottersman/lats/helpers"
	"github.com/jrottersman/lats/pgstate"
	"github.com/jrottersman/lats/stack"
	"github.com/jrottersman/lats/state"
)

//InstanceStackInputs struct to generate stack for an instance
type InstanceStackInputs struct {
	R                 state.RDSRestorationStore
	StackName         string
	InstanceFileName  string
	ParameterFileName string
	ParameterGroups   []pgstate.ParameterGroup
}

//GenerateRDSInstanceStack creates a stack for restoration for an RDS instance
func GenerateRDSInstanceStack(i InstanceStackInputs) (*stack.Stack, error) {
	if i.InstanceFileName == "" {
		i.InstanceFileName = *helpers.RandomStateFileName()
	}

	if i.ParameterFileName == "" {
		i.ParameterFileName = *helpers.RandomStateFileName()
	}

	b := pgstate.EncodeParameterGroups(i.ParameterGroups)
	_, err := state.WriteOutput(i.ParameterFileName, b)
	if err != nil {
		return nil, fmt.Errorf("error writing parameter groups %s", err)
	}
	paramObj := stack.NewObject(i.ParameterFileName, 1, stack.DBParameterGroup)
	var paramObjects []stack.Object
	paramObjects = append(paramObjects, paramObj)

	DBInput := state.GenerateRestoreDBInstanceFromDBSnapshotInput(i.R)
	b = state.EncodeRestoreDBInstanceFromDBSnapshotInput(DBInput)
	_, err = state.WriteOutput(i.InstanceFileName, b)
	if err != nil {
		return nil, err
	}

	instanceObj := stack.NewObject(i.InstanceFileName, 2, stack.LoneInstance)

	var instanceObjects []stack.Object
	instanceObjects = append(instanceObjects, instanceObj)

	m := make(map[int][]stack.Object)
	m[1] = paramObjects
	m[2] = instanceObjects

	return &stack.Stack{
		Name:                  i.StackName,
		RestorationObjectName: stack.LoneInstance,
		Objects:               m,
	}, nil
}
