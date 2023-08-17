package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func GenPubandPriKey(bits int, filePath string) error {
	// crate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	priBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}

	err = os.WriteFile(filePath+"private.pem", pem.EncodeToMemory(priBlock), 0644)
	if err != nil {
		return err
	}
	fmt.Println("private key crate success!")

	publicKeyr := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKeyr)
	if err != nil {
		return err
	}

	publicBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}

	err = os.WriteFile(filePath+"public.pem", pem.EncodeToMemory(publicBlock), 0644)
	if err != nil {
		return err
	}

	fmt.Println("public key crate success!")

	return nil
}
