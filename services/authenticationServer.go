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
	"time"

	"P2pSecuritySimulator/cryptoalgs"
)

type AuthenticationServer2 struct {
	once sync.Once
	mux sync.RWMutex

	countVal int
	userInfos map[string]string

	cryptoMachine cryptoalgs.CryptoMachine
	signMachine cryptoalgs.CryptoMachine

	//nPeers int
}

func (a *AuthenticationServer2) Done() bool {
	//a.mux.RLock()
	//defer a.mux.RUnlock()
	//
	//if a.nPeers <= 0 {
	//	return true
	//} else {
	//	return false
	//}

	return false
}

func (a *AuthenticationServer2) checkUserInfo(name, password string) (bool, string) {
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

func call3(machime string, rpcname string, req interface{}, rsp interface{}) {
	c, err := rpc.DialHTTP("unix", machime)
	defer c.Close()
	if err != nil {
		log.Fatal("SERVER: dialing error: ", err.Error())
		return
	}

	err = c.Call(rpcname, req, rsp)
	if err != nil {
		log.Fatalf("SERVER: call %s method %s error: ", err.Error())
	}
}


func (a *AuthenticationServer2) server() {
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

func (a *AuthenticationServer2) Register(req *RegisterReq, rsp *RegisterRsp) error {
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

	rsp.ServerCryptoPubKeyBytes = a.cryptoMachine.GetPublicKeyBytes()
	rsp.ServerSignPubKeyBytes = a.signMachine.GetPublicKeyBytes()
	rsp.Error = NoError

	log.Printf("SERVER: user %s successfully registered", req.Name)

	return nil
}

func (a *AuthenticationServer2) ReportDone(req *ReportWorkDoneReq, rsp *ReportWorkDoneRsp) error {
	a.mux.Lock()
	a.countVal += 1;
	a.mux.Unlock()

	return nil
}

func (a *AuthenticationServer2) checkCountVal() {
	for true {
		a.mux.Lock()
		if a.countVal >= 3 {
			req := CanRestartWorkReq{}
			rsp := CanRestartWorkRsp{}
			for peerName, _ := range a.userInfos {
				call3(peerName, "Peer2.CanRestartWork", &req, &rsp)
			}

			a.countVal = 0
		}

		a.mux.Unlock()

		time.Sleep(10 * time.Millisecond)
	}
}

func (a *AuthenticationServer2) AssignCertificate(req *GetCertificateReq, rsp *GetCertificateRsp) error {
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
	peerCert := a.signMachine.GenerateCertificate(peerInfo.PeerCryptoPublicKeyBytes)

	// 使用用户公钥对证书进行加密
	encryptedPeerCert := a.cryptoMachine.EncryptWithPubKey(peerCert, peerInfo.PeerCryptoPublicKeyBytes)

	rsp.EncryptedPeerCertificateBytes = encryptedPeerCert
	rsp.Error = NoError

	//a.mux.Lock()
	//a.nPeers -= 1;
	//a.mux.Unlock()

	log.Printf("SERVER: user %s successfully get certificate", peerInfo.Name)
	return nil
}

func MakeAuthenticationServer2() *AuthenticationServer2 {
	a := AuthenticationServer2{}

	a.once.Do(func() {
		a.countVal = 0
		a.userInfos = make(map[string]string)
		//a.nPeers = nPeers

		a.cryptoMachine = &cryptoalgs.Ecc{}
		a.signMachine = &cryptoalgs.Dsa{}
	})

	a.cryptoMachine.GenerateKeys()
	a.signMachine.GenerateKeys()
	log.Printf("SERVER: server has generated Keys")

	a.server()
	log.Printf("SERVER: server start to running")

	go a.checkCountVal()

	return &a
}



