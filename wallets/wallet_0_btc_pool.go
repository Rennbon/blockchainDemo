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

	/*//确认相关

	tb4check     chan *txblcok
	txHash4check chan *chainhash.Hash //txId
	*/
	blockHeight    int64
	confirmNum     = int32(6)
	localPoolCount = int(10) //本地容器上限

	btcGPL      = &btcGlobalPool{}                //全局池
	btcLPL      = &btcLocalPool{}                 //本地池
	btcHPL      = &btcHistoryPool{}               //历史池,交易处理等待验证的
	tick        = time.NewTicker(5 * time.Second) //扫描计区块等周期计时器
	btcTimeM, _ = time.ParseDuration("10m")       //同上延迟localpool deadline用
	historyWG   = new(sync.WaitGroup)             //监听历史池
	excuPool    = make(chan *txexcuting, 20)      //等到处理的通道，理论上也是单位时间段最多10个左右，最多同一周期并发出现20个，所以cap设置20就够了
)

//local池，
type btcLocalPool struct {
	m        sync.RWMutex
	deadline time.Time //死亡时间
	size     int       //数量
}

//全局池，只维护计数
type btcGlobalPool struct {
	m    sync.RWMutex
	size int        //数量
	txcs []*txcache //tx组
}

//历史池，这里才开始执行增伤查
type btcHistoryPool struct {
	m       sync.RWMutex
	size    int           //总数
	txcsing []*txexcuting //处理的交易
	offset  int           //当前处理位置，指提交交易并广播
}

type txexcuting struct {
	txHash  *chainhash.Hash //txhash
	blockH  int32           //所在的区块高度
	targetH int32           //目标区块高度
	status  bool            //状态
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

type txblcok struct {
	txHash   *chainhash.Hash
	blockNum int64 //创建时区块高度
	TargetBN int64 //需要验证的区块高度
}

//todo 消费处理池，执行TX
func consumeeExcuPool() {
	tick4Excu := time.NewTicker(5 * time.Second)
	for {
		select {
		case tick4Excu.C:
			for _ := range excuPool {
				//todo 执行tx并广播到共链，(需要分离SendAddressToAddress)
				//todo 广播成功后推入历史池监听
			}
		}
	}
}

//todo 监听历史池，处理交易状态
func consumeHistoryPool() {
	//监听历史池，有更新直接消费

}

//todo 计时器每满一次，清空local池，首先去全局同步到本地
func (*btcLocalPool) Restart() {
	for {
		select {
		case tick.C:
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
							excuPool <- &txexcuting{
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
func (*btcLocalPool) Push(addrFrom, addrTo string, transfer, fee coins.CoinAmounter, txrchan chan<- *TxResult) {
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
			excuPool <- tcing
		}
	}
}

//监听区块高度，需要放入到Init函数
func monitoringBtcBlockHeight() {
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

func pushTxResultIntoBtcTxRet() {
	for tb := range tb4check {
	JUSTDOIT:
		if tb.TargetBN >= blockHeight {
			txInfo, err := btcClient.GetRawTransactionVerbose(tb.txHash)
			if err != nil {
				log.Println(err)
				continue
			}
			if txInfo.Confirmations >= uint64(confirmNum) {
				tr := &TxResult{
					TxId:   tb.txHash.String(),
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
				btcTxRet <- tr
			}
		} else {
			//延时操作
			tick := time.NewTicker(5 * time.Second)
			for {
				select {
				case <-tick.C:
					goto JUSTDOIT
				}
			}
		}
	}
}

//todo 执行交易SendAddressToAddress
func (ex *txexcuting) excuteTx() {

}

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
		ex.blockH = blockInfo.Height
		ex.targetH = blockInfo.Height + confirmNum
	}
}
