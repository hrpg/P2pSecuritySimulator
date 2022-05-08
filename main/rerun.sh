#!/usr/bin/env bash

go build startServer.go
go build startPeer.go

for i in $(seq 1 200)
do
  echo "第 ${i} 次运行"
  ./test.sh startServer startPeer
done

echo "|****---------------------------------------------****|"
echo "RSAWithDSA: "
python3 count.py
mv authentificateTime.csv authentificateTimeRSAWithDSA.csv
mv requireCertificateTime.csv requireCertificateTimeRSAWithDSA.csv
