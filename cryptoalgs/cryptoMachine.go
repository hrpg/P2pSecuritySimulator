package cryptoalgs

type CryptoMachine interface {
	GenerateKeys() error
	GenerateCertificate() error
	Encrypt(PrivateKey string, value string) (string, error)
	Decrypt(PublicKey string, value string) (string, error)
	GetPublicKey() string
	GetPrivateKey() string
	GetCertificate() []byte
}

