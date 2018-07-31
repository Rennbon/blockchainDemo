/*
维护本地队列，以及全局，历史队列
 * 所有请求参数都需缓存下来，带上时间，方便追踪，和挂死的重启
 - 本地队列限制单位时间只能允许定量的tx提交到公链
 - 全局队列用来维护，超出本地队列长度限制后的将要执行tx交易的参数
 - 进本地队里计数的直接通过exec通道处理tx广播到共链，成的直接推入历史队里，失败的日志
 - 历史队里通过区块高度兼容tx状态

1. 需要轮询读取block height
2. 需要给tx赋值当前区块高度及验证区块高度
3. 历史池需要遍历获取池中数据的状态并维护，每隔一定时间清除状态已确认的数据

流程图:http://note.youdao.com/noteshare?id=ecb974714a77656cd8343d9282bcfbf6
btc 一区块1MB，10分钟一区块，一区块大致存储1000比交易
*/
package wallets

import (
	"github.com/Rennbon/blockchainDemo/coins"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"log"
	"sync"
	"time"
)

var (

	blockHeight    int64
	confirmNum     = int32(6)
	localPoolCount = int(10) //本地容器上限

	btcGPL      = &btcGlobalPool{}                //全局池
	btcLPL      = &btcLocalPool{}                 //本地池
	btcHPL      = &btcHistoryPool{}               //历史池,交易处理等待验证的
	tick        = time.NewTicker(5 * time.Second) //扫描计区块等周期计时器
	btcTimeM, _ = time.ParseDuration("10m")       //同上延迟localpool deadline用
	historyWG   = new(sync.WaitGroup)             //监听历史池

	cq = &confmQ{
		q: make([]*txexcuting, 0, localPoolCount*3),
		m: new(sync.Mutex),
	}
	//前后预留2批，共3批，不进块的情况下这个会上升，当恢复后需要释放
	excuCH = make(chan *txexcuting, localPoolCount*2) //等到处理的通道，理论上也是单位时间段最多10个左右，最多同一周期并发出现20个，所以cap设置20就够了


)

type btcDaemon struct {
	tick  *time.Ticker  //周期计时器
	blkHt int64 //btc块高
	gpl *btcGlobalPool
	lpl *btcLocalPool
	hpl *btcHistoryPool
	cq *confmQ
	exch chan *txexcuting

}

//小于当前块的需要验证的tx池
type confmQ struct {
	q []*txexcuting
	m *sync.Mutex
}

//local池，
type btcLocalPool struct {
	m        *sync.Mutex
	deadline time.Time //死亡时间
	size     int       //数量
}

//全局池，只维护计数
type btcGlobalPool struct {
	m    *sync.Mutex
	size int        //数量
	txcs []*txcache //tx组
}

//历史池，这里才开始执行增伤查
type btcHistoryPool struct {
	m       *sync.Mutex
	size    int           //总数
	txcsing []*txexcuting //处理的交易
}

type txexcuting struct {
	txHash  *chainhash.Hash //txhash
	blockH  int64           //所在的区块高度
	targetH int64           //目标区块高度
	status  bool            //状态
	excount int8            //执行过的次数，超过一定次数，执行plan X（移出并日志）
	*txcache
}
type txcache struct {
	birthday time.Time          //出生时间
	addf     string             //来自地址
	addt     string             //去向地址
	transfer coins.CoinAmounter //交易金额
	fee      coins.CoinAmounter //交易小费
	txrchan  chan<- *TxResult   //chan维持
}
func NewBTCDaemon()(daemon *btcDaemon,err error){
	d := &btcDaemon{}
	d.tick = time.NewTicker(5*time.Second)
	d.blkHt,err = btcClient.GetBlockCount()
	if err!=nil{
		panic(err)
	}
	d.cq = &confmQ{
		q: make([]*txexcuting, 0, localPoolCount*3),
		m: new(sync.Mutex),
	}
	d.exch = make(chan *txexcuting, localPoolCount*2)
	d.gpl = &btcGlobalPool{
		m: new(sync.Mutex),
		size:0,
		txcs:[]*txcache{},
	}
	tm, _ = time.ParseDuration("10m")
	d.lpl = &btcLocalPool{
		m: new(sync.Mutex),
		size:0,
		deadline:time.Now().Add(tm),//当前时间+10分钟
	}
	d.hpl = &btcHistoryPool{
		m: new(sync.Mutex),
		size:0,
		txcsing:[]*txexcuting{},
	}


}


