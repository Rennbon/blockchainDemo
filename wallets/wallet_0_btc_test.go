package wallets

import (
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
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
	getRawTransaction("e3f0129da50920c2a01492eecc9c13acbefd9b8a9a46af72626141f23c774030")
}
func TestGetBlockInfo(t *testing.T) {
	getBlockInfo("53dc56749eac5f46820fcdee93e0c1e4242b07a129b4635cc6e4e57c4d69ba76")
}

///////////////////////////////////pool-test////////////////////////////////////////////////
func TestFillBlockHeight(t *testing.T) {
	hash, _ := chainhash.NewHashFromStr("88af33f7e455751ae130746f7a7bd6538fc8b791aa67692dba33ab450ada9c92")
	txc := &txexcuting{
		txHash: hash,
		txcache: &txcache{
			txrchan: make(chan *TxResult),
		},
	}
	txc.fillBlockHeight()
	t.Log(txc.targetH, txc.blockH)
}
func TestMonitoringBtcBlockHeight(t *testing.T) {
	monitoringBtcBlockHeight()
}
