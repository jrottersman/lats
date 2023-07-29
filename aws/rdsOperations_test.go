package aws

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
)

func TestDbInstances_GetInstance(t *testing.T) {
	type fields struct {
		RdsClient *rds.Client
	}
	type args struct {
		instanceName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.DBInstance
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instances := &DbInstances{
				RdsClient: tt.fields.RdsClient,
			}
			got, err := instances.GetInstance(tt.args.instanceName)
			if (err != nil) != tt.wantErr {
				t.Errorf("DbInstances.GetInstance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DbInstances.GetInstance() = %v, want %v", got, tt.want)
			}
		})
	}
}
