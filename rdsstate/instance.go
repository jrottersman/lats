package rdsstate

import (
	"fmt"

	"github.com/jrottersman/lats/helpers"
	"github.com/jrottersman/lats/pgstate"
	"github.com/jrottersman/lats/stack"
	"github.com/jrottersman/lats/state"

	"github.com/aws/aws-sdk-go-v2/service/rds/types"
)

// InstanceStackInputs struct to generate stack for an instance
type InstanceStackInputs struct {
	R                           state.RDSRestorationStore
	StackName                   string
	InstanceFileName            string
	ParameterFileName           string
	OptionGroupFileName         string
	ParameterGroups             []pgstate.ParameterGroup
	OptionGroup                 *types.OptionGroup
	SecurityGroups              *state.SecurityGroupOutput
	SecurityGroupsFileName      string
	SecurityGroupsRulesFileName string
}

// GenerateRDSInstanceStack creates a stack for restoration for an RDS instance
func GenerateRDSInstanceStack(i InstanceStackInputs) (*stack.Stack, error) {
	if i.InstanceFileName == "" {
		i.InstanceFileName = fmt.Sprintf(".state/%s", *helpers.RandomStateFileName())
	}

	if i.ParameterFileName == "" {
		i.ParameterFileName = fmt.Sprintf(".state/%s", *helpers.RandomStateFileName())
	}

	if i.OptionGroupFileName == "" {
		i.OptionGroupFileName = fmt.Sprintf(".state/%s", *helpers.RandomStateFileName())
	}

	if i.SecurityGroupsFileName == "" {
		i.SecurityGroupsFileName = fmt.Sprintf(".state/%s", *helpers.RandomStateFileName())
	}

	if i.SecurityGroupsRulesFileName == "" {
		i.SecurityGroupsRulesFileName = fmt.Sprintf(".state/%s", *helpers.RandomStateFileName())
	}

	b := pgstate.EncodeParameterGroups(i.ParameterGroups)
	_, err := state.WriteOutput(i.ParameterFileName, b)
	if err != nil {
		return nil, fmt.Errorf("error writing parameter groups %s", err)
	}
	paramObj := stack.NewObject(i.ParameterFileName, 1, stack.DBParameterGroup)
	var paramObjects []stack.Object
	paramObjects = append(paramObjects, paramObj)

	if i.OptionGroup != nil {
		b := state.EncodeOptionGroup(i.OptionGroup)
		_, err := state.WriteOutput(i.OptionGroupFileName, b)
		if err != nil {
			return nil, fmt.Errorf("error writing option groups %s", err)
		}
		optionObj := stack.NewObject(i.OptionGroupFileName, 1, stack.OptionGroup)
		paramObjects = append(paramObjects, optionObj)
	}

	if i.SecurityGroups != nil {
		b := state.EncodeSecurityGroups(*i.SecurityGroups)
		_, err := state.WriteOutput(i.SecurityGroupsFileName, b)
		if err != nil {
			return nil, fmt.Errorf("Error saving security groups %s", err)
		}
		sgObj := stack.NewObject(i.SecurityGroupsFileName, 1, stack.SecurityGroup)
		paramObjects = append(paramObjects, sgObj)

		sgRules := state.SecurityGroupNeeds(*i.SecurityGroups)
		b2 := state.EncodeSGRulesStorage(sgRules)
		_, err = state.WriteOutput(i.SecurityGroupsRulesFileName, b2)
		if err != nil {
			return nil, fmt.Errorf("Error saving security group rules %s", err)
		}
		sgRulesObj := stack.NewObject(i.SecurityGroupsRulesFileName, 1, stack.SecurityGroupRules)
		paramObjects = append(paramObjects, sgRulesObj)
	}

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
