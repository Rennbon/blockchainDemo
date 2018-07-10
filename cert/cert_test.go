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
	/* privatekey:92QiFfPkAfafdtTW5a8eCqLgCKK1tEZKMcAGA3PVi79cJpZeujc
	   publickey:046c9bbd1c67db7a99bb45a98c592ec89bffe65174ddd130395d632cb428f7423c3cc4de7d623bc4da321451ddede0e39e8bec0105103268e609cb175ea2fedf91
	*/
	addr, _ := certService.GetNewAddress("046c9bbd1c67db7a99bb45a98c592ec89bffe65174ddd130395d632cb428f7423c3cc4de7d623bc4da321451ddede0e39e8bec0105103268e609cb175ea2fedf91")
	fmt.Println(addr)
	if "n4Wxwu3xQe7vWQoqjzbjPmMMewBYjhcZzn" != addr {
		t.Error("失败")
	}
}
