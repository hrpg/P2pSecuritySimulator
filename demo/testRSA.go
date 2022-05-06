package main

import (
	"P2pSecuritySimulator/cryptoalgs"
	"fmt"
)

func main() {
	var machine cryptoalgs.CryptoMachine
	machine = &cryptoalgs.Rsa{}

	machine.GenerateKeys()
	cert := machine.GenerateCertificate(machine.GetPublicKeyBytes())
	encryptedCert := machine.EncryptWithPubKey(cert, machine.GetPublicKeyBytes())
	decryptedCert := machine.Decrypt(encryptedCert)

	flag := machine.VerifyCertificate(decryptedCert, machine.GetPublicKeyBytes())
	fmt.Println(flag)
}
