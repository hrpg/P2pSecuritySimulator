package main

import (
	"P2pSecuritySimulator/services"
	"fmt"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Fprintf(os.Stderr, "Usage: peer addresses...\n")
		os.Exit(1)
	}

	peer := services.MakePeer(os.Args[1])
	peerAddresses := os.Args[2:]
	for _, addr := range peerAddresses  {
		peer.RequireAuthentication(addr)
		//rand.Seed(time.Now().UnixNano())
		//
		//sleepTime := 10 + rand.Intn(200)
		//time.Sleep(time.Duration(sleepTime) * time.Millisecond)
	}

	time.Sleep(time.Second)
}
