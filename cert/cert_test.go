package cert_test

import (
	"blockchainDemo/cert"
	"encoding/json"
	"fmt"
	"testing"
)

var btc *cert.BtcCertService
var xlm *cert.XlmCertService
var handler cert.CertHandler

func TestGenerateSimpleKey(t *testing.T) {
	handler.LoadService(xlm)
	key, err := handler.GenerateSimpleKey()
	if err != nil {
		fmt.Println(err)
	}
	str, _ := json.Marshal(key)
	fmt.Println(string(str))
	//fmt.Printf("privatekey:%s\n\rpublickey:%s", key.PrivKey, key.PubKey)
}
func TestNewAddress(t *testing.T) {
	handler.LoadService(xlm)
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
