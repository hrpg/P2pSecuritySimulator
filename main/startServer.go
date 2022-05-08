package main

import (
	"P2pSecuritySimulator/services"
)

func main() {
	//if len(os.Args) < 2 {
	//	fmt.Fprintf(os.Stderr, "usage: server input nPeers...\n")
	//	os.Exit(1)
	//}
	//
	//nPees, _ := strconv.Atoi(os.Args[1])
	server := services.MakeAuthenticationServer2()
	ch := make(chan bool)
	for server.Done() == false {
		ch <- true
	}

	//log.Printf("server's task has completed")
	//time.Sleep(time.Second * 2)
}
