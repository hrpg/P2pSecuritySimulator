package main

import (
	"P2pSecuritySimulator/services"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: server input nPeers...\n")
		os.Exit(1)
	}

	nPees, _ := strconv.Atoi(os.Args[1])
	server := services.MakeAuthenticationServer2(nPees)
	for server.Done() == false {
		time.Sleep(time.Second)
	}

	log.Printf("server's task has completed")
	time.Sleep(time.Second)
}
