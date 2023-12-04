package state

import (
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
)

func TestRDSRestorationStoreBuilder(t *testing.T) {
	filename := "/tmp/foo"
	var storage int32
	var progress int32
	storage = 1000
	progress = 100
	tu := true
	snap := types.DBSnapshot{
		AllocatedStorage:     &storage,
		Encrypted:            &tu,
		PercentProgress:      &progress,
		DBInstanceIdentifier: aws.String("foobar"),
	}

	defer os.Remove(filename)
	r := EncodeRDSSnapshotOutput(&snap)
	_, err := WriteOutput(filename, r)
	if err != nil {
		t.Errorf("error writing file %s", err)
	}

	var mu sync.Mutex
	var s []StateKV
	kv := StateKV{
		Object:       "foo",
		FileLocation: filename,
		ObjectType:   "RDSSnapshot",
	}
	s = append(s, kv)

	filename2 := "/tmp/foo2"
	dbi := types.DBInstance{
		AllocatedStorage:     &storage,
		DBInstanceIdentifier: aws.String("foobar"),
	}

	defer os.Remove(filename2)
	r2 := EncodeRDSDatabaseOutput(&dbi)
	_, err = WriteOutput(filename2, r2)
	if err != nil {
		t.Errorf("error writing file %s", err)
	}
	kv2 := StateKV{
		Object:       "foobar",
		FileLocation: filename2,
		ObjectType:   RdsInstanceType,
	}
	s = append(s, kv2)
	sm := StateManager{
		&mu,
		s,
	}

	resp, err := RDSRestorationStoreBuilder(sm, "foo")
	if err != nil {
		t.Errorf("error writing file %s", err)
	}

	if *resp.Snapshot.AllocatedStorage != 1000 {
		t.Errorf("oops WTF wanted 1000 got %d", *resp.Snapshot.AllocatedStorage)
	}
}

