#!/bin/bash

DATE=`date +%Y%m%dT%H%M%S`
mkdir $DATE

mkdir $DATE/graphdb

curl -X GET -H "Accept:application/n-quads" "http://localhost:7200/repositories/SYSTEM/statements?infer=false" > "$DATE/graphdb/SYSTEM.nq"

for REPO in `grep -oP '(?<=#repositoryID> ")[^"]+' $DATE/graphdb/SYSTEM.nq`; do
	echo "dumping $REPO ..."
	curl -X GET -H "Accept:application/n-quads" "http://localhost:7200/repositories/$REPO/statements?infer=false" > "$DATE/graphdb/${REPO}.nq"
done

tar cfz "$DATE.tgz" "$DATE/"
rm -r "$DATE"