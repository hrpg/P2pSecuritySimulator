package services

import (
	"P2pSecuritySimulator/cryptoalgs"
	"bytes"
	"encoding/gob"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
	"time"
)

type Peer struct {
	mux sync.Mutex
	once sync.Once

	name string
	password string

	cryptoMachine cryptoalgs.CryptoMachine
	serverPubKeyBytes []byte
	myCertificate []byte
}

func (p *Peer) Authenticate(req *AuthenticateReq, rsp *AuthenticateRsp) {
	flag := p.cryptoMachine.VerifyCertificate(req.PeerACertificateBytes, req.PeerAPublicKeyBytes)
	if !flag {

	}

	rsp.PeerBCertificateBytes = p.cryptoMachine.GetCertificateBytes()
	rsp.PeerBPublicKeyBytes = p.cryptoMachine.GetPublicKeyBytes()
	rsp.Error = NoError
}

func (p *Peer) register() {
	for i := 0; i < 3; i++  {
		req := RegisterReq{}
		rsp := RegisterRsp{}

		req.Name, req.PassWord = p.name, p.password
		call("/var/tmp/server", "AuthenticationServer.Register", &req, &rsp)
		if rsp.Error == NoError {
			break
		}

		log.Printf("peer %s failed to register, error: %s", p.name, rsp.Error)
		rand.Seed(time.Now().UnixNano())

		sleepTime := 200 + rand.Intn(200)
		time.Sleep(time.Duration(sleepTime) * time.Millisecond)
	}
}

func (p *Peer) requireCertification() {
	for i := 0; i < 3; i++ {
		req := GetCertificateReq{}
		rsp := GetCertificateRsp{}
		var peerInfo PeerInfo
		peerInfo.Name = p.name
		peerInfo.Password = p.password
		peerInfo.PeerPublicKeyBytes = p.cryptoMachine.GetPublicKeyBytes()

		var buffer bytes.Buffer
		enc := gob.NewEncoder(&buffer)
		enc.Encode(peerInfo)
		peerInfoBytes := buffer.Bytes()

		encryptedPeerInfoBytes := p.cryptoMachine.EncryptWithPubKey(peerInfoBytes, p.serverPubKeyBytes)
		req.EncryptedPeerInfoBytes = encryptedPeerInfoBytes

		call("/var/tmp/server", "Authentication.AssignCertification", &req, &rsp)
		if rsp.Error == NoError {
			p.myCertificate = p.cryptoMachine.Decrypt(rsp.EncryptedPeerCertificateBytes)
			break
		}

		log.Printf("peer %s failed to ask certificate, error: %s", p.name, rsp.Error)
		rand.Seed(time.Now().UnixNano())

		sleepTime := 200 + rand.Intn(200)
		time.Sleep(time.Duration(sleepTime) * time.Millisecond)
	}
}

func (p *Peer) RequireAuthentication(peeraddr string) {
	req := AuthenticateReq{}
	rsp := AuthenticateRsp{}
	req.PeerACertificateBytes = p.myCertificate
	req.PeerAPublicKeyBytes = p.cryptoMachine.GetPublicKeyBytes()

	call(peeraddr, "Authentication.Authenticate", &req, &rsp)
	if rsp.Error != NoError {
		log.Printf("peerA % failed to be authenticated, err: %s", p.name, rsp.Error)
		return
	}

	flag := p.cryptoMachine.VerifyCertificate(rsp.PeerBCertificateBytes, rsp.PeerBPublicKeyBytes)
	if !flag {
		log.Printf("peerB %s authentication failed", peeraddr)
		return
	}
	log.Printf("peerA %s and peerB %s authentication succeeded", p.name, peeraddr)
}

func (p *Peer) server(peerName string) {
	rpc.Register(p)
    rpc.HandleHTTP()

	os.Remove(peerName)
	l, e := net.Listen("unix", peerName)
	if e !=nil {
		log.Printf("peer %s listen error: %s", peerName, e.Error())
	}

	go http.Serve(l, nil)
}

func call(machime string, rpcname string, req interface{}, rsp interface{}) {
	c, err := rpc.DialHTTP("unix", machime)
	defer c.Close()
	if err != nil {
		log.Fatal("dialing error: ", err.Error())
		return
	}

	err = c.Call(rpcname, req, rsp)
	if err != nil {
		log.Fatal("call %s method %s error: ", err.Error())
	}
}

func MakePeer(peerName string) *Peer {
	p := &Peer{}

	p.once.Do(func() {
		p.name = peerName
		p.password = "123456.a"

		p.cryptoMachine = &cryptoalgs.Ecc{}
	})
	p.cryptoMachine.GenerateKeys()
	p.register()
	p.requireCertification()

	p.server(peerName)

	return p
}
