package services

const (
	NoError string = "NoError"
	ErrUserExisted = "UserExisted"
	ErrUserNotExist = "UserNotExist"
	ErrAuthenticationFailed = "AuthenticationFailed"
	ErrPassword = "PasswordError"
	ErrBCertInfo = "ErrBCertInfoError"

)

type PeerInfo struct {
	Name string
	Password string
	PeerPublicKeyBytes []byte
}

type RegisterReq struct {
	Name string
	PassWord string
}

type RegisterRsp struct {
	ServerPubKeyBytes []byte
	Error string
}

type GetCertificateReq struct {
	EncryptedPeerInfoBytes []byte
}

type GetCertificateRsp struct {
	EncryptedPeerCertificateBytes []byte
	Error string
}

type AuthenticateReq struct {
	PeerName string
	PeerACertificateBytes []byte
	PeerAPublicKeyBytes []byte
}

type AuthenticateRsp struct {
	PeerBCertAndPubKeyInfoBytes []byte
	Error string
}

type PeerBCertAndPubKeyInfo struct {
	PeerBCertificateBytes []byte
	PeerBPublicKeyBytes  []byte
}

type FinalizeReq struct {
	PeerName string
	Echo string
}

type FinalizeRsp struct {

}




