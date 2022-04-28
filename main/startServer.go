package main

import (
	"P2pSecuritySimulator/services"
	"fmt"
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
	server := services.MakeAuthentication(nPees)
	for server.Done() == false {
		time.Sleep(time.Second)
	}

	time.Sleep(time.Second)
}
