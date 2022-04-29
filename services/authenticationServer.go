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

	userInfos map[string]string

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

	userPassword, _ := a.userInfos[name]
	if userPassword != password {
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
		log.Fatal("SERVER: server listen error: ", e.Error())
	}

	go http.Serve(l, nil)
}

func (a *AuthenticationServer) Register(req *RegisterReq, rsp *RegisterRsp) error {
	log.Printf("SERVER: user %s request registration", req.Name)
	a.mux.Lock()
	log.Printf("SERVER: user %s request registration, get lock", req.Name)
	defer log.Printf("SERVER: user %s request registration, release lock", req.Name)
	defer a.mux.Unlock()
	if _, ok := a.userInfos[req.Name]; ok {
		rsp.Error = ErrUserExisted
		return nil
	}

	a.userInfos[req.Name] = req.PassWord

	rsp.ServerPubKeyBytes = a.cryptoMachine.GetPublicKeyBytes()
	rsp.Error = NoError

	log.Printf("SERVER: user %s successfully registered", req.Name)
	return nil
}

func (a *AuthenticationServer) AssignCertificate(req *GetCertificateReq, rsp *GetCertificateRsp) error {
	decryptedPeerInfoBytes := a.cryptoMachine.Decrypt(req.EncryptedPeerInfoBytes)

	// 使用gob进行解码
	b := bytes.NewBuffer(decryptedPeerInfoBytes)
	dec := gob.NewDecoder(b)
	var peerInfo PeerInfo
	dec.Decode(&peerInfo)

	log.Printf("SERVER: user %s request certificate", peerInfo.Name)
	flag, err := a.checkUserInfo(peerInfo.Name, peerInfo.Password)
	if !flag {
		rsp.Error = err
		return nil
	}

	log.Printf("SERVER: server start to assign certificate to user %s", peerInfo.Name)
	peerCert := a.cryptoMachine.GenerateCertificate(peerInfo.PeerPublicKeyBytes)

	// 使用私钥对证书进行加密
	encryptedPeerCert := a.cryptoMachine.EncryptWithPubKey(peerCert, peerInfo.PeerPublicKeyBytes)

	rsp.EncryptedPeerCertificateBytes = encryptedPeerCert
	rsp.Error = NoError

	a.mux.Lock()
	a.nPeers -= 1;
	a.mux.Unlock()

	log.Printf("SERVER: user %s successfully get certificate", peerInfo.Name)
	return nil
}

func MakeAuthentication(nPeers int) *AuthenticationServer {
	a := AuthenticationServer{}

	a.once.Do(func() {
		a.userInfos = make(map[string]string)
		a.nPeers = nPeers

		a.cryptoMachine = &cryptoalgs.Ecc{}
	})

	a.cryptoMachine.GenerateKeys()
	log.Printf("SERVER: server has generated Keys")

	a.server()
	log.Printf("SERVER: server start to running")

	return &a
}



