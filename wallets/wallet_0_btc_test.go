package wallets

import (
	"encoding/json"
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

func TestGetUnspentByAddress(t *testing.T) {
	unspents, err := getUnspentByAddress("mkxMPobtVtgYVXfY2yw8jKfaWHxSbEyGoQ")
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

func TestGetTxByAddress(t *testing.T) {
	txs, err := btc.GetTxByAddress([]string{"2NBpzw8BLKhES9MyM7gt7Crp1PWckFvsYFn"}, "")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(txs)
}

func TestGetRawTransaction(t *testing.T) {
	getRawTransaction("f1951a8a7f5028a58cef7dace3ec64b1a6ab5f0711bff6baf2d7cd7ee64d9424")
}
func TestGetBlockInfo(t *testing.T) {
	getBlockInfo("53dc56749eac5f46820fcdee93e0c1e4242b07a129b4635cc6e4e57c4d69ba76")
}
