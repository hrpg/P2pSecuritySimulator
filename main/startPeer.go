package main

import (
	"P2pSecuritySimulator/services"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: peer addresses...\n")
		os.Exit(1)
	}

	myaddr := os.Args[1]
	peer := services.MakePeer2(os.Args[1])
	peerAddresses := os.Args[2:]
	for _, addr := range peerAddresses  {
		if myaddr == addr {
			continue
		}

		peer.RequestAuthentication(addr)
		rand.Seed(time.Now().UnixNano())

		sleepTime := 10 + rand.Intn(200)
		time.Sleep(time.Duration(sleepTime) * time.Millisecond)
	}

	log.Printf("peer %s task has completed", myaddr)
	time.Sleep(time.Second * 5)
}
