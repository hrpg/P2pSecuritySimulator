package cryptoalgs

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
)

type Sm2 struct {
	PrivateKey *sm2.PrivateKey
	PrivateKeyBytes []byte
	PublicKey *sm2.PublicKey
	PublicKeyBytes []byte
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

	s.PrivateKey = privateKey
	s.PrivateKeyBytes = privateKeyBytes
	s.PublicKey = publicKey
	s.PublicKeyBytes = publicKeyBytes
}

func (s *Sm2) GenerateCertificate(text []byte) []byte {
	hasher := sha256.New()
	hasher.Write(text)
	hashVal := hasher.Sum(nil)

	signature, err := s.PrivateKey.Sign(rand.Reader, hashVal, crypto.SHA256)
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
	res, err := s.PublicKey.EncryptAsn1(text, rand.Reader)
	if err != nil {
		panic(err)
	}

	return res
}

func (s *Sm2) Decrypt(text []byte) []byte {
	res, err := s.PrivateKey.DecryptAsn1(text)
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

	hasher := sha256.New()
	hasher.Write(peerCert.Text)
	hashVal := hasher.Sum(nil)

	return pubKey.Verify(hashVal, peerCert.Signature)
}

func (s *Sm2) GetPublicKeyBytes() []byte {
	return s.PublicKeyBytes
}

func (s *Sm2) GetPrivateKeyBytes() []byte {
	return s.PrivateKeyBytes
}