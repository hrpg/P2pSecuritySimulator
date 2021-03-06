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

type StatusCode2 int

const (
	FirstConformed StatusCode2 = 1
	LastConformed StatusCode2 = 2
)

type Peer struct {
	mux sync.RWMutex
	once sync.Once
	RestartWork chan int

	name string
	password string

	cryptoMachine cryptoalgs.CryptoMachine
	signMachine cryptoalgs.CryptoMachine

	serverCryptoPubKeyBytes []byte
	serverSignPubKeyBytes []byte
	myCertificate []byte

	connectStatus map[string]StatusCode2
}

func (p *Peer) Report() {
	req := ReportWorkDoneReq{}
	rsp := ReportWorkDoneRsp{}

	call("/var/tmp/server", "AuthenticationServer2.ReportDone", &req, &rsp)
}

func (p *Peer) CanRestartWork(req *CanRestartWorkReq, rsp *CanRestartWorkRsp) error {
	p.RestartWork <- 1
	return nil
}

func (p *Peer) Finalize(req *FinalizeReq, rsp *FinalizeRsp) error {
	p.mux.Lock()
	defer p.mux.Unlock()
	if status, ok := p.connectStatus[req.PeerName]; !ok || status != FirstConformed || req.Echo != NoError {
		delete(p.connectStatus, req.PeerName)
		log.Printf("PEER: peer %s's connection failed", req.PeerName)
		return nil
	}

	p.connectStatus[req.PeerName] = LastConformed
	return nil
}

func (p *Peer) Authenticate(req *AuthenticateReq, rsp *AuthenticateRsp) error {
	flag := p.signMachine.VerifyCertificate(req.PeerACertificateBytes, p.serverSignPubKeyBytes)
	if !flag {
		rsp.Error = ErrAuthenticationFailed
		return nil
	}

	myCertAndPubkeyInfo := PeerBCertAndPubKeyInfo{
		PeerBCertificateBytes: p.myCertificate,
		PeerBCryptoPublicKeyBytes: p.cryptoMachine.GetPublicKeyBytes(),
	}
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	enc.Encode(myCertAndPubkeyInfo)

	rsp.PeerBCertAndCryptoPubKeyInfoBytes = p.cryptoMachine.EncryptWithPubKey(buffer.Bytes(), req.PeerACryptoPublicKeyBytes)
	rsp.Error = NoError

	p.mux.Lock()
	p.connectStatus[req.PeerName] = FirstConformed
	p.mux.Unlock()

	return nil
}


func (p *Peer) Register() {
	req := RegisterReq{}
	rsp := RegisterRsp{}
	req.Name, req.PassWord = p.name, p.password
	for i := 0; i < 3; i++  {

		call("/var/tmp/server", "AuthenticationServer2.Register", &req, &rsp)
		if rsp.Error == NoError {
			p.serverCryptoPubKeyBytes = rsp.ServerCryptoPubKeyBytes
			p.serverSignPubKeyBytes = rsp.ServerSignPubKeyBytes
			log.Printf("PEER: peer %s successfully registered", p.name)
			break
		}

		log.Printf("PEER: peer %s failed to register, error: %s", p.name, rsp.Error)
		rand.Seed(time.Now().UnixNano())

		sleepTime := 20 + rand.Intn(20)
		time.Sleep(time.Duration(sleepTime) * time.Millisecond)
	}
}

func (p *Peer) RequestCertification() {
	log.Printf("PEER: peer %s start to request certificate", p.name)
	req := GetCertificateReq{}
	rsp := GetCertificateRsp{}
	var peerInfo PeerInfo
	peerInfo.Name = p.name
	peerInfo.Password = p.password
	peerInfo.PeerCryptoPublicKeyBytes = p.cryptoMachine.GetPublicKeyBytes()
	for i := 0; i < 3; i++ {
		start := time.Now()

		var buffer bytes.Buffer
		enc := gob.NewEncoder(&buffer)
		enc.Encode(peerInfo)
		peerInfoBytes := buffer.Bytes()

		encryptedPeerInfoBytes := p.cryptoMachine.EncryptWithPubKey(peerInfoBytes, p.serverCryptoPubKeyBytes)
		req.EncryptedPeerInfoBytes = encryptedPeerInfoBytes

		call("/var/tmp/server", "AuthenticationServer2.AssignCertificate", &req, &rsp)
		if rsp.Error == NoError {
			p.myCertificate = p.cryptoMachine.Decrypt(rsp.EncryptedPeerCertificateBytes)
			dataCollector.AppendRequireCertificateTime(time.Since(start).Nanoseconds())
			log.Printf("PEER: peer %s get cerfificate", p.name)
			break
		}

		log.Printf("PEER: peer %s failed to ask certificate, error: %s", p.name, rsp.Error)
		rand.Seed(time.Now().UnixNano())

		sleepTime := 20 + rand.Intn(20)
		time.Sleep(time.Duration(sleepTime) * time.Millisecond)
	}
}

func (p *Peer) RequestAuthentication(peeraddr string) {
	start := time.Now()
	req := AuthenticateReq{}
	rsp := AuthenticateRsp{}
	req.PeerACertificateBytes = p.myCertificate
	req.PeerName = p.name
	req.PeerACryptoPublicKeyBytes = p.cryptoMachine.GetPublicKeyBytes()

	call(peeraddr, "Peer2.Authenticate", &req, &rsp)
	if rsp.Error != NoError {
		log.Printf("PEER: peerA %s failed to be authenticated, err: %s", p.name, rsp.Error)
		return
	}

	log.Printf("PEER: peerA %s successfully authenticated", p.name)
	// ???????????????
	decryptedText := p.cryptoMachine.Decrypt(rsp.PeerBCertAndCryptoPubKeyInfoBytes)
	// ???????????????
	buffer := bytes.NewBuffer(decryptedText)
	dec := gob.NewDecoder(buffer)
	var peerBCertAndPubKey PeerBCertAndPubKeyInfo
	dec.Decode(&peerBCertAndPubKey)

	flag := p.signMachine.VerifyCertificate(peerBCertAndPubKey.PeerBCertificateBytes, p.serverSignPubKeyBytes)
	replyB := FinalizeReq{
		p.name,
		NoError,
	}
	if !flag {
		replyB.Echo = ErrBCertInfo
		call(peeraddr, "Peer2.Finalize", &replyB, nil)
		log.Printf("PEER: peerB %s authentication failed", peeraddr)
		return
	}

	call(peeraddr, "Peer2.Finalize", &replyB, nil)
	dataCollector.AppendAuthentificateTime(time.Since(start).Nanoseconds())
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
	p := &Peer{}

	p.once.Do(func() {
		p.RestartWork = make(chan int, 1)
		p.name = peerName
		p.password = "123456.a"

		p.cryptoMachine = &cryptoalgs.Ecc{}
		p.signMachine = &cryptoalgs.Dsa{}

		p.connectStatus = make(map[string]StatusCode2)
	})

	p.cryptoMachine.GenerateKeys()
	log.Printf("PEER: peer %s has generated keys", peerName)
	p.Register()
	log.Printf("PEER: peer %s has registered", peerName)
	p.RequestCertification()

	p.server(peerName)
	time.Sleep(time.Second * 5)
	log.Printf("PEER: peer %s start to run", peerName)
	return p
}
