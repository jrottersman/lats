package pgstate

import "github.com/aws/aws-sdk-go-v2/service/rds/types"

type ParameterGroup struct {
	ParameterGroup        types.DBParameterGroup
	ClusterParameterGroup types.DBClusterParameterGroup
	Params                []types.Parameter
}
