package services

import (
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
	privateKey string
	publicKey string
	authenticationServerKey string
	certification []byte
}

func (p *Peer) register() {
	for true {
		req := RegisterReq{}
		rsp := RegisterRsp{}

		req.Name, req.PassWord = p.name, p.password
		call("/var/tmp/server", "AuthenticationServer.Register", &req, &rsp)
		if rsp.Error == NoError {
			break
		}

		log.Fatalf("peer %s failed to register, error: %s", p.name, rsp.Error)
		rand.Seed(time.Now().UnixNano())

		sleepTime := 200 + rand.Intn(200)
		time.Sleep(time.Duration(sleepTime) * time.Millisecond)
	}
}

func (p *Peer) server(peerName string) {
	rpc.Register(p)
    rpc.HandleHTTP()

	os.Remove(peerName)
	l, e := net.Listen("unix", peerName)
	if e !=nil {
		log.Fatalf("peer %s listen error: %s", peerName, e.Error())
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
	})

	p.server(peerName)

	return p
}
