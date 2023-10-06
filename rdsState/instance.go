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

//TODO refactor this to use a struct this is messy
//GenerateRDSInstaceStack creates a stack for restoration for an RDS instance
func GenerateRDSInstanceStack(r state.RDSRestorationStore, name string, fn *string, paramfn *string, pgs []aws.ParameterGroup) (*state.Stack, error) {
	if fn == nil {
		fn = helpers.RandomStateFileName()
	}

	if paramfn == nil {
		paramfn = helpers.RandomStateFileName()
	}

	b := encodeParameterGroups(pgs)
	_, err := state.WriteOutput(*paramfn, b)
	if err != nil {
		return nil, fmt.Errorf("error writing parameter groups %s", err)
	}
	paramObj := state.NewObject(*paramfn, 1, state.DBParameterGroup)
	var paramObjects []state.Object
	paramObjects = append(paramObjects, paramObj)

	DBInput := state.GenerateRestoreDBInstanceFromDBSnapshotInput(r)
	b = state.EncodeRestoreDBInstanceFromDBSnapshotInput(DBInput)
	_, err = state.WriteOutput(*fn, b)
	if err != nil {
		return nil, err
	}

	instanceObj := state.NewObject(*fn, 2, state.LoneInstance)

	var instanceObjects []state.Object
	instanceObjects = append(instanceObjects, instanceObj)

	m := make(map[int][]state.Object)
	m[1] = paramObjects
	m[2] = instanceObjects

	return &state.Stack{
		Name:                  name,
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
