package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/jrottersman/lats/state"
)

type ParameterGroup struct {
	ParameterGroup types.DBParameterGroup
	Params         []types.Parameter
}

func GetParameterGroups(r state.RDSRestorationStore, i DbInstances) ([]ParameterGroup, error) {
	pgs := r.GetParameterGroups()
	groups := []ParameterGroup{}
	for _, pg := range pgs {
		group, err := i.GetParameterGroup(*pg.DBParameterGroupName)
		if err != nil {
			return nil, fmt.Errorf("error getting parameter group %s", err)
		}
		if err != nil {
			return nil, fmt.Errorf("error getting parameters %s for group %s", err, *pg.DBParameterGroupName)
		}
		params, err := i.GetParametersForGroup(*pg.DBParameterGroupName)
		fpg := ParameterGroup{
			ParameterGroup: *group,
			Params:         *params,
		}
		groups = append(groups, fpg)
	}
	return groups, nil
}
