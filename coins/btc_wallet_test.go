package coins

import (
	"fmt"
	"testing"
)

var btc BtcService

func TestGetAccounts(t *testing.T) {
	acts, _ := btc.GetAccounts()
	for _, v := range acts {
		fmt.Println(*v)
	}
}

func TestCheckAddressExisted(t *testing.T) {
	//msKTSXKQKXjFsHZit1LDiV4SWacY5kNDEz
	//privatekey:5KE6Brfr1bVN16md5UVaY4kAtx8GkThzBVk2QeLS422D65Cypkf
	//publickey:04ecf7652cbaa4e504d9ab37032a18386a771baf6f01e15047cc8d42a29a99f2eab6a32b69a6f2a779f4202826c287aedc3b2d27278b56a7cd1312eba4e72a8a35
	address, err := btc.CheckAddressExisted("04ecf7652cbaa4e504d9ab37032a18386a771baf6f01e15047cc8d42a29a99f2eab6a32b69a6f2a779f4202826c287aedc3b2d27278b56a7cd1312eba4e72a8a35")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(address.EncodeAddress())
	if address.EncodeAddress() != "msKTSXKQKXjFsHZit1LDiV4SWacY5kNDEz" {
		t.Error("失败")
	}
}