func TestRDSRestorationStore_GetAllocatedStorage(t *testing.T) {
	var valueExpected int32 = 1000
	type fields struct {
		Snapshot *types.DBSnapshot
		Instance *types.DBInstance
		Cluster  *types.DBCluster
	}
	tests := []struct {
		name   string
		fields fields
		want   *int32
	}{
		{name: "totalNil", fields: fields{nil, nil, nil}, want: nil},
		{name: "RegularNil", fields: fields{&types.DBSnapshot{}, nil, nil}, want: nil},
		{name: "Value", fields: fields{&types.DBSnapshot{AllocatedStorage: &valueExpected}, nil, nil}, want: &valueExpected},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RDSRestorationStore{
				Snapshot: tt.fields.Snapshot,
				Instance: tt.fields.Instance,
				Cluster:  tt.fields.Cluster,
			}
			got := r.GetAllocatedStorage()
			if tt.want == nil {
				if got != nil {
					t.Errorf("RDSRestorationStore.GetAllocatedStorage() = %v, want %v", got, tt.want)
				}
			} else {
				if *got != *tt.want {
					t.Errorf("RDSRestorationStore.GetAllocatedStorage() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestRDSRestorationStore_GetInstanceIdentifier(t *testing.T) {
	type fields struct {
		Snapshot *types.DBSnapshot
		Instance *types.DBInstance
		Cluster  *types.DBCluster
	}
	tests := []struct {
		name   string
		fields fields
		want   *string
	}{
		{name: "totalNil", fields: fields{nil, nil, nil}, want: nil},
		{name: "RegularNil", fields: fields{nil, &types.DBInstance{}, nil}, want: nil},
		{name: "GetData", fields: fields{nil, &types.DBInstance{DBInstanceIdentifier: aws.String("foo")}, nil}, want: aws.String("foo")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RDSRestorationStore{
				Snapshot: tt.fields.Snapshot,
				Instance: tt.fields.Instance,
				Cluster:  tt.fields.Cluster,
			}
			got := r.GetInstanceIdentifier()
			if tt.want == nil {
				if got != tt.want {
					t.Errorf("RDSRestorationStore.GetInstanceIdentifier() = %v, want %v", got, tt.want)
				}
			} else {
				if *got != *tt.want {
					t.Errorf("RDSRestorationStore.GetInstanceIdentifier() = %v, want %v", *got, *tt.want)
				}
			}
		})
	}
}

func TestRDSRestorationStore_GetInstanceClass(t *testing.T) {
	type fields struct {
		Snapshot *types.DBSnapshot
		Instance *types.DBInstance
		Cluster  *types.DBCluster
	}
	tests := []struct {
		name   string
		fields fields
		want   *string
	}{
		{name: "totalNil", fields: fields{nil, nil, nil}, want: nil},
		{name: "RegularNil", fields: fields{nil, &types.DBInstance{}, nil}, want: nil},
		{name: "GetData", fields: fields{nil, &types.DBInstance{DBInstanceClass: aws.String("t3.micro")}, nil}, want: aws.String("t3.micro")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RDSRestorationStore{
				Snapshot: tt.fields.Snapshot,
				Instance: tt.fields.Instance,
				Cluster:  tt.fields.Cluster,
			}
			got := r.GetInstanceClass()
			if tt.want == nil {
				if got != tt.want {
					t.Errorf("RDSRestorationStore.GetInstanceClass() = %v, want %v", got, tt.want)
				}
			} else {
				if *got != *tt.want {
					t.Errorf("RDSRestorationStore.GetInstanceClass() = %v, want %v", *got, *tt.want)
				}
			}
		})
	}
}

func TestRDSRestorationStore_GetSnapshotIdentifier(t *testing.T) {
	type fields struct {
		Snapshot *types.DBSnapshot
		Instance *types.DBInstance
		Cluster  *types.DBCluster
	}
	tests := []struct {
		name   string
		fields fields
		want   *string
	}{
		{name: "totalNil", fields: fields{nil, nil, nil}, want: nil},
		{name: "RegularNil", fields: fields{&types.DBSnapshot{}, nil, nil}, want: nil},
		{name: "Value", fields: fields{&types.DBSnapshot{DBSnapshotIdentifier: aws.String("foo")}, nil, nil}, want: aws.String("foo")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RDSRestorationStore{
				Snapshot: tt.fields.Snapshot,
				Instance: tt.fields.Instance,
				Cluster:  tt.fields.Cluster,
			}
			got := r.GetSnapshotIdentifier()
			if tt.want == nil {
				if got != tt.want {
					t.Errorf("RDSRestorationStore.GetSnapshotIdentifier() = %v, want %v", got, tt.want)
				}
			} else {
				if *got != *tt.want {
					t.Errorf("RDSRestorationStore.GetSnapshotIdentifier() = %v, want %v", *got, *tt.want)
				}
			}
		})
	}
}

func TestRDSRestorationStore_GetClusterSnapshotIdentifier(t *testing.T) {
	type fields struct {
		Snapshot        *types.DBSnapshot
		Instance        *types.DBInstance
		Cluster         *types.DBCluster
		ClusterSnapshot *types.DBClusterSnapshot
	}
	tests := []struct {
		name   string
		fields fields
		want   *string
	}{
		{name: "totalNil", fields: fields{nil, nil, nil, nil}, want: nil},
		{name: "RegularNil", fields: fields{nil, nil, nil, &types.DBClusterSnapshot{}}, want: nil},
		{name: "Value", fields: fields{nil, nil, nil, &types.DBClusterSnapshot{DBClusterSnapshotIdentifier: aws.String("foo")}}, want: aws.String("foo")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RDSRestorationStore{
				Snapshot:        tt.fields.Snapshot,
				Instance:        tt.fields.Instance,
				Cluster:         tt.fields.Cluster,
				ClusterSnapshot: tt.fields.ClusterSnapshot,
			}
			got := r.GetClusterSnapshotIdentifier()
			if tt.want == nil {
				if got != tt.want {
					t.Errorf("RDSRestorationStore.GetClusterSnapshotIdentifier() = %v, want %v", got, tt.want)
				}
			} else {
				if *got != *tt.want {
					t.Errorf("RDSRestorationStore.GetClusterSnapshotIdentifier() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestRDSRestorationStore_GetDBClusterIdentifier(t *testing.T) {
	type fields struct {
		Snapshot        *types.DBSnapshot
		Instance        *types.DBInstance
		Cluster         *types.DBCluster
		ClusterSnapshot *types.DBClusterSnapshot
	}
	tests := []struct {
		name   string
		fields fields
		want   *string
	}{
		{name: "totalNil", fields: fields{nil, nil, nil, nil}, want: nil},
		{name: "RegularNil", fields: fields{nil, nil, nil, &types.DBClusterSnapshot{}}, want: nil},
		{name: "Value", fields: fields{nil, nil, nil, &types.DBClusterSnapshot{DBClusterIdentifier: aws.String("foo")}}, want: aws.String("foo")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RDSRestorationStore{
				Snapshot:        tt.fields.Snapshot,
				Instance:        tt.fields.Instance,
				Cluster:         tt.fields.Cluster,
				ClusterSnapshot: tt.fields.ClusterSnapshot,
			}
			got := r.GetDBClusterIdentifier()
			if tt.want == nil {
				if got != tt.want {
					t.Errorf("RDSRestorationStore.GetDBClusterIdentifier() = %v, want %v", got, tt.want)
				}
			} else {
				if *got != *tt.want {
					t.Errorf("RDSRestorationStore.GetDBClusterIdentifier() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestRDSRestorationStore_GetClusterEngine(t *testing.T) {
	type fields struct {
		Snapshot        *types.DBSnapshot
		Instance        *types.DBInstance
		Cluster         *types.DBCluster
		ClusterSnapshot *types.DBClusterSnapshot
	}
	tests := []struct {
		name   string
		fields fields
		want   *string
	}{
		{name: "totalNil", fields: fields{nil, nil, nil, nil}, want: nil},
		{name: "RegularNil", fields: fields{nil, nil, nil, &types.DBClusterSnapshot{}}, want: nil},
		{name: "Value", fields: fields{nil, nil, nil, &types.DBClusterSnapshot{Engine: aws.String("foo")}}, want: aws.String("foo")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RDSRestorationStore{
				Snapshot:        tt.fields.Snapshot,
				Instance:        tt.fields.Instance,
				Cluster:         tt.fields.Cluster,
				ClusterSnapshot: tt.fields.ClusterSnapshot,
			}
			got := r.GetClusterEngine()
			if tt.want == nil {
				if got != tt.want {
					t.Errorf("RDSRestorationStore.GetClusterEngine() = %v, want %v", got, tt.want)
				}
			} else {
				if *got != *tt.want {
					t.Errorf("RDSRestorationStore.GetClusterEngine() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestRDSRestorationStore_GetKmsKey(t *testing.T) {
	type fields struct {
		Snapshot        *types.DBSnapshot
		Instance        *types.DBInstance
		Cluster         *types.DBCluster
		ClusterSnapshot *types.DBClusterSnapshot
	}
	tests := []struct {
		name   string
		fields fields
		want   *string
	}{
		{name: "totalNil", fields: fields{nil, nil, nil, nil}, want: nil},
		{name: "RegularNil", fields: fields{nil, nil, nil, &types.DBClusterSnapshot{}}, want: nil},
		{name: "Value", fields: fields{nil, nil, nil, &types.DBClusterSnapshot{KmsKeyId: aws.String("foo")}}, want: aws.String("foo")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RDSRestorationStore{
				Snapshot:        tt.fields.Snapshot,
				Instance:        tt.fields.Instance,
				Cluster:         tt.fields.Cluster,
				ClusterSnapshot: tt.fields.ClusterSnapshot,
			}
			got := r.GetKmsKey()
			if tt.want == nil {
				if got != tt.want {
					t.Errorf("RDSRestorationStore.GetKmsKey() = %v, want %v", got, tt.want)
				}
			} else {
				if *got != *tt.want {
					t.Errorf("RDSRestorationStore.GetKmsKey() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestRDSRestorationStore_GetClusterAZs(t *testing.T) {
	type fields struct {
		Snapshot        *types.DBSnapshot
		Instance        *types.DBInstance
		Cluster         *types.DBCluster
		ClusterSnapshot *types.DBClusterSnapshot
	}
	tests := []struct {
		name   string
		fields fields
		want   *[]string
	}{
		{name: "totalNil", fields: fields{nil, nil, nil, nil}, want: nil},
		{name: "RegularNil", fields: fields{nil, nil, nil, &types.DBClusterSnapshot{}}, want: nil},
		{name: "Value", fields: fields{nil, nil, nil, &types.DBClusterSnapshot{AvailabilityZones: []string{"foo"}}}, want: &[]string{"foo"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RDSRestorationStore{
				Snapshot:        tt.fields.Snapshot,
				Instance:        tt.fields.Instance,
				Cluster:         tt.fields.Cluster,
				ClusterSnapshot: tt.fields.ClusterSnapshot,
			}
			got := r.GetClusterAZs()
			if tt.want == nil {
				if got != tt.want {
					t.Errorf("RDSRestorationStore.GetClusterAZs() = %v, want %v", got, tt.want)
				}
			} else {
				gotSlice := *got
				gotZero := gotSlice[0]

				wantSlice := *tt.want
				wantZero := wantSlice[0]
				if gotZero != wantZero {
					t.Errorf("RDSRestorationStore.GetClusterAZs() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestRDSRestorationStore_GetDBClusterInstanceClass(t *testing.T) {
	type fields struct {
		Snapshot        *types.DBSnapshot
		Instance        *types.DBInstance
		Cluster         *types.DBCluster
		ClusterSnapshot *types.DBClusterSnapshot
	}
	tests := []struct {
		name   string
		fields fields
		want   *string
	}{
		{name: "totalNil", fields: fields{nil, nil, nil, nil}, want: nil},
		{name: "RegularNil", fields: fields{nil, nil, nil, &types.DBClusterSnapshot{}}, want: nil},
		{name: "Value", fields: fields{nil, nil, &types.DBCluster{DBClusterInstanceClass: aws.String("foo")}, nil}, want: aws.String("foo")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RDSRestorationStore{
				Snapshot:        tt.fields.Snapshot,
				Instance:        tt.fields.Instance,
				Cluster:         tt.fields.Cluster,
				ClusterSnapshot: tt.fields.ClusterSnapshot,
			}
			got := r.GetDBClusterInstanceClass()
			if tt.want == nil {
				if got != tt.want {
					t.Errorf("RDSRestorationStore.GetDBClusterInstanceClass() = %v, want %v", got, tt.want)
				}
			} else {
				if *got != *tt.want {
					t.Errorf("RDSRestorationStore.GetDBClusterInstanceClass() = %v, want %v", *got, *tt.want)
				}
			}
		})
	}
}

func TestRDSRestorationStore_GetAutoMinorVersionUpgrade(t *testing.T) {
	pTrue := true
	type fields struct {
		Snapshot        *types.DBSnapshot
		Instance        *types.DBInstance
		Cluster         *types.DBCluster
		ClusterSnapshot *types.DBClusterSnapshot
	}
	tests := []struct {
		name   string
		fields fields
		want   *bool
	}{
		{name: "totalNil", fields: fields{nil, nil, nil, nil}, want: nil},
		{name: "GetData", fields: fields{nil, &types.DBInstance{AutoMinorVersionUpgrade: &pTrue}, nil, nil}, want: &pTrue},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RDSRestorationStore{
				Snapshot:        tt.fields.Snapshot,
				Instance:        tt.fields.Instance,
				Cluster:         tt.fields.Cluster,
				ClusterSnapshot: tt.fields.ClusterSnapshot,
			}
			got := r.GetAutoMinorVersionUpgrade()
			if got == nil {
				if got != tt.want {
					t.Errorf("RDSRestorationStore.GetAutoMinorVersionUpgrade() = %v, want %v", got, tt.want)
				}
			} else {
				if *got != *tt.want {
					t.Errorf("RDSRestorationStore.GetAutoMinorVersionUpgrade() = %v, want %v", *got, *tt.want)
				}
			}
		})
	}
}

func TestRDSRestorationStore_GetBackupTarget(t *testing.T) {
	type fields struct {
		Snapshot        *types.DBSnapshot
		Instance        *types.DBInstance
		Cluster         *types.DBCluster
		ClusterSnapshot *types.DBClusterSnapshot
	}
	tests := []struct {
		name   string
		fields fields
		want   *string
	}{
		{name: "totalNil", fields: fields{nil, nil, nil, nil}, want: nil},
		{name: "RegularNil", fields: fields{nil, &types.DBInstance{}, nil, nil}, want: nil},
		{name: "GetData", fields: fields{nil, &types.DBInstance{BackupTarget: aws.String("foo")}, nil, nil}, want: aws.String("foo")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RDSRestorationStore{
				Snapshot:        tt.fields.Snapshot,
				Instance:        tt.fields.Instance,
				Cluster:         tt.fields.Cluster,
				ClusterSnapshot: tt.fields.ClusterSnapshot,
			}
			got := r.GetBackupTarget()
			if got == nil {
				if got != tt.want {
					t.Errorf("RDSRestorationStore.GetBackupTarget() = %v, want %v", got, tt.want)
				}
			} else {
				if *got != *tt.want {
					t.Errorf("RDSRestorationStore.GetBackupTarget() = %v, want %v", *got, *tt.want)
				}
			}
		})
	}
}

func TestRDSRestorationStore_GetDeleteProtection(t *testing.T) {
	pTrue := true
	type fields struct {
		Snapshot        *types.DBSnapshot
		Instance        *types.DBInstance
		Cluster         *types.DBCluster
		ClusterSnapshot *types.DBClusterSnapshot
	}
	tests := []struct {
		name   string
		fields fields
		want   *bool
	}{
		{name: "totalNil", fields: fields{nil, nil, nil, nil}, want: nil},
		{name: "GetData", fields: fields{nil, &types.DBInstance{DeletionProtection: &pTrue}, nil, nil}, want: &pTrue},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RDSRestorationStore{
				Snapshot:        tt.fields.Snapshot,
				Instance:        tt.fields.Instance,
				Cluster:         tt.fields.Cluster,
				ClusterSnapshot: tt.fields.ClusterSnapshot,
			}
			got := r.GetDeleteProtection()
			if got == nil {
				if got != tt.want {
					t.Errorf("RDSRestorationStore.GetDeleteProtection() = %v, want %v", got, tt.want)
				}
			} else {
				if *got != *tt.want {
					t.Errorf("RDSRestorationStore.GetDeleteProtection() = %v, want %v", *got, *tt.want)
				}
			}
		})
	}
}

func TestRDSRestorationStore_GetEnabledCloudwatchLogsExports(t *testing.T) {
	type fields struct {
		Snapshot        *types.DBSnapshot
		Instance        *types.DBInstance
		Cluster         *types.DBCluster
		ClusterSnapshot *types.DBClusterSnapshot
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{name: "totalNil", fields: fields{nil, nil, nil, nil}, want: []string{}},
		{name: "RegularNil", fields: fields{nil, &types.DBInstance{}, nil, nil}, want: []string{}},
		{name: "GetData", fields: fields{nil, &types.DBInstance{EnabledCloudwatchLogsExports: []string{"error"}}, nil, nil}, want: []string{"error"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RDSRestorationStore{
				Snapshot:        tt.fields.Snapshot,
				Instance:        tt.fields.Instance,
				Cluster:         tt.fields.Cluster,
				ClusterSnapshot: tt.fields.ClusterSnapshot,
			}
			if got := r.GetEnabledCloudwatchLogsExports(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RDSRestorationStore.GetEnabledCloudwatchLogsExports() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRDSRestorationStore_GetDBClusterMembers(t *testing.T) {
	type fields struct {
		Snapshot        *types.DBSnapshot
		Instance        *types.DBInstance
		Cluster         *types.DBCluster
		ClusterSnapshot *types.DBClusterSnapshot
	}
	pTrue := true
	tests := []struct {
		name   string
		fields fields
		want   *[]types.DBClusterMember
	}{
		{name: "totalNil", fields: fields{nil, nil, nil, nil}, want: nil},
		{name: "RegularNil", fields: fields{nil, nil, &types.DBCluster{}, nil}, want: nil},
		{name: "GetData", fields: fields{nil, &types.DBInstance{EnabledCloudwatchLogsExports: []string{"error"}}, &types.DBCluster{DBClusterMembers: []types.DBClusterMember{{IsClusterWriter: &pTrue}}}, nil}, want: &[]types.DBClusterMember{{IsClusterWriter: &pTrue}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RDSRestorationStore{
				Snapshot:        tt.fields.Snapshot,
				Instance:        tt.fields.Instance,
				Cluster:         tt.fields.Cluster,
				ClusterSnapshot: tt.fields.ClusterSnapshot,
			}
			if got := r.GetDBClusterMembers(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RDSRestorationStore.GetDBClusterMembers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRDSRestorationStore_GetParameterGroups(t *testing.T) {
	type fields struct {
		Snapshot        *types.DBSnapshot
		Instance        *types.DBInstance
		Cluster         *types.DBCluster
		ClusterSnapshot *types.DBClusterSnapshot
	}
	pTrue := true
	tests := []struct {
		name   string
		fields fields
		want   []types.DBParameterGroupStatus
	}{
		{name: "totalNil", fields: fields{nil, nil, nil, nil}, want: nil},
		{name: "RegularNil", fields: fields{nil, nil, &types.DBCluster{}, nil}, want: nil},
		{name: "GetData", fields: fields{nil, &types.DBInstance{DBParameterGroups: []types.DBParameterGroupStatus{}}, &types.DBCluster{DBClusterMembers: []types.DBClusterMember{{IsClusterWriter: &pTrue}}}, nil}, want: []types.DBParameterGroupStatus{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RDSRestorationStore{
				Snapshot:        tt.fields.Snapshot,
				Instance:        tt.fields.Instance,
				Cluster:         tt.fields.Cluster,
				ClusterSnapshot: tt.fields.ClusterSnapshot,
			}
			if got := r.GetParameterGroups(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RDSRestorationStore.GetParameterGroups() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRDSRestorationStore_GetClusterParameterGroups(t *testing.T) {
	type fields struct {
		Snapshot        *types.DBSnapshot
		Instance        *types.DBInstance
		Cluster         *types.DBCluster
		ClusterSnapshot *types.DBClusterSnapshot
	}
	field := fields{nil, &types.DBInstance{}, &types.DBCluster{DBClusterParameterGroup: aws.String("foo")}, nil}
	tests := []struct {
		name   string
		fields fields
		want   *string
	}{
		{name: "totalNil", fields: fields{nil, nil, nil, nil}, want: nil},
		{name: "RegularNil", fields: fields{nil, nil, &types.DBCluster{}, nil}, want: nil},
		{name: "GetData", fields: field, want: aws.String("foo")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RDSRestorationStore{
				Snapshot:        tt.fields.Snapshot,
				Instance:        tt.fields.Instance,
				Cluster:         tt.fields.Cluster,
				ClusterSnapshot: tt.fields.ClusterSnapshot,
			}
			got := r.GetClusterParameterGroups()
			if got != nil && tt.want != nil {
				if *got != *tt.want {
					t.Errorf("RDSRestorationStore.GetClusterParameterGroups() = %v, want %v", *got, *tt.want)
				}
			} else {
				if got != tt.want {
					t.Errorf("RDSRestorationStore.GetClusterParameterGroups() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
