package cert_test

import (
	"btcDemo/cert"
	"fmt"
	"testing"
)

var certService cert.CertService

func TestGenerateSimplePrivateKey(t *testing.T) {
	key, err := certService.GenerateSimplePrivateKey()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("privatekey:%s\n\rpublickey:%s", key.PrivKey, key.PubKey)
}
