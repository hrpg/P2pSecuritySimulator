package dataCollector

import (
	"fmt"
	"os"
	"strconv"
)

type collector struct {
	authentificateTimeFileHandler *os.File
	requireCertificateTimeFileHandler *os.File
}

var myCollector collector

func init() {
	authentificateTimeFilePath := "authentificateTime.csv"
	requireCertificateTimeFile := "requireCertificateTime.csv"

	fileHandler, err := os.OpenFile(authentificateTimeFilePath, os.O_CREATE | os.O_RDWR | os.O_APPEND, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create authentificateTime file failed")
		os.Exit(1)
	}
	myCollector.authentificateTimeFileHandler = fileHandler

	fileHandler, err = os.OpenFile(requireCertificateTimeFile, os.O_CREATE | os.O_RDWR | os.O_APPEND, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create requireCertificateTime file failed")
		os.Exit(1)
	}
	myCollector.requireCertificateTimeFileHandler = fileHandler
}

func AppendAuthentificateTime(elapsed int64) {
	myCollector.authentificateTimeFileHandler.WriteString(strconv.FormatInt(elapsed, 10) + "\n")
}

func AppendRequireCertificateTime(elapsed int64) {
	myCollector.requireCertificateTimeFileHandler.WriteString(strconv.FormatInt(elapsed, 10) + "\n")
}


