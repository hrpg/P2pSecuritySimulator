package cryptoalgs

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"encoding/gob"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
)

type Sm2 struct {
	privateKey *sm2.PrivateKey
	privateKeyBytes []byte
	publicKey *sm2.PublicKey
	publicKeyBytes []byte
}

type CertificateSM2 struct {
	Text []byte
	Signature []byte
}

func (s *Sm2) GenerateKeys() {
	privateKey, err := sm2.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	publicKey := &privateKey.PublicKey

	privateKeyBytes, err := x509.MarshalSm2UnecryptedPrivateKey(privateKey)
	if err != nil {
		panic(err)
	}

	publicKeyBytes, err := x509.MarshalSm2PublicKey(publicKey)
	if err != nil {
		panic(err)
	}

	s.privateKey = privateKey
	s.privateKeyBytes = privateKeyBytes
	s.publicKey = publicKey
	s.publicKeyBytes = publicKeyBytes
}

func (s *Sm2) GenerateCertificate(text []byte) []byte {
	signature, err := s.privateKey.Sign(rand.Reader, text, crypto.SHA256)
	if err != nil {
		panic(err)
	}

	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	var peerCert CertificateSM2
	peerCert.Text = text
	peerCert.Signature = signature
	enc.Encode(peerCert)

	return buffer.Bytes()
}

func (s *Sm2) Encrypt(text []byte) []byte {
	res, err := s.publicKey.EncryptAsn1(text, rand.Reader)
	if err != nil {
		panic(err)
	}

	return res
}

func (s *Sm2) Decrypt(text []byte) []byte {
	res, err := s.privateKey.DecryptAsn1(text)
	if err != nil {
		panic(err)
	}

	return res
}

func (s *Sm2) EncryptWithPubKey(text []byte, pubKeyBytes []byte) []byte {
	pubKey, err := x509.ParseSm2PublicKey(pubKeyBytes)
	if err != nil {
		panic(err)
	}

	res, err := pubKey.EncryptAsn1(text, rand.Reader)
	if err != nil {
		panic(err)
	}

	return res
}

func (s *Sm2) DecryptWithPriKey(text []byte, priKeyBytes []byte) []byte {
	priKey, err := x509.ParseSm2PrivateKey(priKeyBytes)
	if err != nil {
		panic(priKey)
	}

	res, err := priKey.DecryptAsn1(text)
	if err != nil {
		panic(err)
	}

	return res
}

func (s *Sm2) VerifyCertificate(CertificateBytes []byte, pubKeyABytes []byte) bool {
	buffer := bytes.NewBuffer(CertificateBytes)
	dec := gob.NewDecoder(buffer)
	var peerCert CertificateSM2
	dec.Decode(&peerCert)

	pubKey, err := x509.ParseSm2PublicKey(pubKeyABytes)
	if err != nil {
		panic(err)
	}

	return pubKey.Verify(peerCert.Text, peerCert.Signature)
}

func (s *Sm2) GetPublicKeyBytes() []byte {
	return s.publicKeyBytes
}

func (s *Sm2) GetPrivateKeyBytes() []byte {
	return s.privateKeyBytes
}