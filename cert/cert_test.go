package cert_test

import (
	"btcDemo/cert"
	"fmt"
	"testing"
)

var certService cert.CertService

func TestGenerateSimpleKey(t *testing.T) {
	key, err := certService.GenerateSimpleKey()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("privatekey:%s\n\rpublickey:%s", key.PrivKey, key.PubKey)
}
func TestNewAddress(t *testing.T) {
	//msKTSXKQKXjFsHZit1LDiV4SWacY5kNDEz
	//privatekey:5KE6Brfr1bVN16md5UVaY4kAtx8GkThzBVk2QeLS422D65Cypkf
	//publickey:04ecf7652cbaa4e504d9ab37032a18386a771baf6f01e15047cc8d42a29a99f2eab6a32b69a6f2a779f4202826c287aedc3b2d27278b56a7cd1312eba4e72a8a35
	addr, _ := certService.GetNewAddress("04ecf7652cbaa4e504d9ab37032a18386a771baf6f01e15047cc8d42a29a99f2eab6a32b69a6f2a779f4202826c287aedc3b2d27278b56a7cd1312eba4e72a8a35")
	fmt.Println(addr)
	if "msKTSXKQKXjFsHZit1LDiV4SWacY5kNDEz" != addr {
		t.Error("失败")
	}
}
