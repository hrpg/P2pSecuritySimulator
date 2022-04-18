package cryptoalgs

import "crypto"

type Ecc struct {
	PublicKey   string
	PrivateKey  string
	Certificate []byte
}

func (e *Ecc) GenerateKeys() error {

	return nil
}

func (e *Ecc) GenerateCertificate() error {

	return nil
}

func (e *Ecc) Encrypt(PrivateKey string, value string) (string, error) {

}

func (e *Ecc) Decrypt(PublicKey string, value string) (string, error) {

}

func (e *Ecc) GetPublicKey() string {
	return e.PublicKey
}

func (e *Ecc) GetPrivateKey() string {
	return e.PrivateKey
}

func (e *Ecc) GetCertificate() []byte {
	return e.Certificate
}
