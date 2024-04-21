package state

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"log/slog"
	"os"

	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/jrottersman/lats/helpers"
)

// EncodeRDSDatabaseOutput converts a dbInstace to an array of bytes in preperation for wrtiing it to disk
func EncodeRDSDatabaseOutput(db *types.DBInstance) bytes.Buffer {
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)

	err := enc.Encode(db)
	if err != nil {
		slog.Error("Error encoding our database", "error", err)
	}
	return encoder
}

func EncodeCreateDBClusterInput(c *rds.CreateDBClusterInput) bytes.Buffer {
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)
	err := enc.Encode(&c)
	if err != nil {
		slog.Error("Error encoding our database", "error", err)
	}
	return encoder
}

func DecodeCreateDBClusterInput(b bytes.Buffer) *rds.CreateDBClusterInput {
	var dbCluster rds.CreateDBClusterInput
	dec := gob.NewDecoder(&b)
	err := dec.Decode(&dbCluster)
	if err != nil {
		slog.Error("Error decoding state for RDS Cluster", "error", err)
	}
	return &dbCluster
}

// EncodeOptionGroup convers an option group struct to bytes
func EncodeOptionGroup(og *types.OptionGroup) bytes.Buffer {
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)

	err := enc.Encode(og)
	if err != nil {
		slog.Error("Error encoding our option group", "error", err)
	}
	return encoder
}

// DecodeOptionGroup takes a bytes buffer and returns it to a option group
func DecodeOptionGroup(b bytes.Buffer) types.OptionGroup {
	var optionGroup types.OptionGroup
	dec := gob.NewDecoder(&b)
	err := dec.Decode(&optionGroup)
	if err != nil {
		slog.Error("Error decoding state for Option Group", "error", err)
	}
	return optionGroup
}

// EncodeSecurityGroup converts a security group struct to bytes
func EncodeSecurityGroup(s ec2types.SecurityGroup) bytes.Buffer {
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)

	err := enc.Encode(s)
	if err != nil {
		slog.Error("Error encoding our option group", "error", err)
	}
	return encoder
}

// DecodeSecurityGroup takes a bytes buffer and returns it to a option group
func DecodeSecurityGroup(b bytes.Buffer) ec2types.SecurityGroup {
	var securityGroup ec2types.SecurityGroup
	dec := gob.NewDecoder(&b)
	err := dec.Decode(&securityGroup)
	if err != nil {
		slog.Error("Error decoding state for Option Group", "error", err)
	}
	return securityGroup
}

// DecodeRDSClusterOutput takes a bytes buffer and returns it to a DbCluster type in preperation of restoring the database
func DecodeRDSClusterOutput(b bytes.Buffer) types.DBCluster {
	var dbCluster types.DBCluster
	dec := gob.NewDecoder(&b)
	err := dec.Decode(&dbCluster)
	if err != nil {
		slog.Error("Error decoding state for RDS Cluster", "error", err)
	}
	return dbCluster
}

// EncodeRDSClusterOutput takes a cluster snapshot and creates bytes
func EncodeRDSClusterOutput(db *types.DBCluster) bytes.Buffer {
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)

	err := enc.Encode(db)
	if err != nil {
		slog.Error("Error encoding our database", "error", err)
	}
	return encoder
}

// DecodeRDSDatabaseOutput takes a bytes buffer and returns it to a DbInstance type in preperation of restoring the database
func DecodeRDSDatabaseOutput(b bytes.Buffer) types.DBInstance {
	var dbInstance types.DBInstance
	dec := gob.NewDecoder(&b)
	err := dec.Decode(&dbInstance)
	if err != nil {
		slog.Error("Error decoding state for RDS Instance", "error", err)
	}
	return dbInstance
}

// EncodeRDSSnapshotOutput converts a DbSnapshot struct to an array of bytes in preperation for wrtiing it to disk
func EncodeRDSSnapshotOutput(snapshot *types.DBSnapshot) bytes.Buffer {
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)

	err := enc.Encode(snapshot)
	if err != nil {
		slog.Error("Error encoding our snapshot", "error", err)
	}
	return encoder
}

// GetRDSSnapshotOutput reads a snapshot
func GetRDSSnapshotOutput(s StateManager, snap string) (*types.DBSnapshot, error) {
	i := s.GetStateObject(snap)
	snapshot, ok := i.(types.DBSnapshot)
	if !ok {
		str := fmt.Sprintf("error decoding snapshot from interface %v", i)
		return nil, errors.New(str)
	}
	return &snapshot, nil
}

// EncodeRDSClusterSnapshotOutput cluster output as bytes
func EncodeRDSClusterSnapshotOutput(snapshot *types.DBClusterSnapshot) bytes.Buffer {
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)

	err := enc.Encode(snapshot)
	if err != nil {
		slog.Error("Error encoding our snapshot", "error", err)
	}
	return encoder
}