//todo 消费处理池，执行TX
func (d *btcDaemon)consumeeExcuCH() {
	tick4Excu := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-tick4Excu.C:
			for ch := range excuCH {
				//todo 执行tx并广播到共链，(需要分离SendAddressToAddress)
				//todo 广播成功后推入历史池监听
				btcSer := &BtcService{}
				txid, err := btcSer.SendAddressToAddress(ch.addf, ch.addt, ch.transfer, ch.fee)
				if err != nil {
					//TODO 日志
				} else {
					//填充块高度
					txe := &txexcuting{}
					txe.fillBlockHeight()
					txe.txHash, _ = chainhash.NewHashFromStr(txid)
					//todo 扔给历史池
					btcHPL.m.Lock()
					defer btcHPL.m.Unlock()
					{
						btcHPL.size++
						btcHPL.txcsing = append(btcHPL.txcsing, txe)
					}
				}
			}
		}
	}
}

//填充需要当前时间检测的tx的公链状态
func (d *btcDaemon)fillConfmQ() {
	for {
		select {
		case <-tick.C:
			btcHPL.m.Lock()
			txhasharr := []*chainhash.Hash{}
			defer btcHPL.m.Unlock()
			{ //锁池
				qrm := []int{}
				for k, v := range btcHPL.txcsing {
					if v.targetH <= blockHeight && !v.status {
						cq.m.Lock()
						defer cq.m.Unlock()
						{
							txhasharr = append(txhasharr, v.txHash)
							cq.q = append(cq.q, v)
							qrm = append(qrm, k)
						}
					}
				}
				//移除老数据
				if len(qrm) > 0 {
					btcHPL.txcsing = RmoveSliceByIndex(btcHPL.txcsing, qrm)
				}
			}
		}
	}
}

//监听公链
//todo 	轮询tx,检测到ok的需要close 内部chan状态
func (d *btcDaemon)listenMainNet() {
	cq.m.Lock()
	defer cq.m.Unlock()
	{
		qneed := []int{}
		for k, cur := range cq.q {
			txInfo, err := btcClient.GetRawTransactionVerbose(cur.txHash)
			if err != nil {
				log.Println(err)
				continue
			}
			if txInfo.Confirmations >= uint64(confirmNum) {
				tr := &TxResult{
					TxId:   cur.txHash.String(),
					Status: true,
					Err:    nil,
				}
				for _, v := range txInfo.Vout {
					amout, _ := btcCoin.FloatToCoinAmout(v.Value)
					tai := &TxAddressInfo{
						Address: v.ScriptPubKey.Addresses,
						Amount:  amout,
					}
					tr.AddInfos = append(tr.AddInfos, tai)
				}
				cq.q[k].txrchan <- tr
				close(cq.q[k].txrchan)
			} else {
				qneed = append(qneed)
			}
		}
		//这里可以封装成方法
		if len(qneed) == 0 {
			cq.q = make([]*txexcuting, 0, localPoolCount*3)
		} else {
			cq.q = GetSliceByIndex(cq.q, qneed)
		}
	}
}

