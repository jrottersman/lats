#AWS

aws operations will probably split into it's own library at some point.

Currently there are three main parts of this library
1. RDS operations which directly operates on instances and clusters
1. rds Parameter groups which are for parameter groups for database configuration 
1. KMS operations which is purely for copying snapshots right now though that will probably change.