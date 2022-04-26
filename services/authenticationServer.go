package services

import (
	"log"
	"net"
	"net/rpc"
	"net/http"
	"os"
	"sync"

	"P2pSecuritySimulator/cryptoalgs"
)

type AuthenticationServer struct {
	once sync.Once
	mux sync.RWMutex

	userInfos map[string]string

	cryptoMachine cryptoalgs.CryptoMachine

	alive bool
}

func (a *AuthenticationServer) IsAlive() bool {
	return a.alive
}

func (a *AuthenticationServer) checkUserInfo(name, password string) (bool, string) {
	if _, ok := a.userInfos[name]; !ok {
		return false, ErrUserNotExist
	}

	curPassword, _ := a.userInfos[name]
	if curPassword != password {
		return false, ErrPassword
	}

	return true, NoError
}


func (a *AuthenticationServer) server() {
	rpc.Register(a)
	rpc.HandleHTTP()

	serverName := "/var/tmp/server"
	os.Remove(serverName)
	l, e := net.Listen("unix", serverName)
	if e != nil {
		log.Fatal("Server listen error: ", e.Error())
	}

	go http.Serve(l, nil)
}

func (a *AuthenticationServer) Register(req *RegisterReq, rsp *RegisterRsp) {
	a.mux.RLock()
	if _, ok := a.userInfos[req.Name]; ok {
		a.mux.RUnlock()
		rsp.Error = ErrUserExisted
		return
	}

	a.mux.Lock()
	a.userInfos[req.Name] = req.PassWord
	a.mux.Unlock()

	rsp.ServerPubKey = a.cryptoMachine.GetPublicKeyBytes()
	rsp.Error = NoError
}

func (a *AuthenticationServer) AssignCertificate(req *GetCertificateReq, rsp *GetCertificateRsp) {
	flag, err := a.checkUserInfo(req.Name, req.Password)
	if !flag {
		rsp.Error = err
		return
	}

	peerCert := a.cryptoMachine.GenerateCertificate(req.PeerPublicKey)
	//
}

func MakeAuthentication() *AuthenticationServer {
	a := AuthenticationServer{}

	a.once.Do(func() {
		a.userInfos = make(map[string]string)
		a.alive = true

		a.cryptoMachine = &cryptoalgs.Ecc{}
	})

	a.cryptoMachine.GenerateKeys()

	a.server()

	return a
}