//todo 计时器每满一次，清空local池，首先去全局同步到本地
func (d *btcDaemon)restart() {
	for {
		select {
		case <-tick.C:
			btcLPL.m.Lock()
			defer btcLPL.m.Unlock()
			{ //本地锁池
				if btcLPL.size > 0 {
					btcLPL.size = 0
				}
				btcLPL.deadline.Add(btcTimeM)

				btcGPL.m.Lock()
				defer btcGPL.m.Unlock()
				{ //全局锁池
					if btcGPL.size > 0 {

						size := 0
						if btcGPL.size < localPoolCount {
							size = btcGPL.size
						} else {
							size = localPoolCount
						}
						btcLPL.size = size

						//处理chan+size
						for _, v := range btcGPL.txcs {
							excuCH <- &txexcuting{
								txcache: v,
							}
						}
						//全局缩减
						btcGPL.size -= size
						btcGPL.txcs = btcGPL.txcs[size:]
					}
				}
			}
		}
	}
}

//往local池写数据
//若local池已满，则写到global池
//local池的直接写到history池执行tx交易并监听
func (d *btcDaemon) push(addrFrom, addrTo string, transfer, fee coins.CoinAmounter, txrchan chan<- *TxResult) {
	btcLPL.m.Lock()
	defer btcLPL.m.Unlock()
	{ //本地锁池
		tc := &txcache{
			birthday: time.Now(),
			addf:     addrFrom,
			addt:     addrTo,
			transfer: transfer,
			fee:      fee,
			txrchan:  txrchan,
		}
		if btcLPL.size >= localPoolCount {
			btcGPL.m.Lock()
			defer btcGPL.m.Unlock()
			{ //全局锁池
				btcGPL.size += 1
				btcGPL.txcs = append(btcGPL.txcs, tc)
			}
		} else {
			//本地+1
			btcLPL.size += 1
			tcing := &txexcuting{
				txcache: tc,
			}
			//处理chan+1
			excuCH <- tcing
		}
	}
}

//OK
//监听区块高度，需要放入到Init函数
func (d *btcDaemon)monitoringBtcBlockHeight() {
	for {
		select {
		case <-tick.C:
			height, err := btcClient.GetBlockCount()
			if err != nil {
				log.Println(err)
			} else {
				if height > blockHeight {
					blockHeight = height
				}
			}
		}
	}
}

//todo 执行交易SendAddressToAddress
func (ex *txexcuting) excuteTx() {

}

//ok
//处理tx交易信息
//通过txHash获取tx详情，然后通过tx详情中的blockHash获取当前tx所在的block高度
//将txHash 和 block高度，以及确认的block高度推入tb4check channel
func (ex *txexcuting) fillBlockHeight() {

	txinfo, err := btcClient.GetTransaction(ex.txHash)
	if err != nil {
		log.Printf("txId:%s 获取tx详情失败\r\n", ex.txHash.String())
		return
	} else {
		blockHash, err := chainhash.NewHashFromStr(txinfo.BlockHash)
		if err != nil {
			log.Printf("blockHash:%s string to hash失败\r\n", txinfo.BlockHash)
			return
		}
		blockInfo, err := btcClient.GetBlockHeaderVerbose(blockHash)
		if err != nil {
			log.Printf("txId:%s 获取block详情失败\r\n", txinfo.BlockHash)
			return
		}
		ex.blockH = int64(blockInfo.Height)
		ex.targetH = int64(blockInfo.Height + confirmNum)
	}
}

func RmoveSliceByIndex(source []*txexcuting, indes []int) []*txexcuting {
	qnew := make([]*txexcuting, 0, len(source)-len(indes))
	if len(indes) > 0 {
		mdl := 0
		for _, v := range indes {
			if v == 0 {
				mdl = 0
			} else {
				qnew = append(qnew, source[mdl:v]...)
				mdl = v + 1
			}
		}
	}
	return qnew
}
func GetSliceByIndex(source []*txexcuting, indes []int) []*txexcuting {
	qnew := make([]*txexcuting, 0, len(indes))
	for k, v := range source {
		for _, qv := range indes {
			if qv == k {
				qnew = append(qnew, v)
			} else if qv > k {
				break
			}
		}
	}
	return qnew
}
