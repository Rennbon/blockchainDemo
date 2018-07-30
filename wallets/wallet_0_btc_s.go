
/*
维护本地队列，以及全局，历史队列
 - 本地队列限制单位时间只能允许定量的tx提交到公链
 - 全局队列用来维护，超出本地队列长度限制后的将要执行tx交易的参数
 - 历史队列为发起tx交易，并检测tx状态用

1. 需要轮询读取block height
2. 需要给tx赋值当前区块高度及验证区块高度
3. 历史池需要遍历获取池中数据的状态并维护，每隔一定时间清除状态已确认的数据

 */
package wallets

import (
	"github.com/Rennbon/blockchainDemo/coins"
	"sync"
	"time"
	"btcd/chaincfg/chainhash"
)

//本地容器上限
const localPoolCount = int(10)

//local池，
type btcLocalPool struct {
	m         sync.RWMutex
	deadline  time.Time  //死亡时间
	size      int        //数量
	txcs []*txcache //tx组
}
//全局池，只维护计数
type btcGlobalPool struct {
	m         sync.RWMutex
	size int        //数量
}
//历史池，这里才开始执行增伤查
type btcHistoryPool struct {
	m sync.RWMutex
	size int //总数
	txcsing []* txexcuting//处理的交易
	offset int //当前处理位置，指提交交易并广播
}



type txexcuting struct {
	txHash *chainhash.Hash //
	status bool //状态
	*txcache
}
type txcache struct {
	birthday time.Time          //出生时间
	addf     string             //来自地址
	addt     string             //去向地址
	transfer coins.CoinAmounter //交易金额
	fee      coins.CoinAmounter //交易小费
	txrchan  <-chan *TxResult //chan维持
}
var btcGPL *btcGlobalPool
var btcLPL *btcLocalPool
var btcHPL *btcHistoryPool

var btcTimeM,_ = time.ParseDuration("10m")

//计时器每满一次，清空local池，移动到history池
func (*btcLocalPool) Restart() {
	btcLPL.m.Lock()
	defer btcLPL.m.Unlock()
	if btcLPL.size>0{
		btcH.m.Lock()
		defer btcH.m.Unlock()
		btcH.size+=btcLPL.size
		btcH.txcs = append(btcH.txcs,btcLPL.txcs[0:btcLPL.size-1]...)
		btcLPL.size=0
		btcLPL.txcs=make([]*txcache ,localPoolCount,localPoolCount)
		btcLPL.deadline = btcLPL.deadline.Add(btcTimeM)
	}
}
//往local池写数据
//若local池已满，则写到global池
//local池的直接写到history池执行tx交易并监听
func (*btcLocalPool) Push(addrFrom, addrTo string, transfer, fee coins.CoinAmounter,txrchan chan<- *TxResult){
	btcLPL.m.Lock()
	defer btcLPL.m.Unlock()
	tc := &txcache{
		birthday: time.Now(),
		addf:     addrFrom,
		addt:     addrTo,
		transfer: transfer,
		fee:      fee,
		txchan:  txrchan,
	}
	//全局+1
	if btcLPL.size >= localPoolCount {
		 btcGPL.m.Lock()
		 defer  btcGPL.m.Unlock()

		 btcGPL.size+=1
		 btcGPL.txcs = append(btcGPL.txcs,tc)
	} else {
		//本地+1
		btcLPL.size+=1
		tcing := &txexcuting{
			txcache:tc,
		}
		//历史+1
		btcHPL.m.Lock()
		defer btcHPL.m.Unlock()
		btcHPL.size+=1
		btcHPL.txcsing = append(btcHPL.txcsing,tcing)
	}
}
func (g *btcGlobalPool) Push(addrFrom, addrTo string, transfer, fee coins.CoinAmounter) {

}
func (g *btcGlobalPool) Pull() {

}
