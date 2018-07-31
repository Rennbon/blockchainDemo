/*
维护本地队列，以及全局，历史队列
 - 本地队列限制单位时间只能允许定量的tx提交到公链
 - 全局队列用来维护，超出本地队列长度限制后的将要执行tx交易的参数
 - 历史队列为发起tx交易，并检测tx状态用

1. 需要轮询读取block height
2. 需要给tx赋值当前区块高度及验证区块高度
3. 历史池需要遍历获取池中数据的状态并维护，每隔一定时间清除状态已确认的数据


btc 一区块1MB，10分钟一区块，一区块大致存储1000比交易
*/
package wallets

import (
	"btcd/chaincfg/chainhash"
	"github.com/Rennbon/blockchainDemo/coins"
	"log"
	"sync"
	"time"
)

//本地容器上限
const localPoolCount = int(10)

var (
	blockHeight int64

	btcGPL      *btcGlobalPool                    //全局池
	btcLPL      *btcLocalPool                     //本地池
	btcHPL      *btcHistoryPool                   //历史池
	tick        = time.NewTicker(5 * time.Second) //扫描计区块等周期计时器
	btcTimeM, _ = time.ParseDuration("10m")       //同上延迟localpool deadline用
	historyWG   *sync.WaitGroup                   //监听历史池
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
	txHash *chainhash.Hash //
	status bool            //状态
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

//计时器每满一次，清空local池，移动到history池
func (*btcLocalPool) Restart() {
	for {
		select {
		case tick.C:
			//本地池锁
			btcLPL.m.Lock()
			defer btcLPL.m.Unlock()
			if btcLPL.size > 0 {
				btcLPL.size = 0
			}
			btcLPL.deadline.Add(btcTimeM)

			//全局池锁
			btcGPL.m.Lock()
			defer btcGPL.m.Unlock()
			if btcGPL.size > 0 {

				btcHPL.m.Lock()
				defer btcHPL.m.Unlock()
				size := 0
				if btcGPL.size < localPoolCount {
					size = btcGPL.size
				} else {
					size = localPoolCount
				}
				btcLPL.size = size

				//历史池
				btcHPL.size += size
				for _, v := range btcGPL.txcs {
					btcHPL.txcsing = append(btcHPL.txcsing, &txexcuting{
						txcache: v,
					})
				}
				btcGPL.txcs = btcGPL.txcs[size:]
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
	tc := &txcache{
		birthday: time.Now(),
		addf:     addrFrom,
		addt:     addrTo,
		transfer: transfer,
		fee:      fee,
		txrchan:  txrchan,
	}
	//全局+1
	if btcLPL.size >= localPoolCount {
		btcGPL.m.Lock()
		defer btcGPL.m.Unlock()

		btcGPL.size += 1
		btcGPL.txcs = append(btcGPL.txcs, tc)
	} else {
		//本地+1
		btcLPL.size += 1
		tcing := &txexcuting{
			txcache: tc,
		}
		//历史+1
		btcHPL.m.Lock()
		defer btcHPL.m.Unlock()
		btcHPL.size += 1
		btcHPL.txcsing = append(btcHPL.txcsing, tcing)
	}
}

func consumeHistoryPool() {
	//监听历史池，有更新直接消费

}

//监听区块高度
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

//处理tx交易信息
//通过txHash获取tx详情，然后通过tx详情中的blockHash获取当前tx所在的block高度
//将txHash 和 block高度，以及确认的block高度推入tb4check channel
func btcExcuteTxHash() {
	for txHash := range txHash4check {
		go func(txHash *chainhash.Hash) {
			txinfo, err := btcClient.GetTransaction(txHash)
			if err != nil {
				log.Printf("txId:%s 获取tx详情失败\r\n", txHash.String())
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
				tb := &txblcok{
					txHash:   txHash,
					blockNum: int64(blockInfo.Height),
					TargetBN: int64(blockInfo.Height + confirmNum),
				}
				tb4check <- tb
			}
		}(txHash)
	}
}
