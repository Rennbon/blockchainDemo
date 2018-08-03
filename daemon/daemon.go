package daemon

import (
	"github.com/Rennbon/blockchainDemo/coins"
	"github.com/Rennbon/blockchainDemo/wallets"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"sync"
	"time"
)

type simpleTx struct {
	addrF      string                //提币地址
	addrTs     []*wallets.AddrAmount //充币地址
	createTime time.Time             //创建时间
	deadline   time.Time             //执行期限，晚于此线执行失败逻辑
}
type TxOngoing struct {
	*simpleTx
	txHash  *chainhash.Hash //txhash
	blockH  int64           //所在的区块高度
	confirm int64           //公链确认次数

}
type lockType bool

const (
	lock   lockType = true
	unlock lockType = false
)

type localPool struct {
	m         *sync.Mutex        //抢占用锁，确保原子性
	balance   coins.CoinAmounter //账户余额
	txs       []*simpleTx        //交易单s
	size      int                //txs的数量
	txsAmount coins.CoinAmounter //txs聚合总额
	inlock    bool               //填充锁，决定txs是否可以append新tx
	outlock   bool               //消费锁，决定txs是否可以执行消费流程
	inlockch  chan lockType      //inlock锁通道
	outlockch chan lockType      //outlock锁通道
}

type globalPool struct {
	m    *sync.Mutex //抢占用锁，确保原子性
	txs  []*simpleTx //交易单s
	size int         //txs的数量
}

type daemon struct {
	ispkg    bool         //是否支持打包处理tx
	tick     *time.Ticker //周期计时器
	blkHt    int64        //区块块高
	l        *localPool
	g        *globalPool
	excuTx   chan []*simpleTx //执行tx交易通道
	listenTx chan *TxOngoing  //tx监听
	dd       distributionDaemoner
}

func newDaemon(dd distributionDaemoner) daemoner {
	d := &daemon{
		ispkg: dd.isPkg(),
		tick:  dd.blockTick(),
	}
	return d
}
func (l *localPool) newLocalPool() {

}

func (d *daemon) run() {

}
func (d *daemon) stop() {

}

//todo 消费excutx中的数据进行交易
func (d *daemon) consumeExcutx() error {

}

type daemoner interface {
	run()
	stop()
	//消费交易
	consumeExcutx() error
}
type distributionDaemoner interface {
	isPkg() bool
	//获取总块高
	getBlockHeight() (height int64, err error)
	//获取tx成功状态的区块确认数，秒过的返回1
	getSuccessfulConfirmedNum() int64
	//验证txId对应的tx的确认状态
	checkConfirm(txId string) (int64, error)
	blockTick() (tick *time.Ticker)
}
