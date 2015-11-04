#!/bin/bash

PackageRoot='github.com/griesbacher/Yarbes/'

echo "mode: count" > profile.cov
for dir in $(find `ls` -type d);
do
if ls $dir/*.go &> /dev/null; then
	echo $dir
	go test -v -race -covermode=count -coverprofile=profile.tmp $PackageRoot$dir
	if [ -f profile.tmp ]
    then
        cat profile.tmp | tail -n +2 >> profile.cov
        rm profile.tmp
    fi
fi
done

