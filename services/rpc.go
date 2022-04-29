package services

const (
	NoError string = "NoError"
	ErrUserExisted = "UserExisted"
	ErrUserNotExist = "UserNotExist"
	ErrAuthenticationFailed = "AuthenticationFailed"
	ErrPassword = "PasswordError"
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
	PeerACertificateBytes []byte
	PeerAPublicKeyBytes []byte
}

type AuthenticateRsp struct {
	PeerBCertificateBytes []byte
	PeerBPublicKeyBytes  []byte
	Error string
}





