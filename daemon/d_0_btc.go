package daemon

import (
	"github.com/Rennbon/blockchainDemo/wallets"
	"time"
)

type btcDaemon struct {
	*wallets.BtcService
	daemoner
}

func newBtcDaemon() *btcDaemon {
	d := &btcDaemon{
		BtcService: &wallets.BtcService{},
	}
	d.daemoner = newDaemon(d)
	return d
}

func (*btcDaemon) getBlockHeight() (height int64, err error) {
	return
}

//获取tx成功状态的区块确认数，秒过的返回1
func (*btcDaemon) getSuccessfulConfirmedNum() int64 {
	return 0
}

//验证txId对应的tx的确认状态
func (*btcDaemon) checkConfirm(txId string) (int64, error) {
	return 0, nil
}
func (*btcDaemon) isPkg() bool {
	return true
}

func (*btcDaemon) blockTick() (tick *time.Ticker) {
	return time.NewTicker(5 * time.Second)
}
