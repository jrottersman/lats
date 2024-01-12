package helpers

import "github.com/aws/aws-sdk-go-v2/service/rds/types"

// GetClusterID get's the clsuter id or nil as a getter this is the only one in helpers because it's needed for the initial load
func GetClusterID(t *types.DBInstance) *string {
	if t.DBClusterIdentifier == nil {
		return nil
	}
	return t.DBClusterIdentifier
}
