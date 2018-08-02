/*
发起交易的生命周期（大致流程图）：
					  ┏—————————————---—┓
					  |  传参调用交易方法 |
					  ┗-----------------┛
								↓
					  ┏—————————————---—┓
					  |      参数拦截	    |
					  ┗-----------------┛
								↓
						   ----------
						﹤   参数拦截	   >
						   ----------
				默认	↓  					  ↓  本地池已满
		  ┏—————————————---—┓    ┏—————————————---—┓
		  |    本地池计数     | <=|   缓存到全局池	   |    (如本地池上限是10个，本地池每隔一个周期清一次计数，清完后优先从全局池获取前10个进入本地计数)
		  ┗-----------------┛  	 ┗-----------------┛
					↓
		  ┏—————————————---—┓     ┏—————————————---—┓    ┏—————————————---—┓
		  |  推入tx处理队列   |  -> |   执行TX操作     | -> | 填充预期处理块高  |
		  ┗-----------------┛	  ┗-----------------┛	 ┗-----------------┛
																  ↓
		  ┏—————————————---—┓     ┏—————————————---—┓    ┏—————————————---—┓
		  |    处理完毕记录   |  -> | 监听块高拉取处理  | <- |    放入历史池    |
		  ┗-----------------┛	  ┗-----------------┛	 ┗-----------------┛
生命周期


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

	//blockHeight    int64
	confirmNum     = int32(6)
	localPoolCount = int(10) //本地容器上限
	btcTimeM, _ = time.ParseDuration("10m")       //同上延迟localpool deadline用
	//btcD = NewBTCDaemon(time.NewTicker(10*time.Minute))
/*	btcGPL      = &btcGlobalPool{}                //全局池
	btcLPL      = &btcLocalPool{}                 //本地池
	btcHPL      = &btcHistoryPool{}               //历史池,交易处理等待验证的
	tick        = time.NewTicker(5 * time.Second) //扫描计区块等周期计时器

	historyWG   = new(sync.WaitGroup)             //监听历史池

	cq = &confmQ{
		q: make([]*txexcuting, 0, localPoolCount*3),
		m: new(sync.Mutex),
	}
	//前后预留2批，共3批，不进块的情况下这个会上升，当恢复后需要释放
	excuCH = make(chan *txexcuting, localPoolCount*2) //等到处理的通道，理论上也是单位时间段最多10个左右，最多同一周期并发出现20个，所以cap设置20就够了*/
)

type btcDaemon struct {
	tick  *time.Ticker  //周期计时器
	blkHt int64 //btc块高
	gpl *btcGlobalPool  //全局池
	lpl *btcLocalPool //本地池
	hpl *btcHistoryPool //历史池（等待确认tx状态的）
	cq *confmQ  //当前可以确认tx状态的，指block height now > 6块 + created block height
	exch chan *txexcuting //需要处理提交到共链的tx，该chan下数据需要提交离线签名=>处理块高=>投入历史池

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

func (d *btcDaemon)Run(){
	if d==nil{
		panic("老兄，不存在的，你需要调用NewBTCDaemon创建一个对象")
	}
	//消费，永动机，全局d即可
	go func(){
		d.consumeeExcuCH()
	}()
	for {
		select {
		case <-d.tick.C:
			//刷区块高度
			go d.monitoringBtcBlockHeight()

			//重置单位时间本地限流
			go func(d *btcDaemon){
				time.Sleep(1*time.Minute)
				d.restart()
			}(d)
			//从历史池抓取即将处理的数据到待验证队列等到消费，扔到cq队列
			go func(d *btcDaemon){
				time.Sleep(2*time.Minute)
				d.fillConfmQ()
			}(d)

			// 轮询cq队列，监听公链状态同步
			go func(d *btcDaemon){
				time.Sleep(3*time.Minute)
				d.listenMainNet()
			}(d)
		}
	}
}




//btc后端运行机制，单例跑
func NewBTCDaemon(tick *time.Ticker )(daemon *btcDaemon){
	var err error
	d := &btcDaemon{}
	//1
	d.tick = tick// time.NewTicker(5*time.Second)
	//2
	d.blkHt,err = btcClient.GetBlockCount()
	if err!=nil{
		panic(err)
	}
	//3
	d.cq = &confmQ{
		q: make([]*txexcuting, 0, localPoolCount*3),
		m: new(sync.Mutex),
	}
	//4 处理的通道，理论上也是单位时间段最多10个左右，最多同一周期并发出现20个，所以cap设置20就够了
	d.exch = make(chan *txexcuting, localPoolCount*2)
	//5
	d.gpl = &btcGlobalPool{
		m: new(sync.Mutex),
		size:0,
		txcs:[]*txcache{},
	}
	tm, _ := time.ParseDuration("10m")
	//6
	d.lpl = &btcLocalPool{
		m: new(sync.Mutex),
		size:0,
		deadline:time.Now().Add(tm),//当前时间+10分钟
	}
	//7
	d.hpl = &btcHistoryPool{
		m: new(sync.Mutex),
		size:0,
		txcsing:[]*txexcuting{},
	}
	return  d
}


//（被动方法，触发器触发,需要间隔时间段触发）
//
//填充需要当前时间检测的tx的公链状态
func (d *btcDaemon)fillConfmQ() {
	d.hpl.m.Lock()
	defer d.hpl.m.Unlock()
	{ //锁池
		txhasharr := []*chainhash.Hash{}
		qrm := []int{}
		d.cq.m.Lock()
		defer d.cq.m.Unlock()
		{
			for k, v := range d.hpl.txcsing {
				if v.targetH <= d.blkHt {

					txhasharr = append(txhasharr, v.txHash)
					d.cq.q = append(d.cq.q, v)
					qrm = append(qrm, k)
				}
			}
		}
		//移除老数据
		if len(qrm) > 0 {
			d.hpl.txcsing = RmoveSliceByIndex(d.hpl.txcsing, qrm)
			d.hpl.size = d.hpl.size-len(qrm)
		}
	}
}
//（被动方法，触发器触发,需要间隔时间段触发）
//监听公链
//todo 	轮询tx,检测到ok的需要close 内部chan状态
func (d *btcDaemon)listenMainNet() {
	d.cq.m.Lock()
	defer d.cq.m.Unlock()
	{
		qneed := []int{}
		for k, cur := range d.cq.q {
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
				d.cq.q[k].txrchan <- tr
				close(d.cq.q[k].txrchan)
			} else {
				qneed = append(qneed)
			}
		}
		//这里可以封装成方法
		if len(qneed) == 0 {
			d.cq.q = make([]*txexcuting, 0, localPoolCount*3)
		} else {
			d.cq.q = GetSliceByIndex(d.cq.q, qneed)
		}
	}
}
//（被动方法,一次触发，跑到死）
//启动一次，跑个没完
//todo 消费处理池，执行TX,需要轮询
func (d *btcDaemon)consumeeExcuCH() {
	for ch := range d.exch  {
		//todo 执行tx并广播到共链，(需要分离SendAddressToAddress)
		//todo 广播成功后推入历史池监听
		btcSer := &BtcService{}
		txid, err := btcSer.sendAddressToAddress(ch.addf, ch.addt, ch.transfer, ch.fee)
		if err != nil {
			//TODO 日志
			log.Println(err)
			 txr := &TxResult{
				 Err:err,
			 }
			 ch.txrchan<-txr
			 close(ch.txrchan)
		} else {
			//填充块高度
			txe := &txexcuting{}
			txe.fillBlockHeight()
			txe.txHash, _ = chainhash.NewHashFromStr(txid)
			// 扔给历史池
			d.hpl.m.Lock()
			{
				d.hpl.size++
				d.hpl.txcsing = append(d.hpl.txcsing, txe)
			}
			d.hpl.m.Unlock()
		}
	}
}