// GenerateRestoreDBInstanceFromDBSnapshotInput create a db instance input
func GenerateRestoreDBInstanceFromDBSnapshotInput(r RDSRestorationStore) *rds.RestoreDBInstanceFromDBSnapshotInput {
	return &rds.RestoreDBInstanceFromDBSnapshotInput{
		DBInstanceClass:             r.GetInstanceClass(),
		DBInstanceIdentifier:        r.GetInstanceIdentifier(),
		AutoMinorVersionUpgrade:     r.GetAutoMinorVersionUpgrade(),
		AllocatedStorage:            r.GetAllocatedStorage(),
		BackupTarget:                r.GetBackupTarget(),
		DBSnapshotIdentifier:        r.GetSnapshotIdentifier(),
		DeletionProtection:          r.GetDeleteProtection(),
		EnableCloudwatchLogsExports: r.GetEnabledCloudwatchLogsExports(),
	}
}

// GenerateRestoreDBInstanceFromDBClusterSnapshotInput create a db cluster input
func GenerateRestoreDBInstanceFromDBClusterSnapshotInput(r RDSRestorationStore) *rds.RestoreDBInstanceFromDBSnapshotInput {
	return &rds.RestoreDBInstanceFromDBSnapshotInput{
		DBInstanceClass:             r.GetInstanceClass(),
		DBInstanceIdentifier:        r.GetInstanceIdentifier(),
		AutoMinorVersionUpgrade:     r.GetAutoMinorVersionUpgrade(),
		AllocatedStorage:            r.GetAllocatedStorage(),
		BackupTarget:                r.GetBackupTarget(),
		DBClusterSnapshotIdentifier: r.GetClusterSnapshotIdentifier(),
		DeletionProtection:          r.GetDeleteProtection(),
		EnableCloudwatchLogsExports: r.GetEnabledCloudwatchLogsExports(),
	}
}

// EncodeRestoreDBInstanceFromDBSnapshotInput encode snapshot as bytes
func EncodeRestoreDBInstanceFromDBSnapshotInput(r *rds.RestoreDBInstanceFromDBSnapshotInput) bytes.Buffer {
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)
	slog.Info("encoding db instance")
	err := enc.Encode(r)
	if err != nil {
		slog.Error("Error encoding our snapshot", "error", err)
	}
	return encoder
}

// DecodeRestoreDBInstanceFromDBSnapshotInput  decodes snapshot from bytes
func DecodeRestoreDBInstanceFromDBSnapshotInput(b bytes.Buffer) *rds.RestoreDBInstanceFromDBSnapshotInput {
	var Restore rds.RestoreDBInstanceFromDBSnapshotInput
	dec := gob.NewDecoder(&b)
	err := dec.Decode(&Restore)
	if err != nil {
		slog.Error("Error decoding state for RestoreDBInstance struct", "error", err)
	}
	return &Restore
}

// GenerateRestoreDBClusterFromSnapshotInput create a snapshot input
func GenerateRestoreDBClusterFromSnapshotInput(r RDSRestorationStore) *rds.RestoreDBClusterFromSnapshotInput {
	return &rds.RestoreDBClusterFromSnapshotInput{
		DBClusterIdentifier: r.GetDBClusterIdentifier(),
		Engine:              r.GetClusterEngine(),
		SnapshotIdentifier:  r.GetClusterSnapshotIdentifier(),
	}
}

// EncodeRestoreDBClusterFromSnapshotInput takes a cluster snapshot and turns it into bytes
func EncodeRestoreDBClusterFromSnapshotInput(r *rds.RestoreDBClusterFromSnapshotInput) bytes.Buffer {
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)

	err := enc.Encode(r)
	if err != nil {
		slog.Error("Error encoding our snapshot", "error", err)
	}
	return encoder
}

// DecodeRestoreDBClusterFromSnapshotInput takes bytes and retruns a db cluster from snapshot input which is needed for restoring a db cluster
func DecodeRestoreDBClusterFromSnapshotInput(b bytes.Buffer) *rds.RestoreDBClusterFromSnapshotInput {
	var Restore rds.RestoreDBClusterFromSnapshotInput
	dec := gob.NewDecoder(&b)
	err := dec.Decode(&Restore)
	if err != nil {
		slog.Error("Error decoding state for RestoreDBCluster struct", "error", err)
	}
	return &Restore
}

