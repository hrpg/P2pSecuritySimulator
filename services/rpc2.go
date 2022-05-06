package services

const (
	NoError2 string = "NoError"
	ErrUserExisted2 = "UserExisted"
	ErrUserNotExist2 = "UserNotExist"
	ErrAuthenticationFailed2 = "AuthenticationFailed"
	ErrPassword2 = "PasswordError"
	ErrBCertInfo2 = "ErrBCertInfoError"

)

type PeerInfo2 struct {
	Name string
	Password string
	PeerCryptoPublicKeyBytes []byte
}

type RegisterReq2 struct {
	Name string
	PassWord string
}

type RegisterRsp2 struct {
	ServerCryptoPubKeyBytes []byte
	ServerSignPubKeyBytes []byte
	Error string
}

type GetCertificateReq2 struct {
	EncryptedPeerInfoBytes []byte
}

type GetCertificateRsp2 struct {
	EncryptedPeerCertificateBytes []byte
	Error string
}

type AuthenticateReq2 struct {
	PeerName string
	PeerACertificateBytes []byte
	PeerACryptoPublicKeyBytes []byte
}

type AuthenticateRsp2 struct {
	PeerBCertAndCryptoPubKeyInfoBytes []byte
	Error string
}

type PeerBCertAndPubKeyInfo2 struct {
	PeerBCertificateBytes []byte
	PeerBCryptoPublicKeyBytes  []byte
}

type FinalizeReq2 struct {
	PeerName string
	Echo string
}

type FinalizeRsp2 struct {

}




