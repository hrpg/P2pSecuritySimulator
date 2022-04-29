package services

import (
	"P2pSecuritySimulator/cryptoalgs"
	"P2pSecuritySimulator/dataCollector"
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

func (p *Peer) Authenticate(req *AuthenticateReq, rsp *AuthenticateRsp) error {
	flag := p.cryptoMachine.VerifyCertificate(req.PeerACertificateBytes, p.serverPubKeyBytes)
	if !flag {
		rsp.Error = ErrAuthenticationFailed
		return nil
	}

	rsp.PeerBCertificateBytes = p.myCertificate
	rsp.PeerBPublicKeyBytes = p.cryptoMachine.GetPublicKeyBytes()
	rsp.Error = NoError

	return nil
}

func (p *Peer) register() {
	for i := 0; i < 3; i++  {
		req := RegisterReq{}
		rsp := RegisterRsp{}

		req.Name, req.PassWord = p.name, p.password
		call("/var/tmp/server", "AuthenticationServer.Register", &req, &rsp)
		if rsp.Error == NoError {
			p.serverPubKeyBytes = rsp.ServerPubKeyBytes
			log.Printf("PEER: peer %s successfully registered", p.name)
			break
		}

		log.Printf("PEER: peer %s failed to register, error: %s", p.name, rsp.Error)
		rand.Seed(time.Now().UnixNano())

		sleepTime := 200 + rand.Intn(200)
		time.Sleep(time.Duration(sleepTime) * time.Millisecond)
	}
}

func (p *Peer) requestCertification() {
	log.Printf("PEER: peer %s start to request certificate", p.name)
	for i := 0; i < 3; i++ {
		start := time.Now()
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

		call("/var/tmp/server", "AuthenticationServer.AssignCertificate", &req, &rsp)
		if rsp.Error == NoError {
			p.myCertificate = p.cryptoMachine.Decrypt(rsp.EncryptedPeerCertificateBytes)
			dataCollector.AppendRequireCertificateTime(time.Since(start).Milliseconds())
			log.Printf("PEER: peer %s get cerfificate", p.name)
			break
		}

		log.Printf("PEER: peer %s failed to ask certificate, error: %s", p.name, rsp.Error)
		rand.Seed(time.Now().UnixNano())

		sleepTime := 200 + rand.Intn(200)
		time.Sleep(time.Duration(sleepTime) * time.Millisecond)
	}
}

func (p *Peer) RequestAuthentication(peeraddr string) {
	start := time.Now()
	req := AuthenticateReq{}
	rsp := AuthenticateRsp{}
	req.PeerACertificateBytes = p.myCertificate
	req.PeerAPublicKeyBytes = p.cryptoMachine.GetPublicKeyBytes()

	call(peeraddr, "Peer.Authenticate", &req, &rsp)
	if rsp.Error != NoError {
		log.Printf("PEER: peerA %s failed to be authenticated, err: %s", p.name, rsp.Error)
		return
	}

	log.Printf("PEER: peerA %s successfully authenticated", p.name)
	flag := p.cryptoMachine.VerifyCertificate(rsp.PeerBCertificateBytes, p.serverPubKeyBytes)
	if !flag {
		log.Printf("PEER: peerB %s authentication failed", peeraddr)
		return
	}

	dataCollector.AppendAuthentificateTime(time.Since(start).Milliseconds())
	log.Printf("PEER: peerA %s and peerB %s authentication succeeded", p.name, peeraddr)
}

func (p *Peer) server(peerName string) {
	rpc.Register(p)
    rpc.HandleHTTP()

	os.Remove(peerName)
	l, e := net.Listen("unix", peerName)
	if e !=nil {
		log.Fatalf("PEER: peer %s listen error: %s", peerName, e.Error())
	}

	go http.Serve(l, nil)
}

func call(machime string, rpcname string, req interface{}, rsp interface{}) {
	c, err := rpc.DialHTTP("unix", machime)
	defer c.Close()
	if err != nil {
		log.Fatal("PEER: dialing error: ", err.Error())
		return
	}

	err = c.Call(rpcname, req, rsp)
	if err != nil {
		log.Fatalf("PEER: call %s method %s error: ", err.Error())
	}
}

func MakePeer(peerName string) *Peer {
	time.Sleep(time.Second * 2)
	p := &Peer{}

	p.once.Do(func() {
		p.name = peerName
		p.password = "123456.a"

		p.cryptoMachine = &cryptoalgs.Ecc{}
	})
	p.cryptoMachine.GenerateKeys()
	log.Printf("PEER: peer %s has generated keys", peerName)
	p.register()
	log.Printf("PEER: peer %s has registered", peerName)
	p.requestCertification()

	p.server(peerName)
	time.Sleep(time.Minute * 5)
	log.Printf("PEER: peer %s start to run", peerName)
	return p
}
