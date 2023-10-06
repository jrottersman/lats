package aws

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	mock "github.com/jrottersman/lats/mocks"
	"github.com/jrottersman/lats/state"
)

func TestGetParameterGroups(t *testing.T) {
	type args struct {
		r state.RDSRestorationStore
		i DbInstances
	}

	arg := args{
		r: state.RDSRestorationStore{},
		i: DbInstances{mock.MockRDSClient{}},
	}
	tests := []struct {
		name    string
		args    args
		want    []ParameterGroup
		wantErr bool
	}{
		{"test", arg, []ParameterGroup{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetParameterGroups(tt.args.r, tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetParameterGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetParameterGroups() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetClusterParameterGroup(t *testing.T) {
	type args struct {
		r state.RDSRestorationStore
		i DbInstances
	}
	cpg := types.DBClusterParameterGroup{DBClusterParameterGroupName: aws.String("foo")}
	wantParameterGroup := ParameterGroup{
		ClusterParameterGroup: cpg,
		Params:                []types.Parameter{},
	}
	arg := args{
		r: state.RDSRestorationStore{Cluster: &types.DBCluster{DBClusterParameterGroup: aws.String("foo")}},
		i: DbInstances{mock.MockRDSClient{}},
	}
	tests := []struct {
		name    string
		args    args
		want    []ParameterGroup
		wantErr bool
	}{
		{"test", arg, []ParameterGroup{wantParameterGroup}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetClusterParameterGroup(tt.args.r, tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetClusterParameterGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetClusterParameterGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}
