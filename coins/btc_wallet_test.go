package coins

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

var btc BtcService

func TestGetAccounts(t *testing.T) {
	acts, _ := btc.GetAccounts()
	for _, v := range acts {
		fmt.Println(*v)
	}
}

//cert cert_test/TestNewAddress()可以得到key
func TestCheckAddressExisted(t *testing.T) {
	/* privatekey:92QiFfPkAfafdtTW5a8eCqLgCKK1tEZKMcAGA3PVi79cJpZeujc
	   publickey:046c9bbd1c67db7a99bb45a98c592ec89bffe65174ddd130395d632cb428f7423c3cc4de7d623bc4da321451ddede0e39e8bec0105103268e609cb175ea2fedf91
	   n4Wxwu3xQe7vWQoqjzbjPmMMewBYjhcZzn
	*/
	address, err := btc.CheckAddressExisted("046c9bbd1c67db7a99bb45a98c592ec89bffe65174ddd130395d632cb428f7423c3cc4de7d623bc4da321451ddede0e39e8bec0105103268e609cb175ea2fedf91")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(address.EncodeAddress())
	if address.EncodeAddress() != "n4Wxwu3xQe7vWQoqjzbjPmMMewBYjhcZzn" {
		t.Error("失败")
	}
}

func TestGetNewAddress(t *testing.T) {
	address, account, err := btc.GetNewAddress("Test" + strconv.FormatInt(time.Now().Unix(), 10))
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("address:%s\n\raccount:%s", address, account)
}
