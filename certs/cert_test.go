package certs_test

import (
	"fmt"
	"github.com/Rennbon/blockchainDemo/certs"
	"reflect"
	"testing"
)

type CertHandler struct {
	certs.Generater
	TypeName string
}

func (ch *CertHandler) LoadService(g certs.Generater) error {
	if g != nil {
		ch.Generater = g
	}
	typ := reflect.TypeOf(g)
	ch.TypeName = typ.String()
	return nil
}

var (
	btc     *certs.BtcCertService
	xlm     *certs.XlmCertService
	eth     *certs.EthCertService
	handler CertHandler
)

func TestGenerateSimpleKey(t *testing.T) {
	handler.LoadService(eth)
	key, err := handler.GenerateSimpleKey()
	if err != nil {
		fmt.Println(err)
	}
	t.Log(key)
	//fmt.Printf("privatekey:%s\n\rpublickey:%s", key.PrivKey, key.PubKey)
}
func TestNewAddress(t *testing.T) {
	handler.LoadService(btc)
	/* btc
	   privatekey:92QiFfPkAfafdtTW5a8eCqLgCKK1tEZKMcAGA3PVi79cJpZeujc
	   publickey:046c9bbd1c67db7a99bb45a98c592ec89bffe65174ddd130395d632cb428f7423c3cc4de7d623bc4da321451ddede0e39e8bec0105103268e609cb175ea2fedf91
	   address:n4Wxwu3xQe7vWQoqjzbjPmMMewBYjhcZzn
	   xlm
	   seed:
	*/
	input := ""
	output := ""

	switch handler.TypeName {
	case "*cert.BtcCertService":
		input = "046c9bbd1c67db7a99bb45a98c592ec89bffe65174ddd130395d632cb428f7423c3cc4de7d623bc4da321451ddede0e39e8bec0105103268e609cb175ea2fedf91"
		output = "n4Wxwu3xQe7vWQoqjzbjPmMMewBYjhcZzn"
		break
	case "*cert.XlmCertService":
		input = "SCK3NIJ5XLQS3E7OP3TVFYBYXXXIHA2ILKCW6PDYFHMABRPFUIV2HAE4"
		output = "GC7JBI22JROCC5T5ROWT4RDB4C4IHNZPUSQVIAPICJHGSMXWDW5TVKDF"
	}
	addr, _ := handler.GetNewAddress(input)
	fmt.Println(addr)
	if output != addr {
		t.Error("失败")
	}

}

func Test(t *testing.T) {
	fmt.Println(123)
	fmt.Println("12312312")
}
