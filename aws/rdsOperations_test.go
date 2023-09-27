package aws

import (
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	mock "github.com/jrottersman/lats/mocks"
	"github.com/jrottersman/lats/state"
)

func TestGetCluster(t *testing.T) {
	expected := "foo"
	c := mock.MockRDSClient{}
	dbi := DbInstances{
		RdsClient: c,
	}
	resp, err := dbi.GetCluster(expected)
	if err != nil {
		t.Errorf("got error %s", err)
	}
	if *resp.DBClusterIdentifier != expected {
		t.Errorf("got %s expected %s", *resp.DBClusterIdentifier, expected)
	}
}
func TestGetInstance(t *testing.T) {
	expected := "foo"
	c := mock.MockRDSClient{}
	dbi := DbInstances{
		RdsClient: c,
	}
	resp, err := dbi.GetInstance("foo")
	if err != nil {
		t.Errorf("got the following error %v", err)
	}
	if *resp.DBInstanceIdentifier != expected {
		t.Errorf("got %s expected %s", *resp.DBInstanceIdentifier, expected)
	}
}

func TestGetInstanceFromCluster(t *testing.T) {
	expected := "foo"
	c := mock.MockRDSClient{}
	dbi := DbInstances{
		RdsClient: c,
	}
	resp, err := dbi.GetCluster(expected)
	if err != nil {
		t.Errorf("got error %s", err)
	}
	r, err := dbi.GetInstancesFromCluster(resp)
	if err != nil {
		t.Errorf("got error %s", err)
	}
	if *r[0].DBInstanceIdentifier != expected {
		t.Errorf("got %s expected %s", *r[0].DBInstanceIdentifier, expected)
	}
}

func TestCreateSnapshot(t *testing.T) {
	c := mock.MockRDSClient{}
	dbi := DbInstances{
		RdsClient: c,
	}
	resp, err := dbi.CreateSnapshot("foo", "bar")
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if resp.AllocatedStorage != 1000 {
		t.Errorf("got %d expected 1000", resp.AllocatedStorage)
	}
}

func TestCreateClusterSnapshot(t *testing.T) {
	c := mock.MockRDSClient{}
	dbi := DbInstances{
		RdsClient: c,
	}
	resp, err := dbi.CreateClusterSnapshot("foo", "bar")
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if resp.AllocatedStorage != 1000 {
		t.Errorf("got %d expected 1000", resp.AllocatedStorage)
	}
}

func TestDescribeParameterGroup(t *testing.T) {
	c := mock.MockRDSClient{}
	dbi := DbInstances{
		RdsClient: c,
	}
	resp, err := dbi.GetParameterGroup("foo")
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if *resp.DBParameterGroupName != "foo" {
		t.Errorf("expected foo got %s", *resp.DBParameterGroupName)
	}
}

func TestCopySnapshot(t *testing.T) {
	c := mock.MockRDSClient{}
	dbi := DbInstances{
		RdsClient: c,
	}
	resp, err := dbi.CopySnapshot("foo", "bar", "us-east-1", "keyArn")
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if resp.AllocatedStorage != 1000 {
		t.Errorf("got %d expected 1000", resp.AllocatedStorage)
	}
}

func TestRestoreSnapshotInstance(t *testing.T) {

	filename := "/tmp/foo"
	snap := types.DBSnapshot{
		AllocatedStorage:     1000,
		Encrypted:            true,
		PercentProgress:      100,
		DBInstanceIdentifier: aws.String("foobar"),
		DBSnapshotIdentifier: aws.String("foo"),
	}

	defer os.Remove(filename)
	r := state.EncodeRDSSnapshotOutput(&snap)
	_, err := state.WriteOutput(filename, r)
	if err != nil {
		t.Errorf("error writing file %s", err)
	}

	var s []state.StateKV
	kv := state.StateKV{
		Object:       "foo",
		FileLocation: filename,
		ObjectType:   "RDSSnapshot",
	}
	s = append(s, kv)

	filename2 := "/tmp/foo2"
	dbz := types.DBInstance{
		AllocatedStorage:     1000,
		DBInstanceIdentifier: aws.String("foobar"),
	}

	defer os.Remove(filename2)
	r2 := state.EncodeRDSDatabaseOutput(&dbz)
	_, err = state.WriteOutput(filename2, r2)
	if err != nil {
		t.Errorf("error writing file %s", err)
	}
	kv2 := state.StateKV{
		Object:       "foobar",
		FileLocation: filename2,
		ObjectType:   state.RdsInstanceType,
	}
	s = append(s, kv2)
	sm := state.StateManager{
		Mu:             &sync.Mutex{},
		StateLocations: s,
	}

	resp, err := state.RDSRestorationStoreBuilder(sm, "foo")
	if err != nil {
		t.Errorf("error writing file %s", err)
	}

	c := mock.MockRDSClient{}
	dbi := DbInstances{
		RdsClient: c,
	}
	input := state.GenerateRestoreDBInstanceFromDBSnapshotInput(*resp)
	resp2, err := dbi.RestoreSnapshotInstance(*input)
	if err != nil {
		t.Errorf("got error: %s", err)
	}
	if reflect.TypeOf(resp2) != reflect.TypeOf(&rds.RestoreDBInstanceFromDBSnapshotOutput{}) {
		t.Error()
	}
}

