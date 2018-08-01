package wallets

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"testing"
	"time"
)

var _t_btcDaemon  = NewBTCDaemon(time.NewTicker(5*time.Second))

///////////////////////////////////pool-test////////////////////////////////////////////////
func TestFillBlockHeight(t *testing.T) {
	hash, _ := chainhash.NewHashFromStr("88af33f7e455751ae130746f7a7bd6538fc8b791aa67692dba33ab450ada9c92")
	_txc := &txexcuting{
		txHash: hash,
		txcache: &txcache{
			txrchan: make(chan *TxResult),
		},
	}
	_txc.fillBlockHeight()
	t.Log(_txc.targetH, _txc.blockH)
}
func TestMonitoringBtcBlockHeight(t *testing.T) {
	//	monitoringBtcBlockHeight()
}
