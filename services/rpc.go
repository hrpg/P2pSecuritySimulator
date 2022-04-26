package services

const (
	NoError string = "NoError"
	ErrUserExisted = "UserExisted"
	ErrUserNotExist = "UserNotExist"
	ErrPassword = "PasswordError"
)

type PeerInfo struct {
	Name string
	Password string
	PeerPublicKey []byte
}

type RegisterReq struct {
	Name string
	PassWord string
}

type RegisterRsp struct {
	ServerPubKey []byte
	Error string
}

type GetCertificateReq struct {
	PeerInfoBytes []byte
}

type GetCertificateRsp struct {
	DecryptedPeerCertificate []byte
	Error string
}

type AuthenticateReq struct {
	PeerACertificate []byte
	PeerAPublicKey []byte
}

type AuthenticateRsp struct {
	PeerBCertificate []byte
	PeerBPublicKey  []byte
	Error string
}





