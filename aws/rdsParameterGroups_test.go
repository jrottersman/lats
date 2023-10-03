package aws

import (
	"reflect"
	"testing"

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
