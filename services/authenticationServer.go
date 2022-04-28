package services

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"

	"P2pSecuritySimulator/cryptoalgs"
)

type AuthenticationServer struct {
	once sync.Once
	mux sync.RWMutex

	userInfos map[string][]string

	cryptoMachine cryptoalgs.CryptoMachine
	nPeers int
}

func (a *AuthenticationServer) Done() bool {
	a.mux.RLock()
	defer a.mux.RUnlock()

	if a.nPeers <= 0 {
		return true
	} else {
		return false
	}
}

func (a *AuthenticationServer) checkUserInfo(name, password string) (bool, string) {
	a.mux.RLock()
	defer a.mux.RUnlock()

	if _, ok := a.userInfos[name]; !ok {
		return false, ErrUserNotExist
	}

	userInfo, _ := a.userInfos[name]
	if userInfo[0] != password {
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
	a.userInfos[req.Name] = []string{req.PassWord, "No"}
	a.mux.Unlock()

	rsp.ServerPubKeyBytes = a.cryptoMachine.GetPublicKeyBytes()
	rsp.Error = NoError
}

func (a *AuthenticationServer) AssignCertificate(req *GetCertificateReq, rsp *GetCertificateRsp) {
	decryptedPeerInfoBytes := a.cryptoMachine.Decrypt(req.EncryptedPeerInfoBytes)

	// 使用gob进行解码
	b := bytes.NewBuffer(decryptedPeerInfoBytes)
	dec := gob.NewDecoder(b)
	var peerInfo PeerInfo
	dec.Decode(&peerInfo)


	flag, err := a.checkUserInfo(peerInfo.Name, peerInfo.Password)
	if !flag {
		rsp.Error = err
		return
	}

	peerCert := a.cryptoMachine.GenerateCertificate(peerInfo.PeerPublicKeyBytes)

	// 使用私钥对证书进行加密
	encryptedPeerCert := a.cryptoMachine.Encrypt(peerCert)

	rsp.EncryptedPeerCertificateBytes = encryptedPeerCert
	rsp.Error = NoError

	a.mux.Lock()
	a.nPeers -= 1;
	a.mux.Unlock()
}

func MakeAuthentication(nPeers int) *AuthenticationServer {
	a := AuthenticationServer{}

	a.once.Do(func() {
		a.userInfos = make(map[string][]string)
		a.nPeers = nPeers

		a.cryptoMachine = &cryptoalgs.Ecc{}
	})

	a.cryptoMachine.GenerateKeys()

	a.server()

	return &a
}



