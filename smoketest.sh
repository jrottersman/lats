go build .
./lats init
./lats CreateRDSSnapshot --database-name test --snapshot-name test
while :
do 
    snapStatus=`aws rds describe-db-snapshots --db-snapshot-identifier test --output json | jq ".DBSnapshots[] | .Status"`
    if [$snapStatus = "available"]
    then    
        echo "Snapshot complete"
        break
    fi
    echo "snapshot in progress"
    sleep 15
done
./lats CopyRDSSnapshot --snapshot test --new-snapshot test2
while :
do 
    snapStatus=`aws rds describe-db-snapshots --db-snapshot-identifier test2 --output json | jq ".DBSnapshots[] | .Status"`
    if [$snapStatus = "available"]
    then    
        echo "Snapshot complete"
        break
    fi
    echo "snapshot in progress"
    sleep 15
done
./lats restoreRDSSnapshot --snapshot-name test2 --database-name testRestore --region us-east-2 --subnet-group sg-name
