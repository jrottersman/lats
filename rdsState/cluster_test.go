package rdsState

import (
	"reflect"
	"testing"

	"github.com/jrottersman/lats/aws"
	"github.com/jrottersman/lats/state"
)

func TestGenerateRDSClusterStack(t *testing.T) {
	type args struct {
		r      state.RDSRestorationStore
		name   string
		fn     *string
		client aws.DbInstances
		folder string
	}
	tests := []struct {
		name    string
		args    args
		want    *state.Stack
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateRDSClusterStack(tt.args.r, tt.args.name, tt.args.fn, tt.args.client, tt.args.folder)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateRDSClusterStack() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateRDSClusterStack() = %v, want %v", got, tt.want)
			}
		})
	}
}
