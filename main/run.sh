#!/usr/bin/env bash

for i in $(seq 1 100)
do
  echo "第 ${i} 次运行"
  ./test.sh startServerRSAWithECC startPeerRSAWithECC
  sleep 3
done

echo "|****---------------------------------------------****|"
echo "RSAWithECC: "
python3 count.py
mv authentificateTime.csv authentificateTimeRSAWithECC.csv
mv requireCertificateTime.scv requireCertificateTimeRSAWithECC.scv

for i in $(seq 1 100)
do
  echo "第 ${i} 次运行"
  ./test.sh startServerRSAWithSM2 startPeerRSAWithSM2
  sleep 3
done

echo "|****---------------------------------------------****|"
echo "RSAWithSM2: "
python3 count.py
mv authentificateTime.csv authentificateTimeRSAWithSM2.csv
mv requireCertificateTime.scv requireCertificateTimeRSAWithSM2.scv

#for i in $(seq 1 100)
#do
#  echo "第 ${i} 次运行"
#  ./test.sh startServerRSAWithRSA startPeerRSAWithRSA
#  sleep 3
#done

#echo "|****---------------------------------------------****|"
#echo "RSAWithRSA: "
#python3 count.py
#mv authentificateTime.csv authentificateTimeRSAWithRSA.csv
#mv requireCertificateTime.scv requireCertificateTimeRSAWithRSA.scv

for i in $(seq 1 100)
do
  echo "第 ${i} 次运行"
  ./test.sh startServerRSAWithDSA startPeerRSAWithDSA
  sleep 3
done

echo "|****---------------------------------------------****|"
echo "RSAWithDSA: "
python3 count.py
mv authentificateTime.csv authentificateTimeRSAWithDSA.csv
mv requireCertificateTime.scv requireCertificateTimeRSAWithDSA.scv

for i in $(seq 1 100)
do
  echo "第 ${i} 次运行"
  ./test.sh startServerSM2WithECC startPeerSM2WithECC
  sleep 3
done

echo "|****---------------------------------------------****|"
echo "SM2WithECC: "
python3 count.py
mv authentificateTime.csv authentificateTimeSM2WithECC.csv
mv requireCertificateTime.scv requireCertificateTimeSM2WithECC.scv

for i in $(seq 1 100)
do
  echo "第 ${i} 次运行"
  ./test.sh startServerSM2WithSM2 startPeerSM2WithSM2
  sleep 3
done

echo "|****---------------------------------------------****|"
echo "SM2WithSM2: "
python3 count.py
mv authentificateTime.csv authentificateTimeSM2WithSM2.csv
mv requireCertificateTime.scv requireCertificateTimeSM2WithSM2.scv

for i in $(seq 1 100)
do
  echo "第 ${i} 次运行"
  ./test.sh startServerSM2WithRSA startPeerSM2WithRSA
  sleep 3
done

echo "|****---------------------------------------------****|"
echo "SM2WithRSA: "
python3 count.py
mv authentificateTime.csv authentificateTimeSM2WithRSA.csv
mv requireCertificateTime.scv requireCertificateTimeSM2WithRSA.scv

for i in $(seq 1 100)
do
  echo "第 ${i} 次运行"
  ./test.sh startServerSM2WithDSA startPeerSM2WithDSA
  sleep 3
done

echo "|****---------------------------------------------****|"
echo "SM2WithDSA: "
python3 count.py
mv authentificateTime.csv authentificateTimeSM2WithDSA.csv
mv requireCertificateTime.scv requireCertificateTimeSM2WithDSA.scv

for i in $(seq 1 100)
do
  echo "第 ${i} 次运行"
  ./test.sh startServerECCWithECC startPeerECCWithECC
  sleep 3
done

echo "|****---------------------------------------------****|"
echo "ECCWithECC: "
python3 count.py
mv authentificateTime.csv authentificateTimeECCWithECC.csv
mv requireCertificateTime.scv requireCertificateTimeECCWithECC.scv

for i in $(seq 1 100)
do
  echo "第 ${i} 次运行"
  ./test.sh startServerECCWithSM2 startPeerECCWithSM2
  sleep 3
done

echo "|****---------------------------------------------****|"
echo "ECCWithSM2: "
python3 count.py
mv authentificateTime.csv authentificateTimeECCWithSM2.csv
mv requireCertificateTime.scv requireCertificateTimeECCWithSM2.scv

for i in $(seq 1 100)
do
  echo "第 ${i} 次运行"
  ./test.sh startServerECCWithRSA startPeerECCWithRSA
  sleep 3
done

echo "|****---------------------------------------------****|"
echo "ECCWithRSA: "
python3 count.py
mv authentificateTime.csv authentificateTimeECCWithRSA.csv
mv requireCertificateTime.scv requireCertificateTimeECCWithRSA.scv

for i in $(seq 1 100)
do
  echo "第 ${i} 次运行"
  ./test.sh startServerECCWithDSA startPeerECCWithDSA
  sleep 3
done

echo "|****---------------------------------------------****|"
echo "ECCWithDSA: "
python3 count.py
mv authentificateTime.csv authentificateTimeECCWithDSA.csv
mv requireCertificateTime.scv requireCertificateTimeECCWithDSA.scv

