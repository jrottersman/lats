package state

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/rds/types"
)

type rdsDatabase struct {

	// Indicates whether engine-native audit fields are included in the database
	// activity stream.
	ActivityStreamEngineNativeAuditFieldsIncluded *bool `json:"activityStreamEngineNativeAuditFieldsIncluded"`

	// The name of the Amazon Kinesis data stream used for the database activity
	// stream.
	ActivityStreamKinesisStreamName *string `json:"activityStreamKinesisStreamName"`

	// The Amazon Web Services KMS key identifier used for encrypting messages in the
	// database activity stream. The Amazon Web Services KMS key identifier is the key
	// ARN, key ID, alias ARN, or alias name for the KMS key.
	ActivityStreamKmsKeyId *string `json:"activityStreamKmsKeyId"`

	// The mode of the database activity stream. Database events such as a change or
	// access generate an activity stream event. RDS for Oracle always handles these
	// events asynchronously.
	ActivityStreamMode types.ActivityStreamMode

	// The status of the policy state of the activity stream.
	ActivityStreamPolicyStatus types.ActivityStreamPolicyStatus

	// The status of the database activity stream.
	ActivityStreamStatus types.ActivityStreamStatus

	// The amount of storage in gibibytes (GiB) allocated for the DB instance.
	AllocatedStorage int32 `json:"allocatedStorage"`

	// The Amazon Web Services Identity and Access Management (IAM) roles associated
	// with the DB instance.
	AssociatedRoles []types.DBInstanceRole

	// Indicates whether minor version patches are applied automatically.
	AutoMinorVersionUpgrade bool `json:"autoMinorVersionUpgrade"`

	// The time when a stopped DB instance is restarted automatically.
	AutomaticRestartTime *time.Time `json:"automaticRestartTime"`

	// The automation mode of the RDS Custom DB instance: full or all paused . If full
	// , the DB instance automates monitoring and instance recovery. If all paused ,
	// the instance pauses automation for the duration set by
	// --resume-full-automation-mode-minutes .
	AutomationMode types.AutomationMode

	// The name of the Availability Zone where the DB instance is located.
	AvailabilityZone *string

	// The Amazon Resource Name (ARN) of the recovery point in Amazon Web Services
	// Backup.
	AwsBackupRecoveryPointArn *string

	// The number of days for which automatic DB snapshots are retained.
	BackupRetentionPeriod int32

	// The location where automated backups and manual snapshots are stored: Amazon
	// Web Services Outposts or the Amazon Web Services Region.
	BackupTarget *string

	// The identifier of the CA certificate for this DB instance. For more
	// information, see Using SSL/TLS to encrypt a connection to a DB instance (https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.SSL.html)
	// in the Amazon RDS User Guide and Using SSL/TLS to encrypt a connection to a DB
	// cluster (https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/UsingWithRDS.SSL.html)
	// in the Amazon Aurora User Guide.
	CACertificateIdentifier *string

	// The details of the DB instance's server certificate.
	CertificateDetails *types.CertificateDetails

	// If present, specifies the name of the character set that this instance is
	// associated with.
	CharacterSetName *string

	// Indicates whether tags are copied from the DB instance to snapshots of the DB
	// instance. This setting doesn't apply to Amazon Aurora DB instances. Copying tags
	// to snapshots is managed by the DB cluster. Setting this value for an Aurora DB
	// instance has no effect on the DB cluster setting. For more information, see
	// DBCluster .
	CopyTagsToSnapshot bool

	// The instance profile associated with the underlying Amazon EC2 instance of an
	// RDS Custom DB instance. The instance profile must meet the following
	// requirements:
	//   - The profile must exist in your account.
	//   - The profile must have an IAM role that Amazon EC2 has permissions to
	//   assume.
	//   - The instance profile name and the associated IAM role name must start with
	//   the prefix AWSRDSCustom .
	// For the list of permissions required for the IAM role, see  Configure IAM and
	// your VPC (https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/custom-setup-orcl.html#custom-setup-orcl.iam-vpc)
	// in the Amazon RDS User Guide.
	CustomIamInstanceProfile *string

	// Indicates whether a customer-owned IP address (CoIP) is enabled for an RDS on
	// Outposts DB instance. A CoIP provides local or external connectivity to
	// resources in your Outpost subnets through your on-premises network. For some use
	// cases, a CoIP can provide lower latency for connections to the DB instance from
	// outside of its virtual private cloud (VPC) on your local network. For more
	// information about RDS on Outposts, see Working with Amazon RDS on Amazon Web
	// Services Outposts (https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-on-outposts.html)
	// in the Amazon RDS User Guide. For more information about CoIPs, see
	// Customer-owned IP addresses (https://docs.aws.amazon.com/outposts/latest/userguide/routing.html#ip-addressing)
	// in the Amazon Web Services Outposts User Guide.
	CustomerOwnedIpEnabled *bool

	// If the DB instance is a member of a DB cluster, indicates the name of the DB
	// cluster that the DB instance is a member of.
	DBClusterIdentifier *string

	// The Amazon Resource Name (ARN) for the DB instance.
	DBInstanceArn *string

	// The list of replicated automated backups associated with the DB instance.
	DBInstanceAutomatedBackupsReplications []types.DBInstanceAutomatedBackupsReplication

	// The name of the compute and memory capacity class of the DB instance.
	DBInstanceClass *string

	// The user-supplied database identifier. This identifier is the unique key that
	// identifies a DB instance.
	DBInstanceIdentifier *string

	// The current state of this database. For information about DB instance statuses,
	// see Viewing DB instance status (https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/accessing-monitoring.html#Overview.DBInstance.Status)
	// in the Amazon RDS User Guide.
	DBInstanceStatus *string

	// Contains the initial database name that you provided (if required) when you
	// created the DB instance. This name is returned for the life of your DB instance.
	// For an RDS for Oracle CDB instance, the name identifies the PDB rather than the
	// CDB.
	DBName *string

	// The list of DB parameter groups applied to this DB instance.
	DBParameterGroups []types.DBParameterGroupStatus

	// A list of DB security group elements containing DBSecurityGroup.Name and
	// DBSecurityGroup.Status subelements.
	DBSecurityGroups []types.DBSecurityGroupMembership

	// Information about the subnet group associated with the DB instance, including
	// the name, description, and subnets in the subnet group.
	DBSubnetGroup *types.DBSubnetGroup

	// The Oracle system ID (Oracle SID) for a container database (CDB). The Oracle
	// SID is also the name of the CDB. This setting is only valid for RDS Custom DB
	// instances.
	DBSystemId *string

	// The port that the DB instance listens on. If the DB instance is part of a DB
	// cluster, this can be a different port than the DB cluster port.
	DbInstancePort int32

	// The Amazon Web Services Region-unique, immutable identifier for the DB
	// instance. This identifier is found in Amazon Web Services CloudTrail log entries
	// whenever the Amazon Web Services KMS key for the DB instance is accessed.
	DbiResourceId *string

	// Indicates whether the DB instance has deletion protection enabled. The database
	// can't be deleted when deletion protection is enabled. For more information, see
	// Deleting a DB Instance (https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_DeleteInstance.html)
	// .
	DeletionProtection bool

	// The Active Directory Domain membership records associated with the DB instance.
	DomainMemberships []types.DomainMembership

	// A list of log types that this DB instance is configured to export to CloudWatch
	// Logs. Log types vary by DB engine. For information about the log types for each
	// DB engine, see Monitoring Amazon RDS log files (https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_LogAccess.html)
	// in the Amazon RDS User Guide.
	EnabledCloudwatchLogsExports []string

	// The connection endpoint for the DB instance. The endpoint might not be shown
	// for instances with the status of creating .
	Endpoint *types.Endpoint

	// The database engine used for this DB instance.
	Engine *string

	// The version of the database engine.
	EngineVersion *string

	// The Amazon Resource Name (ARN) of the Amazon CloudWatch Logs log stream that
	// receives the Enhanced Monitoring metrics data for the DB instance.
	EnhancedMonitoringResourceArn *string

	// Indicates whether mapping of Amazon Web Services Identity and Access Management
	// (IAM) accounts to database accounts is enabled for the DB instance. For a list
	// of engine versions that support IAM database authentication, see IAM database
	// authentication (https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Concepts.RDS_Fea_Regions_DB-eng.Feature.IamDatabaseAuthentication.html)
	// in the Amazon RDS User Guide and IAM database authentication in Aurora (https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/Concepts.Aurora_Fea_Regions_DB-eng.Feature.IAMdbauth.html)
	// in the Amazon Aurora User Guide.
	IAMDatabaseAuthenticationEnabled bool

	// The date and time when the DB instance was created.
	InstanceCreateTime *time.Time

	// The Provisioned IOPS (I/O operations per second) value for the DB instance.
	Iops *int32

	// If StorageEncrypted is enabled, the Amazon Web Services KMS key identifier for
	// the encrypted DB instance. The Amazon Web Services KMS key identifier is the key
	// ARN, key ID, alias ARN, or alias name for the KMS key.
	KmsKeyId *string

	// The latest time to which a database in this DB instance can be restored with
	// point-in-time restore.
	LatestRestorableTime *time.Time

	// The license model information for this DB instance. This setting doesn't apply
	// to RDS Custom DB instances.
	LicenseModel *string

	// The listener connection endpoint for SQL Server Always On.
	ListenerEndpoint *types.Endpoint

	// The secret managed by RDS in Amazon Web Services Secrets Manager for the master
	// user password. For more information, see Password management with Amazon Web
	// Services Secrets Manager (https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-secrets-manager.html)
	// in the Amazon RDS User Guide.
	MasterUserSecret *types.MasterUserSecret

	// The master username for the DB instance.
	MasterUsername *string

	// The upper limit in gibibytes (GiB) to which Amazon RDS can automatically scale
	// the storage of the DB instance.
	MaxAllocatedStorage *int32

	// The interval, in seconds, between points when Enhanced Monitoring metrics are
	// collected for the DB instance.
	MonitoringInterval *int32

	// The ARN for the IAM role that permits RDS to send Enhanced Monitoring metrics
	// to Amazon CloudWatch Logs.
	MonitoringRoleArn *string

	// Indicates whether the DB instance is a Multi-AZ deployment. This setting
	// doesn't apply to RDS Custom DB instances.
	MultiAZ bool

	// The name of the NCHAR character set for the Oracle DB instance. This character
	// set specifies the Unicode encoding for data stored in table columns of type
	// NCHAR, NCLOB, or NVARCHAR2.
	NcharCharacterSetName *string

	// The network type of the DB instance. The network type is determined by the
	// DBSubnetGroup specified for the DB instance. A DBSubnetGroup can support only
	// the IPv4 protocol or the IPv4 and the IPv6 protocols ( DUAL ). For more
	// information, see Working with a DB instance in a VPC (https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_VPC.WorkingWithRDSInstanceinaVPC.html)
	// in the Amazon RDS User Guide and Working with a DB instance in a VPC (https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/USER_VPC.WorkingWithRDSInstanceinaVPC.html)
	// in the Amazon Aurora User Guide. Valid Values: IPV4 | DUAL
	NetworkType *string

	// The list of option group memberships for this DB instance.
	OptionGroupMemberships []types.OptionGroupMembership

	// Information about pending changes to the DB instance. This information is
	// returned only when there are pending changes. Specific changes are identified by
	// subelements.
	PendingModifiedValues *types.PendingModifiedValues

	// The progress of the storage optimization operation as a percentage.
	PercentProgress *string

	// Indicates whether Performance Insights is enabled for the DB instance.
	PerformanceInsightsEnabled *bool

	// The Amazon Web Services KMS key identifier for encryption of Performance
	// Insights data. The Amazon Web Services KMS key identifier is the key ARN, key
	// ID, alias ARN, or alias name for the KMS key.
	PerformanceInsightsKMSKeyId *string

	// The number of days to retain Performance Insights data. Valid Values:
	//   - 7
	//   - month * 31, where month is a number of months from 1-23. Examples: 93 (3
	//   months * 31), 341 (11 months * 31), 589 (19 months * 31)
	//   - 731
	// Default: 7 days
	PerformanceInsightsRetentionPeriod *int32

	// The daily time range during which automated backups are created if automated
	// backups are enabled, as determined by the BackupRetentionPeriod .
	PreferredBackupWindow *string

	// The weekly time range during which system maintenance can occur, in Universal
	// Coordinated Time (UTC).
	PreferredMaintenanceWindow *string

	// The number of CPU cores and the number of threads per core for the DB instance
	// class of the DB instance.
	ProcessorFeatures []types.ProcessorFeature

	// The order of priority in which an Aurora Replica is promoted to the primary
	// instance after a failure of the existing primary instance. For more information,
	// see Fault Tolerance for an Aurora DB Cluster (https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/Concepts.AuroraHighAvailability.html#Aurora.Managing.FaultTolerance)
	// in the Amazon Aurora User Guide.
	PromotionTier *int32

	// Indicates whether the DB instance is publicly accessible. When the DB cluster
	// is publicly accessible, its Domain Name System (DNS) endpoint resolves to the
	// private IP address from within the DB cluster's virtual private cloud (VPC). It
	// resolves to the public IP address from outside of the DB cluster's VPC. Access
	// to the DB cluster is ultimately controlled by the security group it uses. That
	// public access isn't permitted if the security group assigned to the DB cluster
	// doesn't permit it. When the DB instance isn't publicly accessible, it is an
	// internal DB instance with a DNS name that resolves to a private IP address. For
	// more information, see CreateDBInstance .
	PubliclyAccessible bool

	// The identifiers of Aurora DB clusters to which the RDS DB instance is
	// replicated as a read replica. For example, when you create an Aurora read
	// replica of an RDS for MySQL DB instance, the Aurora MySQL DB cluster for the
	// Aurora read replica is shown. This output doesn't contain information about
	// cross-Region Aurora read replicas. Currently, each RDS DB instance can have only
	// one Aurora read replica.
	ReadReplicaDBClusterIdentifiers []string

	// The identifiers of the read replicas associated with this DB instance.
	ReadReplicaDBInstanceIdentifiers []string

	// The identifier of the source DB cluster if this DB instance is a read replica.
	ReadReplicaSourceDBClusterIdentifier *string

	// The identifier of the source DB instance if this DB instance is a read replica.
	ReadReplicaSourceDBInstanceIdentifier *string

	// The open mode of an Oracle read replica. The default is open-read-only . For
	// more information, see Working with Oracle Read Replicas for Amazon RDS (https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/oracle-read-replicas.html)
	// in the Amazon RDS User Guide. This attribute is only supported in RDS for
	// Oracle.
	ReplicaMode types.ReplicaMode

	// The number of minutes to pause the automation. When the time period ends, RDS
	// Custom resumes full automation. The minimum value is 60 (default). The maximum
	// value is 1,440.
	ResumeFullAutomationModeTime *time.Time

	// If present, specifies the name of the secondary Availability Zone for a DB
	// instance with multi-AZ support.
	SecondaryAvailabilityZone *string

	// The status of a read replica. If the DB instance isn't a read replica, the
	// value is blank.
	StatusInfos []types.DBInstanceStatusInfo

	// Indicates whether the DB instance is encrypted.
	StorageEncrypted bool

	// The storage throughput for the DB instance. This setting applies only to the gp3
	// storage type.
	StorageThroughput *int32

	// The storage type associated with the DB instance.
	StorageType *string

	// A list of tags. For more information, see Tagging Amazon RDS Resources (https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_Tagging.html)
	// in the Amazon RDS User Guide.
	TagList []types.Tag

	// The ARN from the key store with which the instance is associated for TDE
	// encryption.
	TdeCredentialArn *string

	// The time zone of the DB instance. In most cases, the Timezone element is empty.
	// Timezone content appears only for Microsoft SQL Server DB instances that were
	// created with a time zone specified.
	Timezone *string

	// The list of Amazon EC2 VPC security groups that the DB instance belongs to.
	VpcSecurityGroups []types.VpcSecurityGroupMembership
}
