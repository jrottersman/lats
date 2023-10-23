go build .
./lats init
./lats CreateRDSSnapshot --database-name test --snapshot-name test
./lats CopyRDSSnapshot --snapshot test --new-snapshot test2
/lats restoreRDSSnapshot --snapshot-name test2 --database-name testRestore --region us-east-2
