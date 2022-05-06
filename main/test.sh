#!/usr/bin/env bash

echo "start experiment"
start=`date +%s`

commonProfix="/var/tmp/peer-"

peerAddresses=()

# shellcheck disable=SC1073
for i in $(seq 0 10)
do
  peerAddress=${commonProfix}"$i"
  peerAddresses[${i}]=${peerAddress}
done

./$1 11 &
for i in $(seq 0 10)
do {
      ./$2 ${peerAddresses[${i}]} ${peerAddresses[*]}
   } &
done

wait

end=`date +%s`
elapsed=$(($end - $start))

echo "time: ${elapsed} S"