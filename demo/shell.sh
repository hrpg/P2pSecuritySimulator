#!/usr/bin/env bash

start=`date +%s`
sleep 5

end=`date +%s`
elapsed=$(($end-$start))

echo "time: ${elapsed} S"