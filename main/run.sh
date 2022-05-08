#!/usr/bin/env bash

#for i in $(seq 1 100)
#do
#  echo "第 ${i} 次运行"
#  ./test.sh startServerRSAWithECC startPeerRSAWithECC
#  sleep 3
#done
#
#echo "|****---------------------------------------------****|"
#echo "RSAWithECC: "
#python3 count.py
#mv authentificateTime.csv authentificateTimeRSAWithECC.csv
#mv requireCertificateTime.csv requireCertificateTimeRSAWithECC.csv
#
#for i in $(seq 1 100)
#do
#  echo "第 ${i} 次运行"
#  ./test.sh startServerRSAWithSM2 startPeerRSAWithSM2
#  sleep 3
#done
#
#echo "|****---------------------------------------------****|"
#echo "RSAWithSM2: "
#python3 count.py
#mv authentificateTime.csv authentificateTimeRSAWithSM2.csv
#mv requireCertificateTime.csv requireCertificateTimeRSAWithSM2.csv

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
#mv requireCertificateTime.csv requireCertificateTimeRSAWithRSA.csv

#for i in $(seq 1 100)
#do
#  echo "第 ${i} 次运行"
#  ./test.sh startServerRSAWithDSA startPeerRSAWithDSA
#  sleep 3
#done
#
#echo "|****---------------------------------------------****|"
#echo "RSAWithDSA: "
#python3 count.py
#mv authentificateTime.csv authentificateTimeRSAWithDSA.csv
#mv requireCertificateTime.csv requireCertificateTimeRSAWithDSA.csv

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
mv requireCertificateTime.csv requireCertificateTimeSM2WithECC.csv

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
mv requireCertificateTime.csv requireCertificateTimeSM2WithSM2.csv

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
mv requireCertificateTime.csv requireCertificateTimeSM2WithRSA.csv

#for i in $(seq 1 100)
#do
#  echo "第 ${i} 次运行"
#  ./test.sh startServerSM2WithDSA startPeerSM2WithDSA
#  sleep 3
#done
#
#echo "|****---------------------------------------------****|"
#echo "SM2WithDSA: "
#python3 count.py
#mv authentificateTime.csv authentificateTimeSM2WithDSA.csv
#mv requireCertificateTime.csv requireCertificateTimeSM2WithDSA.csv

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
mv requireCertificateTime.csv requireCertificateTimeECCWithECC.csv

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
mv requireCertificateTime.csv requireCertificateTimeECCWithSM2.csv

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
mv requireCertificateTime.csv requireCertificateTimeECCWithRSA.csv

#for i in $(seq 1 100)
#do
#  echo "第 ${i} 次运行"
#  ./test.sh startServerECCWithDSA startPeerECCWithDSA
#  sleep 3
#done
#
#echo "|****---------------------------------------------****|"
#echo "ECCWithDSA: "
#python3 count.py
#mv authentificateTime.csv authentificateTimeECCWithDSA.csv
#mv requireCertificateTime.csv requireCertificateTimeECCWithDSA.csv

