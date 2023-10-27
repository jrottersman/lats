go build .
./lats init
./lats CreateRDSSnapshot --database-name test --snapshot-name test
echo "sleep for five minutes to create snapshot TODO change this to actually check status and then go to next step"
sleep 300
./lats CopyRDSSnapshot --snapshot test --new-snapshot test2
echo "sleep for five minutes to copy snapshot TODO change this as above"
sleep 300
./lats restoreRDSSnapshot --snapshot-name test2 --database-name testRestore --region us-east-2
