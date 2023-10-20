package aws

import (
	"fmt"

	"github.com/jrottersman/lats/pgstate"
	"github.com/jrottersman/lats/state"
)

//GetParameterGroups get a list of paramter groups
func GetParameterGroups(r state.RDSRestorationStore, i DbInstances) ([]pgstate.ParameterGroup, error) {
	pgs := r.GetParameterGroups()
	groups := []pgstate.ParameterGroup{}
	for _, pg := range pgs {
		group, err := i.GetParameterGroup(*pg.DBParameterGroupName)
		if err != nil {
			return nil, fmt.Errorf("error getting parameter group %s", err)
		}

		params, err := i.GetParametersForGroup(*pg.DBParameterGroupName)
		if err != nil {
			return nil, fmt.Errorf("error getting parameters %s for group %s", err, *pg.DBParameterGroupName)
		}
		fpg := pgstate.ParameterGroup{
			ParameterGroup: *group,
			Params:         *params,
		}
		groups = append(groups, fpg)
	}
	return groups, nil
}

func GetClusterParameterGroup(r state.RDSRestorationStore, i DbInstances) ([]pgstate.ParameterGroup, error) {
	pg := r.GetClusterParameterGroups()
	groups := []pgstate.ParameterGroup{}
	group, err := i.GetClusterParameterGroup(*pg)
	if err != nil {
		return nil, fmt.Errorf("error getting cluster parameter group %s", err)
	}
	params, err := i.GetParametersForGroup(*pg)
	if err != nil {
		return nil, fmt.Errorf("error getting parameters %s for group %s", err, *pg)
	}
	fpg := pgstate.ParameterGroup{
		ClusterParameterGroup: *group,
		Params:                *params,
	}
	groups = append(groups, fpg)
	return groups, nil
}
