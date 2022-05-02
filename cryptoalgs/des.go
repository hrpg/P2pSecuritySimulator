package cryptoalgs

import (
	"crypto/des"
)

type Des struct {

}

type CertificateDES struct {
	Text []byte
}

func (d *Des) GenerateKeys() {
	privateKey :=
}