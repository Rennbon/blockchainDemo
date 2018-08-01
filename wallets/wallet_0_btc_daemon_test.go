package wallets

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"testing"
	"time"
	"log"
)


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

func TestNewBTCDaemon(t *testing.T) {
	daemon := NewBTCDaemon(time.NewTicker(5*time.Second))
	t.Log( daemon.blkHt)
}
func TestMonitoringBtcBlockHeight(t *testing.T) {
	daemon := NewBTCDaemon(time.NewTicker(1*time.Second))
	height :=daemon.blkHt
	t.Log("目前高度:",height)
	log.Println("目前高度:",height)
	for i:=0;i<100;i++ {
		select {
			case <-daemon.tick.C:
				log.Println("第",i,"次")
				daemon.monitoringBtcBlockHeight()
				log.Println("current block height",daemon.blkHt)
				log.Println("old block height",height)
			    if daemon.blkHt>height+2{
			    	log.Println("I am coming!!!")
			    	return
				}
		}
	}
}
