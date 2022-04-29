#!/usr/bin/env bash

echo "start experiment"
start=`date +%s`

go build startServer.go

go build startPeer.go

commonProfix="/var/tmp/peer-"

peerAddresses=()

# shellcheck disable=SC1073
for i in $(seq 0 255)
do
  peerAddress=${commonProfix}"$i"
  peerAddresses[${i}]=${peerAddress}
done

./startServer 1000 &
for i in $(seq 0 999)
do {
      ./startPeer ${peerAddresses[${i}]} ${peerAddresses[*]}
   } &
done

wait
end=`date +%s`
elapsed=$(($end - $start))

echo "time: ${elapsed} S"