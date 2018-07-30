package wallets

import (
	"github.com/Rennbon/blockchainDemo/coins"
	"sync"
	"time"
	"btcd/chaincfg/chainhash"
)

//本地容器上限
const localPoolCount = int(10)

//全局容器结构体
type btcLocalPool struct {
	m         sync.RWMutex
	deadline  time.Time  //死亡时间
	size      int        //数量
	txcs []*txcache //tx组
}
type btcGlobalPool struct {
	m         sync.RWMutex
	size int        //数量
	txcs []*txcache //链表
}
type btcHistoryPool struct {
	m sync.RWMutex
	size int //总数
	txcsing []* txexcuting//处理的交易
	offset int //当前处理位置，指提交交易并广播
}

type btcgp btcGlobalPool
type btclp btcLocalPool
type btcHistory globalTxCachePool

type txexcuting struct {
	txHash *chainhash.Hash //
	state bool
	txcache
}
type txcache struct {
	birthday time.Time          //出生时间
	addf     string             //来自地址
	addt     string             //去向地址
	transfer coins.CoinAmounter //交易金额
	fee      coins.CoinAmounter //交易小费
}

var btcG btcGlobal
var btcL btcLocal
var
var btcH btcHistory
var btcTimeM,_ = time.ParseDuration("10m")

func (l *btcLocal) Restart() {
	l.m.Lock()
	defer l.m.Unlock()
	if l.size>0{
		btcH.m.Lock()
		defer btcH.m.Unlock()
		btcH.size+=l.size
		btcH.txcs = append(btcH.txcs,l.txcs[0:l.size-1]...)
		l.size=0
		l.txcs=make([]*txcache ,localPoolCount,localPoolCount)
		l.deadline = l.deadline.Add(btcTimeM)
	}
}
func (l *btcLocal) Push(addrFrom, addrTo string, transfer, fee coins.CoinAmounter) {
	l.m.Lock()
	defer l.m.Unlock()
	if l.size >= localPoolCount {
		//todo 放到全局
	} else {
		tc := &txcache{
			birthday: time.Now(),
			addf:     addrFrom,
			addt:     addrTo,
			transfer: transfer,
			fee:      fee,
		}
		l. = append(l.txcs, tc)
	}

}
func (g *btcGlobal) Push(addrFrom, addrTo string, transfer, fee coins.CoinAmounter) {

}
func (g *btcGlobal) Pull() {

}
