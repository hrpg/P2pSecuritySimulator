package cryptoalgs

import (
	"bytes"

	"crypto/dsa"
	"crypto/rand"

	"encoding/gob"
	"math/big"
)

type Dsa struct {
	privateKey dsa.PrivateKey
	publicKey dsa.PublicKey
	publicKeyBytes []byte
}

type CertificateBasedDSA struct {
	Text []byte
	RText []byte
	SText []byte
}

func (d *Dsa) GenerateKeys() {
	var param dsa.Parameters
	dsa.GenerateParameters(&param, rand.Reader, dsa.L3072N256)
	var privateKey dsa.PrivateKey
	privateKey.Parameters = param
	dsa.GenerateKey(&privateKey, rand.Reader)
	publicKey := privateKey.PublicKey

	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	enc.Encode(publicKey)

	d.privateKey = privateKey
	d.publicKeyBytes = buffer.Bytes()
}

func (d *Dsa) GenerateCertificate(text []byte) []byte {
	r, s, err := dsa.Sign(rand.Reader, &d.privateKey, text)
	if err != nil {
		panic(err)
	}

	var peerCert CertificateBasedDSA
	rText, err := r.MarshalText()
	if err != nil {
		panic(err)
	}

	sText, err := s.MarshalText()
	if err != nil {
		panic(err)
	}

	peerCert.Text = text
	peerCert.RText = rText
	peerCert.SText = sText

	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	enc.Encode(peerCert)

	return buffer.Bytes()
}

func (d *Dsa) Encrypt(text []byte) []byte {
	return nil
}

func (d *Dsa) Decrypt(text []byte) []byte {
	return nil
}

func (d *Dsa) EncryptWithPubKey(text []byte, pubKeyBytes []byte) []byte {
	return nil
}

func (d *Dsa) DecryptWithPriKey(text []byte, priKeyBytes []byte) []byte {
	return nil
}

func (d *Dsa) VerifyCertificate(CertificateBytes []byte, pubKeyBytes []byte) bool {
	buffer := bytes.NewBuffer(CertificateBytes)
	dec := gob.NewDecoder(buffer)
	var peerCert CertificateBasedDSA
	dec.Decode(&peerCert)

	var r, s big.Int
	r.UnmarshalText(peerCert.RText)
	s.UnmarshalText(peerCert.SText)

	buffer = bytes.NewBuffer(pubKeyBytes)
	dec = gob.NewDecoder(buffer)
	var publicKey dsa.PublicKey
	dec.Decode(&publicKey)

	res := dsa.Verify(&publicKey, peerCert.Text, &r, &s)

	return res
}

func (d *Dsa) GetPublicKeyBytes() []byte {
	return d.publicKeyBytes
}

func (d *Dsa) GetPrivateKeyBytes() []byte {
	return nil
}