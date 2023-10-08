package rdsstate

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"

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

	b := encodeParameterGroups(i.ParameterGroups)
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

func encodeParameterGroups(pgs []pgstate.ParameterGroup) bytes.Buffer {
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)
	err := enc.Encode(&pgs)
	if err != nil {
		log.Fatalf("Error encoding our parameters %s", err)
	}
	return encoder
}

//DecodeParameterGroups Decodes the parameter Group
func DecodeParameterGroups(b bytes.Buffer) []pgstate.ParameterGroup {
	var pg []pgstate.ParameterGroup
	dec := gob.NewDecoder(&b)
	err := dec.Decode(&pg)
	if err != nil {
		log.Fatalf("Error decoding parameters %s", err)
	}
	return pg
}
