package main

import (
	"P2pSecuritySimulator/cryptoalgs"
	"fmt"
)

func main() {
	var machine cryptoalgs.CryptoMachine
	machine = &cryptoalgs.Dsa{}

	machine.GenerateKeys()

	text := []byte("hello world")
	cert := machine.GenerateCertificate(text)
	res := machine.VerifyCertificate(cert, machine.GetPublicKeyBytes())

	fmt.Println(res)
}
