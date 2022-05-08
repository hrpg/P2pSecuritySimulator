import csv

csvAuthenTimeFile = open("authentificateTimeECCWithDSA.csv", "r")
csvRequireCertTimeFile = open("requireCertificateTimeECCWithDSA.csv", "r")
authenTimeFileReader = csv.reader(csvAuthenTimeFile)
requireCertTimeFileReader = csv.reader(csvRequireCertTimeFile)

sum = 0
maxTime = 0
minTime = 99999999999999999999999
cnt = 0
for item in authenTimeFileReader:
    sum += int(item[0])
    if int(item[0]) > maxTime:
        maxTime = int(item[0])
    if int(item[0]) < minTime:
        minTime = int(item[0])
    cnt += 1

print("Average Authentificate Time:", sum / cnt, "ns")
print("Max Authentificate Time:", maxTime, "ns")
print("Min Authentificate Time:", minTime, "ns")

sum = 0
maxTime = 0
minTime = 99999999999999999999999
cnt = 0
for item in requireCertTimeFileReader:
    sum += int(item[0])
    if int(item[0]) > maxTime:
        maxTime = int(item[0])
    if int(item[0]) < minTime:
        minTime = int(item[0])
    cnt += 1

print("Average Require Certificate Time:", sum / cnt, "ns")
print("Max Require Certificate Time:", maxTime, "ns")
print("Min Require Certificate Time:", minTime, "ns")

csvAuthenTimeFile.close()
csvRequireCertTimeFile.close()