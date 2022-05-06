import csv

csvAuthenTimeFile = open("authentificateTimeSM2WithSM2.csv", "r")
csvRequireCertTimeFile = open("requireCertificateTime.csv", "r")
authenTimeFileReader = csv.reader(csvAuthenTimeFile)
requireCertTimeFileReader = csv.reader(csvRequireCertTimeFile)

sum = 0
cnt = 0
for item in authenTimeFileReader:
    sum += int(item[0])
    cnt += 1

print("Average Authentificate Time:", sum / cnt, "ns")

sum = 0
cnt = 0
for item in requireCertTimeFileReader:
    sum += int(item[0])
    cnt += 1

print("Average Require Certificate Time:", sum / cnt, "ns")

csvAuthenTimeFile.close()
csvRequireCertTimeFile.close()