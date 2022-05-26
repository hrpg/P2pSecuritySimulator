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
	PeerCryptoPublicKeyBytes []byte
}

type RegisterReq struct {
	Name string
	PassWord string
}

type RegisterRsp struct {
	ServerCryptoPubKeyBytes []byte
	ServerSignPubKeyBytes []byte
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
	PeerACryptoPublicKeyBytes []byte
}

type AuthenticateRsp struct {
	PeerBCertAndCryptoPubKeyInfoBytes []byte
	Error string
}

type PeerBCertAndPubKeyInfo struct {
	PeerBCertificateBytes []byte
	PeerBCryptoPublicKeyBytes  []byte
}

type FinalizeReq struct {
	PeerName string
	Echo string
}

type FinalizeRsp struct {

}

type ReportWorkDoneReq struct {

}

type ReportWorkDoneRsp struct {

}

type CanRestartWorkReq struct {

}

type CanRestartWorkRsp struct {

}



