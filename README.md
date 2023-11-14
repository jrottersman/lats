# Lats

## WARNING: Still in rapid development and not suited for production use

## Features
* Snapshot creation for RDS databases
* Copy Snapshots, parameter groups and option groups to another region
* restore snapshot and create new parameter group and option groups for the restored snapshot

## Development plan
Lats is a tool to simplify disaster recovery and multiregion movement in AWS. It's currently in heavy development and is working on getting RDS operations fully going. 
In the near future it will allow
1. Copying DB snapshots between regions
1. Making DB snapshots immutable using AWS Backup
1. Partial Restores
1. Migration of IAM roles and DB Parameter Groups

## Running Lats

to build lats run `go run .`
Add lats to your path
The first time you run lats you will need to run `./lats init` this will prompt you for your aws regions and create a state file entry for running lats. 
The state file entry are json files however do not edit them manually or lats will fail to resotre snapshots

## Lats commands
* lats init 
* lats create-snapshot --database-name {dbName} --snapshot-name {snapshotName}


## Contributing
1. Fork the repository and clone it
1. Have Go and the golang AWS SDK installed 
