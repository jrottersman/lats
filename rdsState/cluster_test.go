package rdsState_test

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/jrottersman/lats/aws"
	"github.com/jrottersman/lats/rdsState"
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

	clusterId := "foo"
	snapshotId := "bar"
	filen := "/tmp/bar"
	c := mockRDSClient{}
	i := aws.DbInstances{c}

	resto := state.RDSRestorationStore{
		Cluster:         &types.DBCluster{DBClusterIdentifier: &clusterId},
		ClusterSnapshot: &types.DBClusterSnapshot{DBClusterSnapshotIdentifier: &snapshotId},
	}
	arg := args{
		r:      resto,
		name:   "foo",
		fn:     &filen,
		client: i,
		folder: "/tmp",
	}

	tests := []struct {
		name    string
		args    args
		want    *state.Stack
		wantErr bool
	}{
		{"test", arg, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rdsState.GenerateRDSClusterStack(tt.args.r, tt.args.name, tt.args.fn, tt.args.client, tt.args.folder)
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
