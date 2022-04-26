package cryptoalgs

type CryptoMachine interface {
	GenerateKeys()
	GenerateCertificate(text []byte) []byte
	Encrypt(text []byte) []byte
	Decrypt(text []byte) []byte
	EncryptWithPubKey(text []byte, pubKeyBytes []byte) []byte
	DecryptWithPriKey(text []byte, priKeyBytes []byte) []byte
	VerifyCertificate(CertificateBytes []byte, pubKeyABytes []byte) bool
	GetPublicKeyBytes() []byte
	GetPrivateKeyBytes() []byte
	GetCertificateBytes() []byte
}

