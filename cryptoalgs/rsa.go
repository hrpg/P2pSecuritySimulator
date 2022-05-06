package cryptoalgs

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/gob"
)
type Rsa struct {
	PrivateKey *rsa.PrivateKey
	PublicKey rsa.PublicKey
	PrivateKeyBytes []byte
	PublicKeyBytes []byte
}

type CertificateBasedRSA struct {
	Text []byte
	Signature []byte
}

func (r *Rsa) GenerateKeys() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 3072)
	if err != nil {
		panic(err)
	}
	publicKey := privateKey.PublicKey
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	publicKeyBytes := x509.MarshalPKCS1PublicKey(&publicKey)


	r.PrivateKey = privateKey
	r.PublicKey = publicKey
	r.PrivateKeyBytes = privateKeyBytes
	r.PublicKeyBytes = publicKeyBytes
}

func (r *Rsa) GenerateCertificate(text []byte) []byte {
	// 计算hash值
	hasher := sha256.New()
	hasher.Write(text)
	hashVal := hasher.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, r.PrivateKey, crypto.SHA256, hashVal)
	if err != nil {
		panic(err)
	}

	newCertificate := CertificateBasedRSA{
		Text: text,
		Signature: signature,
	}

	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	enc.Encode(newCertificate)

	return buffer.Bytes()
}

func (r *Rsa) Encrypt(text []byte) []byte {
	res, err := rsa.EncryptPKCS1v15(rand.Reader,&r.PublicKey, text)
	if err != nil {
		panic(err)
	}

	return res
}

func (r *Rsa) Decrypt(text []byte) []byte {
	res, err := rsa.DecryptPKCS1v15(rand.Reader, r.PrivateKey, text)
	if err != nil {
		panic(err)
	}

	return res
}

func (r *Rsa) EncryptWithPubKey(text []byte, pubKeyBytes []byte) []byte {
	publicKey, err := x509.ParsePKCS1PublicKey(pubKeyBytes)
	if err != nil {
		panic(err)
	}

	res, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, text)
	if err != nil {
		panic(err)
	}

	return res
}

func (r *Rsa) DecryptWithPriKey(text []byte, priKeyBytes []byte) []byte {
	privateKey, err := x509.ParsePKCS1PrivateKey(priKeyBytes)
	if err != nil {
		panic(err)
	}

	res, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, text)
	if err != nil {
		panic(err)
	}

	return res
}

func (r *Rsa) VerifyCertificate(CertificateBytes []byte, pubKeyBytes []byte) bool {
	buffer := bytes.NewBuffer(CertificateBytes)
	enc := gob.NewDecoder(buffer)
	var peerCerificate CertificateBasedRSA
	enc.Decode(&peerCerificate)

	publicKey, err := x509.ParsePKCS1PublicKey(pubKeyBytes)
	if err != nil {
		panic(err)
	}

	// 计算hash值
	hasher := sha256.New()
	hasher.Write(peerCerificate.Text)
	hashVal := hasher.Sum(nil)

	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashVal, peerCerificate.Signature)
	if err != nil {
		return false
	}

	return true
}

func (r *Rsa) GetPublicKeyBytes() []byte {
	return r.PublicKeyBytes
}

func (r *Rsa) GetPrivateKeyBytes() []byte {
	return r.PrivateKeyBytes
}