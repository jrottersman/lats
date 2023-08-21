package helpers

import "github.com/aws/aws-sdk-go-v2/service/rds/types"

func GetClusterId(t *types.DBInstance) *string {
	if t.DBClusterIdentifier == nil {
		return nil
	}
	return t.DBClusterIdentifier
}