func TestRestoreSnapshotCluster(t *testing.T) {

	snap := types.DBClusterSnapshot{
		AllocatedStorage:            1000,
		PercentProgress:             100,
		DBClusterIdentifier:         aws.String("foobar"),
		DBClusterSnapshotIdentifier: aws.String("foo"),
	}

	dbz := types.DBCluster{
		DBClusterIdentifier: aws.String("foobar"),
	}

	store := state.RDSRestorationStore{
		Cluster:         &dbz,
		ClusterSnapshot: &snap,
	}

	c := mock.MockRDSClient{}
	dbi := DbInstances{
		RdsClient: c,
	}
	input := state.GenerateRestoreDBInstanceFromDBClusterSnapshotInput(store)
	resp, err := dbi.RestoreSnapshotInstance(*input)
	if err != nil {
		t.Errorf("got error: %s", err)
	}
	if reflect.TypeOf(resp) != reflect.TypeOf(&rds.RestoreDBInstanceFromDBSnapshotOutput{}) {
		t.Errorf("got %T expected %T", resp, &rds.RestoreDBInstanceFromDBSnapshotOutput{})
	}
}

func TestDbInstances_getClusterStatus(t *testing.T) {
	type fields struct {
		RdsClient Client
	}
	type args struct {
		name string
	}
	m := mock.MockRDSClient{}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *string
		wantErr bool
	}{
		{name: "test", fields: fields{RdsClient: m}, args: args{"foo"}, want: aws.String("creating"), wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instances := &DbInstances{
				RdsClient: tt.fields.RdsClient,
			}
			got, err := instances.getClusterStatus(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("DbInstances.getClusterStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if *got != *tt.want {
				t.Errorf("DbInstances.getClusterStatus() = %v want %v", *got, *tt.want)
			}
		})
	}
}

func TestDbInstances_GetInstanceSnapshotARN(t *testing.T) {
	type fields struct {
		RdsClient Client
	}
	type args struct {
		name   string
		marker *string
	}

	m := mock.MockRDSClient{}
	field := fields{RdsClient: m}
	arg := args{"foo", nil}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *string
		wantErr bool
	}{
		{name: "test", fields: field, args: arg, want: aws.String("foo"), wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instances := &DbInstances{
				RdsClient: tt.fields.RdsClient,
			}
			got, err := instances.GetInstanceSnapshotARN(tt.args.name, tt.args.marker)
			if (err != nil) != tt.wantErr {
				t.Errorf("DbInstances.GetSnapshotARN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if *got != *tt.want {
				t.Errorf("DbInstances.GetSnapshotARN() = %v, want %v", *got, *tt.want)
			}
		})
	}
}

func TestDbInstances_GetClusterSnapshotARN(t *testing.T) {
	type fields struct {
		RdsClient Client
	}
	type args struct {
		name   string
		marker *string
	}

	m := mock.MockRDSClient{}
	field := fields{RdsClient: m}
	arg := args{"foo", nil}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *string
		wantErr bool
	}{
		{name: "test", fields: field, args: arg, want: aws.String("foo"), wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instances := &DbInstances{
				RdsClient: tt.fields.RdsClient,
			}
			got, err := instances.GetClusterSnapshotARN(tt.args.name, tt.args.marker)
			if (err != nil) != tt.wantErr {
				t.Errorf("DbInstances.GetClusterSnapshotARN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if *got != *tt.want {
				t.Errorf("DbInstances.GetClusterSnapshotARN() = %v, want %v", *got, *tt.want)
			}
		})
	}
}

func TestDbInstances_GetSnapshotARN(t *testing.T) {
	type fields struct {
		RdsClient Client
	}
	type args struct {
		name    string
		cluster bool
	}

	m := mock.MockRDSClient{}
	field := fields{RdsClient: m}
	arg := args{"foo", false}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *string
		wantErr bool
	}{
		{name: "test", fields: field, args: arg, want: aws.String("foo"), wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instances := &DbInstances{
				RdsClient: tt.fields.RdsClient,
			}
			got, err := instances.GetSnapshotARN(tt.args.name, tt.args.cluster)
			if (err != nil) != tt.wantErr {
				t.Errorf("DbInstances.GetSnapshotARN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if *got != *tt.want {
				t.Errorf("DbInstances.GetSnapshotARN() = %v, want %v", *got, *tt.want)
			}
		})
	}
}

