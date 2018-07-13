package coins

import (
	"encoding/json"
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

//获取新地址，同事数据库会存储key以便调试
func TestGetNewAddress(t *testing.T) {
	address, account, err := btc.GetNewAddress("Test" + strconv.FormatInt(time.Now().Unix(), 10))
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("address:%s\n\raccount:%s\n\r", address, account)
}

func TestGetBalanceInAddress(t *testing.T) {
	balance, err := btc.GetBalanceInAddress("ms8d4chAKH9CjTY57HNymFSLZNUkZXFnVY")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(balance)
}
func TestGetUnspentByAddress(t *testing.T) {
	unspents, err := btc.GetUnspentByAddress("mkxMPobtVtgYVXfY2yw8jKfaWHxSbEyGoQ")
	if err != nil {
		t.Error(err)
	}
	for k, v := range unspents {
		model, err := json.MarshalIndent(v, "", " ")
		if err != nil {
			t.Error(err)
		}
		fmt.Printf("%d\r\n%s", k, model)
	}
}

func TestSendAddressToAddress(t *testing.T) {
	err := btc.SendAddressToAddress("mkxMPobtVtgYVXfY2yw8jKfaWHxSbEyGoQ", "mhAfGecTPa9eZaaNkGJcV7fmUPFi3T2Ki8", 40, 0.0001)
	if err != nil {
		t.Error(err)
	}
}

func TestGetTxByAddress(t *testing.T) {
	txs, err := btc.GetTxByAddress("mkxMPobtVtgYVXfY2yw8jKfaWHxSbEyGoQ")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(txs)
}
