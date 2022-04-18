package services

const (
	NoError string = "NoError"
	ErrUserExisted = "UserExisted"
	ErrUserNotExist = "UserNotExist"
	ErrPassword = "PasswordError"
)

type RegisterReq struct {
	Name string
	PassWord string
}

type RegisterRsp struct {
	AuthenticationServerKey string
	Error string
}

type GetCertificateReq struct {
	UserInfo []byte
}

type GetCertificateRsp struct {
	Certificate []byte
	Error string
}

type ConnectReq struct {
	CertificateA []byte
	PublicKeyA string
}

type ConnectRsp struct {
	CertificateB []byte
	PublicKeyB  string
}





