package main

import (
	"P2pSecuritySimulator/cryptoalgs"
	"fmt"
)

func main() {
	var machine cryptoalgs.CryptoMachine

	ecc := cryptoalgs.Ecc{}
	machine = &ecc
	machine.GenerateKeys()

	text := []byte("hello world")
	cert := machine.GenerateCertificate(text)

	res := machine.VerifyCertificate(cert, machine.GetPublicKeyBytes())
	fmt.Println(res)
}