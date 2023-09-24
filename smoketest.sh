go build .
./lats init
./lats CreateRDSSnapshot --database-name test --snapshot-name test
./lats CopyRDSSnapshot --snapshot test --new-snapshot test2
