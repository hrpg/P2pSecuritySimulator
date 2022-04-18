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
	privateKey string
	publicKey string

	alive bool
}

func (a *AuthenticationServer) IsAlive() bool {
	return a.alive
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

	rsp.AuthenticationServerKey = a.publicKey
	rsp.Error = NoError
}

func MakeAuthentication() *AuthenticationServer {
	a := AuthenticationServer{}

	a.once.Do(func() {
		a.userInfos = make(map[string]string)
		a.alive = true

		a.cryptoMachine = &cryptoalgs.Ecc{}
	})

	err := a.cryptoMachine.GenerateKeys()
	if err != nil {
		log.Fatal("Server failed to generate keys, error: ", err.Error())
		return nil
	}

	a.privateKey = a.cryptoMachine.GetPrivateKey()
	a.publicKey = a.cryptoMachine.GetPublicKey()

	a.server()

	return a
}



