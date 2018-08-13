/*这里是单账户模式
todo:依赖注入的方式实现daemon层动态适配不同的币种，主流程只维护一套

*/

package daemon

import (
	"github.com/Rennbon/blockchainDemo/coins"
	"github.com/Rennbon/blockchainDemo/wallets"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"sync"
	"time"
)

type simpleTx struct {
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
	address string //账户地址
	balance   coins.CoinAmounter //账户余额
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
	var err error
	d := &daemon{
		ispkg: dd.isPkg(),
		tick:  dd.blockTick(),
	}
	d.blkHt, err = dd.getBlockHeight()
	d.address = "" //这里应该是本地数据库读去address配置
	d.balance,err = dd.getBalance(d.address)//账户余额
	if err != nil {
		panic(err)
	}
	dd.
    d.l = &localPool{
    	m:new(sync.Mutex),
		txs:[]*simpleTx{},
		size:0,
		txsAmount:
	}

	if err!=nil{
		panic(err)
	}

	txs       []*simpleTx        //交易单s
	size      int                //txs的数量
	txsAmount coins.CoinAmounter //txs聚合总额
	inlock    bool               //填充锁，决定txs是否可以append新tx
	outlock   bool               //消费锁，决定txs是否可以执行消费流程
	inlockch  chan lockType      //inlock锁通道
	outlockch chan lockType      //outlock锁通道
	return d
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

//以下方法需要动态注入到daemon里面
type distributionDaemoner interface {
	isPkg() bool
	//获取总块高
	getBlockHeight() (height int64, err error)
	//获取tx成功状态的区块确认数，秒过的返回1
	getSuccessfulConfirmedNum()(minBlock int64)
	//验证txId对应的tx的确认状态
	checkConfirm(txId string) (confimNum int64,err error)
	//获取账户余额
	getBalance(address string )(balance coins.CoinAmounter, err error)
	//区块同步时间
	blockTick() (tick *time.Ticker)
}
