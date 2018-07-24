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
