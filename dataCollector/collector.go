package dataCollector

import (
	"fmt"
	"os"
)

type collector struct {
	authentificateTimeFileHandler *os.File
	requireCertificateTimeFileHandler *os.File
}

var myCollector collector

func init() {
	authentificateTimeFilePath := "authentificateTime"
	requireCertificateTimeFile := "requireCertificateTime"

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
}

func AppendAuthentificateTime() {

}

func AppendRequireCertificateTime() {

}


