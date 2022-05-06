#!/usr/bin/env bash

for i in $(seq 1 100)
do
  echo "第 ${i} 次运行"
  ./test.sh
  sleep 3
done