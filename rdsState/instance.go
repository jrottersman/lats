package rdsState

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"

	"github.com/jrottersman/lats/aws"
	"github.com/jrottersman/lats/helpers"
	"github.com/jrottersman/lats/state"
)

type InstanceStackInputs struct {
	R                 state.RDSRestorationStore
	StackName         string
	InstanceFileName  string
	ParameterFileName string
	ParameterGroups   []aws.ParameterGroup
}

//GenerateRDSInstaceStack creates a stack for restoration for an RDS instance
func GenerateRDSInstanceStack(i InstanceStackInputs) (*state.Stack, error) {
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
	paramObj := state.NewObject(i.ParameterFileName, 1, state.DBParameterGroup)
	var paramObjects []state.Object
	paramObjects = append(paramObjects, paramObj)

	DBInput := state.GenerateRestoreDBInstanceFromDBSnapshotInput(i.R)
	b = state.EncodeRestoreDBInstanceFromDBSnapshotInput(DBInput)
	_, err = state.WriteOutput(i.InstanceFileName, b)
	if err != nil {
		return nil, err
	}

	instanceObj := state.NewObject(i.InstanceFileName, 2, state.LoneInstance)

	var instanceObjects []state.Object
	instanceObjects = append(instanceObjects, instanceObj)

	m := make(map[int][]state.Object)
	m[1] = paramObjects
	m[2] = instanceObjects

	return &state.Stack{
		Name:                  i.StackName,
		RestorationObjectName: state.LoneInstance,
		Objects:               m,
	}, nil
}

func encodeParameterGroups(pgs []aws.ParameterGroup) bytes.Buffer {
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)
	err := enc.Encode(&pgs)
	if err != nil {
		log.Fatalf("Error encoding our parameters %s", err)
	}
	return encoder
}

//DecodeParameterGroup Decodes the parameter Group
func DecodeParameterGroups(b bytes.Buffer) []aws.ParameterGroup {
	var pg []aws.ParameterGroup
	dec := gob.NewDecoder(&b)
	err := dec.Decode(&pg)
	if err != nil {
		log.Fatalf("Error decoding parameters %s", err)
	}
	return pg
}
