#!/bin/sh
ID=`ps -ef | grep "elasticsearch" | grep -v "grep" | awk '{print $2}'`
echo $ID
echo "---------------"
for id in $ID
do
kill -9 $id
echo "killed $id"
done
echo "---------------"
ulimit -n 65536
ulimit -u 65536
sysctl -p
cd /data/elasticsearch/bin
su elastic
./elasticsearch -d