func DecodeRDSClusterSnapshotOutput(b bytes.Buffer) types.DBClusterSnapshot {
	var dbSnapshot types.DBClusterSnapshot
	dec := gob.NewDecoder(&b)
	err := dec.Decode(&dbSnapshot)
	if err != nil {
		slog.Error("Error decoding state for cluster snapshot", "error", err)
	}
	return dbSnapshot
}

func GetRDSClusterSnapshotOutput(s StateManager, snap string) (*types.DBClusterSnapshot, error) {
	i := s.GetStateObject(snap)
	snapshot, ok := i.(types.DBClusterSnapshot)
	if !ok {
		str := fmt.Sprintf("error decoding cluster snapshot from interface %v", i)
		return nil, errors.New(str)
	}
	return &snapshot, nil
}

func GetRDSDatabaseInstanceOutput(s StateManager, dbName string) (*types.DBInstance, error) {
	i := s.GetStateObject(dbName)
	dbi, ok := i.(types.DBInstance)
	if !ok {
		str := fmt.Sprintf("error decoding instance from interface %v", i)
		return nil, errors.New(str)
	}
	return &dbi, nil
}

func GetRDSDatabaseClusterOutput(s StateManager, dbName string) (*types.DBCluster, error) {
	i := s.GetStateObject(dbName)
	dbi, ok := i.(types.DBCluster)
	if !ok {
		str := fmt.Sprintf("error decoding cluster from interface %v", i)
		return nil, errors.New(str)
	}
	return &dbi, nil
}

// DecodeRDSSnapshhotOutput takes a bytes buffer and returns it to a DbSnapshot type in preperation of restoring the database
func DecodeRDSSnapshotOutput(b bytes.Buffer) types.DBSnapshot {
	var dbSnapshot types.DBSnapshot
	dec := gob.NewDecoder(&b)
	err := dec.Decode(&dbSnapshot)
	if err != nil {
		slog.Error("Error decoding state for snapshot", "error", err)
	}
	return dbSnapshot
}

func WriteOutput(filename string, b bytes.Buffer) (int64, error) {
	f, err := os.Create(filename)
	if err != nil {
		slog.Error("Error creating file", "error", err)
	}
	defer f.Close()
	n, err := b.WriteTo(f)
	if err != nil {
		slog.Error("error writing to file", "error", err)
	}
	return n, err
}