func TestDbInstances_CopyClusterSnaphot(t *testing.T) {
	type fields struct {
		RdsClient Client
	}
	type args struct {
		originalSnapshotName string
		newSnapshotName      string
		sourceRegion         string
		kmsKey               string
	}
	m := mock.MockRDSClient{}
	field := fields{RdsClient: m}
	want := types.DBClusterSnapshot{}
	arg := args{originalSnapshotName: "foo", newSnapshotName: "foo", sourceRegion: "baz", kmsKey: "bat"}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.DBClusterSnapshot
		wantErr bool
	}{
		{name: "test", fields: field, args: arg, want: &want, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instances := &DbInstances{
				RdsClient: tt.fields.RdsClient,
			}
			got, err := instances.CopyClusterSnaphot(tt.args.originalSnapshotName, tt.args.newSnapshotName, tt.args.sourceRegion, tt.args.kmsKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("DbInstances.CopyClusterSnaphot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DbInstances.CopyClusterSnaphot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDbInstances_CreateInstanceFromStack(t *testing.T) {
	type fields struct {
		RdsClient Client
	}
	type args struct {
		s *state.Stack
	}
	field := fields{RdsClient: mock.MockRDSClient{}}
	// Create long args
	objs := []state.Object{}
	obj1 := state.Object{}
	obj2 := state.Object{}
	objs = append(objs, obj1)
	objs = append(objs, obj2)
	objects := make(map[int][]state.Object)
	objects[1] = objs
	longStack := state.Stack{
		Objects: objects,
	}
	failArg := args{s: &longStack}

	//Create a valid object and instance
	filename := "/tmp/foo"
	order := 1
	objType := state.LoneInstance

	defer os.Remove(filename)
	db := rds.RestoreDBInstanceFromDBSnapshotInput{
		DBInstanceIdentifier: aws.String("foo"),
	}
	r := state.EncodeRestoreDBInstanceFromDBSnapshotInput(&db)
	_, err := state.WriteOutput(filename, r)
	if err != nil {
		t.Errorf("failed to write output, %s", err)
	}
	gObj := state.NewObject(filename, order, objType)
	objs2 := []state.Object{}
	objs2 = append(objs2, gObj)
	object := make(map[int][]state.Object)
	object[1] = objs2
	goodstack := state.Stack{
		Objects: object,
	}
	passArgs := args{s: &goodstack}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "fail", fields: field, args: failArg, wantErr: true},
		{name: "pass", fields: field, args: passArgs, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instances := &DbInstances{
				RdsClient: tt.fields.RdsClient,
			}
			if err := instances.CreateInstanceFromStack(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("DbInstances.CreateInstanceFromStack() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDbInstances_RestoreInstanceForCluster(t *testing.T) {
	type fields struct {
		RdsClient Client
	}
	type args struct {
		input rds.CreateDBInstanceInput
	}

	field := fields{RdsClient: mock.MockRDSClient{}}
	arg := args{input: rds.CreateDBInstanceInput{}}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *rds.CreateDBInstanceOutput
		wantErr bool
	}{
		{name: "good", fields: field, args: arg, want: &rds.CreateDBInstanceOutput{}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instances := &DbInstances{
				RdsClient: tt.fields.RdsClient,
			}
			got, err := instances.RestoreInstanceForCluster(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("DbInstances.RestoreInstanceForCluster() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DbInstances.RestoreInstanceForCluster() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDbInstances_CreateClusterFromStack(t *testing.T) {
	type fields struct {
		RdsClient Client
	}
	type args struct {
		s *state.Stack
	}

	field := fields{RdsClient: mock.MockRDSClient{}}
	// Create long args
	objs := []state.Object{}
	obj1 := state.Object{}
	obj2 := state.Object{}
	objs = append(objs, obj1)
	objs = append(objs, obj2)
	objects := make(map[int][]state.Object)
	objects[1] = objs
	longStack := state.Stack{
		Objects: objects,
	}
	failArg := args{s: &longStack}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "fail", fields: field, args: failArg, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instances := &DbInstances{
				RdsClient: tt.fields.RdsClient,
			}
			if err := instances.CreateClusterFromStack(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("DbInstances.CreateClusterFromStack() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDbInstances_GetParametersForGroup(t *testing.T) {
	type fields struct {
		RdsClient Client
	}
	type args struct {
		ParameterGroupName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *[]types.Parameter
		wantErr bool
	}{
		{name: "pass", fields: fields{RdsClient: mock.MockRDSClient{}}, args: args{ParameterGroupName: "foo"}, want: &[]types.Parameter{}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instances := &DbInstances{
				RdsClient: tt.fields.RdsClient,
			}
			got, err := instances.GetParametersForGroup(tt.args.ParameterGroupName)
			if (err != nil) != tt.wantErr {
				t.Errorf("DbInstances.GetParametersForGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DbInstances.GetParametersForGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDbInstances_GetClusterParameterGroup(t *testing.T) {
	c := mock.MockRDSClient{}
	dbi := DbInstances{
		RdsClient: c,
	}
	resp, err := dbi.GetClusterParameterGroup("foo")
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if *resp.DBClusterParameterGroupName != "foo" {
		t.Errorf("expected foo got %s", *resp.DBClusterParameterGroupName)
	}
}
