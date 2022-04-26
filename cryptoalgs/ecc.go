package cryptoalgs

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/gob"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"math/big"
)

type Ecc struct {
	PublicKey   ecdsa.PublicKey
	PublicKeyBytes []byte
	PrivateKey  *ecdsa.PrivateKey
	PrivateKeyBytes []byte
	EciesPublicKey ecies.PublicKey
	EciesPrivateKey ecies.PrivateKey
	MyCertificate Certificate
	MyCertificateBytes []byte
}

type Certificate struct {
	Text []byte
	R []byte
	S []byte
}

func (e *Ecc) GenerateKeys() {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	e.PrivateKey = privateKey
	e.PublicKey = privateKey.PublicKey

	e.EciesPublicKey = ecies.PublicKey{
		X: e.PublicKey.X,
		Y: e.PublicKey.Y,
		Curve: e.PublicKey.Curve,
		Params: ecies.ParamsFromCurve(e.PublicKey.Curve),
	}

	e.EciesPrivateKey = ecies.PrivateKey{
		PublicKey: e.EciesPublicKey,
		D: e.PrivateKey.D,
	}

	// 序列化
	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		panic(err)
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		panic(err)
	}

	e.PrivateKeyBytes = privateKeyBytes
	e.PublicKeyBytes = publicKeyBytes
}

func (e *Ecc) GenerateCertificate(text []byte) []byte {
	hasher := sha256.New()
	hasher.Write(text)
	hashVal := hasher.Sum(nil)

	r, s, err := ecdsa.Sign(rand.Reader, e.PrivateKey, hashVal)
	if err != nil {
		panic(err)
	}

	// 序列化
	rBytes, _ := r.MarshalText()
	sBytes, _ := s.MarshalText()

	newCertificate := Certificate{
		Text: text,
		R: rBytes,
		S: sBytes,
	}

	e.MyCertificate = newCertificate

	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	enc.Encode(newCertificate)

	e.MyCertificateBytes = b.Bytes()

	return e.MyCertificateBytes
}

func (e *Ecc) Encrypt(text []byte) []byte {
	res, err := ecies.Encrypt(rand.Reader, &e.EciesPublicKey, text, nil, nil)
	if err != nil {
		panic(err)
	}

	return res
}

func (e *Ecc) Decrypt(text []byte) []byte {
	res, err := e.EciesPrivateKey.Decrypt(text, nil, nil)
	if err != nil {
		panic(err)
	}

	return res
}

func (e *Ecc) EncryptWithPubKey(text []byte, pubKeyBytes []byte) []byte {
	// 对公钥反序列化
	publicKey, err := x509.ParsePKIXPublicKey(pubKeyBytes)
	if err != nil {
		panic(err)
	}

	ecdsaPublicKey, ok := publicKey.(*ecdsa.PublicKey)
	if ok == false {
		panic("assert publicKey failed")
	}

	eciesPublicKey := ecies.PublicKey{
		X: ecdsaPublicKey.X,
		Y: ecdsaPublicKey.Y,
		Curve: ecdsaPublicKey.Curve,
		Params: ecies.ParamsFromCurve(ecdsaPublicKey.Curve),
	}

	res, err := ecies.Encrypt(rand.Reader, &eciesPublicKey, text, nil, nil)
	if err != nil {
		panic(err)
	}

	return res
}

func (e *Ecc) DecryptWithPriKey(text []byte, priKeyBytes []byte) []byte {
	ecdsaPrivateKey, err := x509.ParseECPrivateKey(priKeyBytes)
	if err != nil {
		panic(err)
	}

	eciesPrivateKey := ecies.PrivateKey{
		PublicKey: ecies.PublicKey{
			X: ecdsaPrivateKey.X,
			Y: ecdsaPrivateKey.Y,
			Curve: ecdsaPrivateKey.Curve,
			Params: ecies.ParamsFromCurve(ecdsaPrivateKey.Curve),
		},
		D: ecdsaPrivateKey.D,
	}

	res, err := eciesPrivateKey.Decrypt(text, nil, nil)
	if err != nil {
		panic(err)
	}

	return res
}

func (e *Ecc) VerifyCertificate(CertificateBytes []byte, pubKeyABytes []byte) bool {
	b := bytes.NewBuffer(CertificateBytes)
	dec := gob.NewDecoder(b)
	var certA Certificate
	dec.Decode(&certA)

	publicKeyA, err := x509.ParsePKIXPublicKey(pubKeyABytes)
	if err != nil {
		panic(err)
	}

	ecdsaPubKeyA, ok := publicKeyA.(*ecdsa.PublicKey)
	if !ok {
		panic("assert pubkey failed")
	}

	hasher := sha256.New()
	hasher.Write(certA.Text)
	hashVal := hasher.Sum(nil)

	var r, s big.Int
	r.UnmarshalText(certA.R)
	s.UnmarshalText(certA.S)

	return ecdsa.Verify(ecdsaPubKeyA, hashVal, &r, &s)
}

func (e *Ecc) GetPublicKeyBytes() []byte {
	return e.PublicKeyBytes
}

func (e *Ecc) GetPrivateKeyBytes() []byte {
	return e.PrivateKeyBytes
}

func (e *Ecc) GetCertificateBytes() []byte {
	return e.MyCertificateBytes
}