// CreateInstanceInput creates an instance to prep for creating our Cluster
func CreateDbInstanceInput(i *types.DBInstance, ci *string) *rds.CreateDBInstanceInput {
	dbID := helpers.InstanceName()
	slog.Info("Creating database input")
	sgs := []string{}
	if i.VpcSecurityGroups != nil {
		for _, sg := range i.VpcSecurityGroups {
			sgs = append(sgs, *sg.VpcSecurityGroupId)
		}
	}
	var og *string
	if len(i.OptionGroupMemberships) > 0 {
		og = i.OptionGroupMemberships[0].OptionGroupName
	}
	var pg *string
	if len(i.DBParameterGroups) > 0 {
		pg = i.DBParameterGroups[0].DBParameterGroupName
	}
	var dbSubnetGroup *string
	if i.DBSubnetGroup != nil {
		dbSubnetGroup = i.DBSubnetGroup.DBSubnetGroupName
	}
	var domain *string
	var authSecretArn *string
	var dnsIps []string
	var domainOu *string
	if i.DomainMemberships != nil {
		if len(i.DomainMemberships) > 0 {
			domain = i.DomainMemberships[0].Domain
			authSecretArn = i.DomainMemberships[0].AuthSecretArn
			dnsIps = i.DomainMemberships[0].DnsIps
			domainOu = i.DomainMemberships[0].OU
		}
	}
	var masterUserPassword *bool
	var kmsKeyId *string
	if i.MasterUserSecret != nil {
		mup := true
		masterUserPassword = &mup
		if i.MasterUserSecret.KmsKeyId != nil {
			kmsKeyId = i.MasterUserSecret.KmsKeyId
		}
	}
	slog.Info("Testing allocated storage", "allocated storage", i.AllocatedStorage)
	return &rds.CreateDBInstanceInput{
		DBInstanceClass:                    i.DBInstanceClass,
		DBInstanceIdentifier:               dbID,
		Engine:                             i.Engine,
		AutoMinorVersionUpgrade:            i.AutoMinorVersionUpgrade,
		EngineVersion:                      i.EngineVersion,
		StorageType:                        i.StorageType,
		DBClusterIdentifier:                ci,
		AllocatedStorage:                   i.AllocatedStorage,
		StorageEncrypted:                   i.StorageEncrypted,
		MaxAllocatedStorage:                i.MaxAllocatedStorage,
		MultiAZ:                            i.MultiAZ,
		LicenseModel:                       i.LicenseModel,
		BackupRetentionPeriod:              i.BackupRetentionPeriod,
		PreferredBackupWindow:              i.PreferredBackupWindow,
		PreferredMaintenanceWindow:         i.PreferredMaintenanceWindow,
		PubliclyAccessible:                 i.PubliclyAccessible,
		Iops:                               i.Iops,
		AvailabilityZone:                   i.AvailabilityZone,
		BackupTarget:                       i.BackupTarget,
		MonitoringRoleArn:                  i.MonitoringRoleArn,
		CopyTagsToSnapshot:                 i.CopyTagsToSnapshot,
		VpcSecurityGroupIds:                sgs,
		EnableIAMDatabaseAuthentication:    i.IAMDatabaseAuthenticationEnabled,
		EnablePerformanceInsights:          i.PerformanceInsightsEnabled,
		OptionGroupName:                    og,
		DeletionProtection:                 i.DeletionProtection,
		EnableCloudwatchLogsExports:        i.EnabledCloudwatchLogsExports,
		CACertificateIdentifier:            i.CACertificateIdentifier,
		CharacterSetName:                   i.CharacterSetName,
		CustomIamInstanceProfile:           i.CustomIamInstanceProfile,
		DBName:                             i.DBName,
		DBParameterGroupName:               pg,
		DBSubnetGroupName:                  dbSubnetGroup,
		DBSystemId:                         i.DBSystemId,
		DedicatedLogVolume:                 i.DedicatedLogVolume,
		Domain:                             domain,
		DomainAuthSecretArn:                authSecretArn,
		DomainDnsIps:                       dnsIps,
		DomainOu:                           domainOu,
		EnableCustomerOwnedIp:              i.CustomerOwnedIpEnabled,
		KmsKeyId:                           i.KmsKeyId,
		ManageMasterUserPassword:           masterUserPassword,
		MasterUserSecretKmsKeyId:           kmsKeyId,
		MasterUsername:                     i.MasterUsername,
		MonitoringInterval:                 i.MonitoringInterval,
		MultiTenant:                        i.MultiTenant,
		NcharCharacterSetName:              i.NcharCharacterSetName,
		NetworkType:                        i.NetworkType,
		PerformanceInsightsKMSKeyId:        i.PerformanceInsightsKMSKeyId,
		PerformanceInsightsRetentionPeriod: i.PerformanceInsightsRetentionPeriod,
		Port:                               i.DbInstancePort,
		ProcessorFeatures:                  i.ProcessorFeatures,
		PromotionTier:                      i.PromotionTier,
		StorageThroughput:                  i.StorageThroughput,
		Tags:                               i.TagList,
		TdeCredentialArn:                   i.TdeCredentialArn,
		Timezone:                           i.Timezone,
	}
}

func CreateDBClusterInput(c *types.DBCluster) *rds.CreateDBClusterInput {
	return &rds.CreateDBClusterInput{
		DBClusterIdentifier:     c.DBClusterIdentifier,
		Engine:                  c.Engine,
		AllocatedStorage:        c.AllocatedStorage,
		AutoMinorVersionUpgrade: c.AutoMinorVersionUpgrade,
		AvailabilityZones:       c.AvailabilityZones,
		BacktrackWindow:         c.BacktrackWindow,
		BackupRetentionPeriod:   c.BackupRetentionPeriod,
	}
}

// EncodeCreateDBInstanceInput bytes buffer for create
func EncodeCreateDBInstanceInput(c *rds.CreateDBInstanceInput) bytes.Buffer {
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)
	err := enc.Encode(&c)
	if err != nil {
		slog.Error("Error encoding our database", "error", err)
	}
	return encoder
}

// DecodeCreateDBInstanceInput creates the instance from our bytes buffer when we want to replay
func DecodeCreateDBInstanceInput(b bytes.Buffer) *rds.CreateDBInstanceInput {
	var dbInstance rds.CreateDBInstanceInput
	dec := gob.NewDecoder(&b)
	err := dec.Decode(&dbInstance)
	if err != nil {
		slog.Error("Error decoding state for RDS Instance", "error", err)
	}
	return &dbInstance
}

func EncodeClusterCreateDBInstanceInput(c []rds.CreateDBInstanceInput) bytes.Buffer {
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)
	err := enc.Encode(&c)
	if err != nil {
		slog.Error("Error encoding our database", "error", err)
	}
	return encoder
}

func DecodeClusterCreateDBInstanceInput(b bytes.Buffer) []rds.CreateDBInstanceInput {
	var dbInstances []rds.CreateDBInstanceInput
	dec := gob.NewDecoder(&b)
	err := dec.Decode(&dbInstances)
	if err != nil {
		slog.Error("Error decoding state for RDS Instance", "error", err)
	}
	return dbInstances
}