//（被动方法，触发器触发,需要间隔时间段触发）
// 清空local池，首先去全局同步到本地
func (d *btcDaemon)restart() {
	d.lpl.m.Lock()
	defer d.lpl.m.Unlock()
	{ //本地锁池
		if d.lpl.size > 0 {
			d.lpl.size = 0
		}
		d.lpl.deadline.Add(btcTimeM)
		d.gpl.m.Lock()
		defer d.gpl.m.Unlock()
		{ //全局锁池
			if d.gpl.size > 0 {

				size := 0
				if d.gpl.size < localPoolCount {
					size = d.gpl.size
				} else {
					size = localPoolCount
				}
				d.lpl.size = size

				//处理chan+size
				for _, v := range d.gpl.txcs {
				d.exch	 <- &txexcuting{
						txcache: v,
					}
				}
				//全局缩减
				d.gpl.size -= size
				d.gpl.txcs = d.gpl.txcs[size:]
			}
		}
	}
}

//(主动方法，交易发起的时候调用)
//往local池写数据
//若local池已满，则写到global池
//local池的直接写到history池执行tx交易并监听
func (d *btcDaemon) push(addrFrom, addrTo string, transfer, fee coins.CoinAmounter, txrchan chan<- *TxResult) {
	d.lpl.m.Lock()
	defer d.lpl.m.Unlock()
	{ //本地锁池
		tc := &txcache{
			birthday: time.Now(),
			addf:     addrFrom,
			addt:     addrTo,
			transfer: transfer,
			fee:      fee,
			txrchan:  txrchan,
		}
		if d.lpl.size >= localPoolCount {
			d.gpl.m.Lock()
			defer d.gpl.m.Unlock()
			{ //全局锁池
				d.gpl.size += 1
				d.gpl.txcs = append(d.gpl.txcs, tc)
			}
		} else {
			//本地+1
			d.lpl.size += 1
			tcing := &txexcuting{
				txcache: tc,
			}
			//处理chan+1
			d.exch <- tcing
		}
	}
}

//OK （被动方法，触发器触发）
//监听区块高度，需要放入到Init函数
func (d *btcDaemon)monitoringBtcBlockHeight() {
	height, err := btcClient.GetBlockCount()
	if err != nil {
		log.Println(err)
	} else {
		if height > d.blkHt  {
			d.blkHt = height
		}
	}
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
	lenI := len(indes)
	lenS := len(source)
	qnew := make([]*txexcuting, 0, lenS-lenI)
	if lenI > 0 {
		mdl := 0
		for _, v := range indes {
			if v == 0 {
				mdl = 0
			} else if mdl != v {

				qnew = append(qnew, source[mdl:v]...)
			}
			mdl = v + 1
		}
		if mdl < lenS {
			qnew = append(qnew, source[mdl:]...)
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
