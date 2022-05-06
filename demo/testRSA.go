package main

import (
	"P2pSecuritySimulator/cryptoalgs"
	"fmt"
)

func main() {
	var machine cryptoalgs.CryptoMachine
	machine = &cryptoalgs.Rsa{}

	machine.GenerateKeys()
	text := []byte("hello world")

	res := machine.Encrypt(text)
	fmt.Println(res)
	decrypted := machine.Decrypt(res)
	fmt.Println(decrypted)

	cert := machine.GenerateCertificate(text)
	fmt.Println(machine.VerifyCertificate(cert, machine.GetPublicKeyBytes()))
}
