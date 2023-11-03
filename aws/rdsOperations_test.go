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
	"github.com/jrottersman/lats/stack"
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
		c CreateInstanceFromStackInput
	}
	field := fields{RdsClient: mock.MockRDSClient{}}
	// Create long args
	objs := []stack.Object{}
	obj1 := stack.Object{}
	obj2 := stack.Object{}
	objs = append(objs, obj1)
	objs = append(objs, obj2)
	objects := make(map[int][]stack.Object)
	objects[2] = objs
	longStack := stack.Stack{
		Objects: objects,
	}
	c := CreateInstanceFromStackInput{
		Stack:  &longStack,
		DBName: aws.String("foo"),
	}
	failArg := args{c: c}

	//Create a valid object and instance
	filename := "/tmp/foo"
	order := 2
	objType := stack.LoneInstance

	defer os.Remove(filename)
	db := rds.RestoreDBInstanceFromDBSnapshotInput{
		DBInstanceIdentifier: aws.String("foo"),
	}
	r := state.EncodeRestoreDBInstanceFromDBSnapshotInput(&db)
	_, err := state.WriteOutput(filename, r)
	if err != nil {
		t.Errorf("failed to write output, %s", err)
	}
	gObj := stack.NewObject(filename, order, objType)
	objs2 := []stack.Object{}
	objs2 = append(objs2, gObj)
	object := make(map[int][]stack.Object)
	object[2] = objs2
	goodstack := stack.Stack{
		Objects: object,
	}
	c2 := CreateInstanceFromStackInput{
		Stack:  &goodstack,
		DBName: aws.String("foo"),
	}
	passArgs := args{c: c2}

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
			if err := instances.CreateInstanceFromStack(tt.args.c); (err != nil) != tt.wantErr {
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
		c CreateClusterFromStackInput
	}

	field := fields{RdsClient: mock.MockRDSClient{}}
	// Create long args
	objs := []stack.Object{}
	obj1 := stack.Object{}
	obj2 := stack.Object{}
	objs = append(objs, obj1)
	objs = append(objs, obj2)
	objects := make(map[int][]stack.Object)
	objects[2] = objs
	longStack := stack.Stack{
		Objects: objects,
	}
	c := CreateClusterFromStackInput{
		S: &longStack,
	}
	failArg := args{c: c}

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
			if err := instances.CreateClusterFromStack(tt.args.c); (err != nil) != tt.wantErr {
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

func TestDbInstances_GetParametersForClusterParameterGroup(t *testing.T) {
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
			got, err := instances.GetParametersForClusterParameterGroup(tt.args.ParameterGroupName)
			if (err != nil) != tt.wantErr {
				t.Errorf("DbInstances.GetParametersForClusterParameterGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DbInstances.GetParametersForClusterParameterGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDbInstances_CreateParameterGroup(t *testing.T) {
	type fields struct {
		RdsClient Client
	}
	type args struct {
		p *types.DBParameterGroup
	}
	params := types.DBParameterGroup{
		DBParameterGroupArn:    aws.String("foo"),
		DBParameterGroupFamily: aws.String("foo"),
		DBParameterGroupName:   aws.String("foo"),
	}
	expect := rds.CreateDBParameterGroupOutput{
		DBParameterGroup: &params,
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *rds.CreateDBParameterGroupOutput
		wantErr bool
	}{
		{name: "foo", fields: fields{RdsClient: mock.MockRDSClient{}}, args: args{p: &params}, want: &expect, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &DbInstances{
				RdsClient: tt.fields.RdsClient,
			}
			got, err := i.CreateParameterGroup(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("DbInstances.CreateParameterGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DbInstances.CreateParameterGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDbInstances_ModifyParameterGroup(t *testing.T) {
	type fields struct {
		RdsClient Client
	}
	type args struct {
		pg         string
		parameters []types.Parameter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"emptyPass", fields{RdsClient: mock.MockRDSClient{}}, args{"foo", []types.Parameter{}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instances := &DbInstances{
				RdsClient: tt.fields.RdsClient,
			}
			if err := instances.ModifyParameterGroup(tt.args.pg, tt.args.parameters); (err != nil) != tt.wantErr {
				t.Errorf("DbInstances.ModifyParameterGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDbInstances_CreateClusterParameterGroup(t *testing.T) {
	type fields struct {
		RdsClient Client
	}
	type args struct {
		p *types.DBClusterParameterGroup
	}

	params := types.DBClusterParameterGroup{
		DBClusterParameterGroupArn:  aws.String("foo"),
		DBParameterGroupFamily:      aws.String("foo"),
		DBClusterParameterGroupName: aws.String("foo"),
	}
	expect := rds.CreateDBClusterParameterGroupOutput{
		DBClusterParameterGroup: &params,
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *rds.CreateDBClusterParameterGroupOutput
		wantErr bool
	}{
		{name: "foo", fields: fields{RdsClient: mock.MockRDSClient{}}, args: args{p: &params}, want: &expect, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instances := &DbInstances{
				RdsClient: tt.fields.RdsClient,
			}
			got, err := instances.CreateClusterParameterGroup(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("DbInstances.CreateClusterParameterGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DbInstances.CreateClusterParameterGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDbInstances_ModifyClusterParameterGroup(t *testing.T) {
	type fields struct {
		RdsClient Client
	}
	type args struct {
		pg         string
		parameters []types.Parameter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"emptyPass", fields{RdsClient: mock.MockRDSClient{}}, args{"foo", []types.Parameter{}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instances := &DbInstances{
				RdsClient: tt.fields.RdsClient,
			}
			if err := instances.ModifyClusterParameterGroup(tt.args.pg, tt.args.parameters); (err != nil) != tt.wantErr {
				t.Errorf("DbInstances.ModifyClusterParameterGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDbInstances_RestoreOptionGroup(t *testing.T) {
	type fields struct {
		RdsClient Client
	}
	type args struct {
		EngineName         string
		MajorEngineVersion string
		OptionGroupName    string
		Description        string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *rds.CreateOptionGroupOutput
		wantErr bool
	}{
		{name: "emptyTest",
			fields:  fields{RdsClient: mock.MockRDSClient{}},
			args:    args{EngineName: "foo", MajorEngineVersion: "bar", OptionGroupName: "bat", Description: "taz"},
			want:    &rds.CreateOptionGroupOutput{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instances := &DbInstances{
				RdsClient: tt.fields.RdsClient,
			}
			got, err := instances.RestoreOptionGroup(tt.args.EngineName, tt.args.MajorEngineVersion, tt.args.OptionGroupName, tt.args.Description)
			if (err != nil) != tt.wantErr {
				t.Errorf("DbInstances.RestoreOptionGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DbInstances.RestoreOptionGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDbInstances_ModifyOptionGroup(t *testing.T) {
	type fields struct {
		RdsClient Client
	}
	type args struct {
		OptionGroupName string
		Include         []types.OptionConfiguration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "pass",
			fields:  fields{RdsClient: mock.MockRDSClient{}},
			args:    args{OptionGroupName: "foo", Include: []types.OptionConfiguration{}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instances := &DbInstances{
				RdsClient: tt.fields.RdsClient,
			}
			if err := instances.ModifyOptionGroup(tt.args.OptionGroupName, tt.args.Include); (err != nil) != tt.wantErr {
				t.Errorf("DbInstances.ModifyOptionGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDbInstances_GetOptionGroup(t *testing.T) {
	type fields struct {
		RdsClient Client
	}
	type args struct {
		OptionGroupName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.OptionGroup
		wantErr bool
	}{
		{
			name:    "pass",
			fields:  fields{RdsClient: mock.MockRDSClient{}},
			args:    args{OptionGroupName: "foo"},
			want:    &types.OptionGroup{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instances := &DbInstances{
				RdsClient: tt.fields.RdsClient,
			}
			got, err := instances.GetOptionGroup(tt.args.OptionGroupName)
			if (err != nil) != tt.wantErr {
				t.Errorf("DbInstances.GetOptionGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DbInstances.GetOptionGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_optionsToConfiguration(t *testing.T) {
	type args struct {
		opts []types.Option
	}
	tests := []struct {
		name string
		args args
		want []types.OptionConfiguration
	}{
		{"empty", args{[]types.Option{}}, []types.OptionConfiguration{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := optionsToConfiguration(tt.args.opts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("optionsToConfiguration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDbInstances_CreateDBSubnetGroup(t *testing.T) {
	type fields struct {
		RdsClient Client
	}
	type args struct {
		name        string
		description string
		subnets     []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *rds.CreateDBSubnetGroupOutput
		wantErr bool
	}{
		{name: "nil", fields: fields{mock.MockRDSClient{}}, args: args{name: "", description: "", subnets: []string{}}, want: &rds.CreateDBSubnetGroupOutput{}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instances := &DbInstances{
				RdsClient: tt.fields.RdsClient,
			}
			got, err := instances.CreateDBSubnetGroup(tt.args.name, tt.args.description, tt.args.subnets)
			if (err != nil) != tt.wantErr {
				t.Errorf("DbInstances.CreateDBSubnetGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DbInstances.CreateDBSubnetGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}
