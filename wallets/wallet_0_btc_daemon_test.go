package wallets

import (
	"bytes"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

///////////////////////////////////pool-test////////////////////////////////////////////////
//验证填充块高
func TestBtcDaemon_fillBlockHeight(t *testing.T) {
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

//验证实例化btcDaemon
func TestNewBTCDaemon(t *testing.T) {
	daemon := NewBTCDaemon(time.NewTicker(5 * time.Second))
	t.Log(daemon.blkHt)
}

//验证监听区块高度并重绘内存中缓存的全局唯一高度
func TestMonitoringBtcBlockHeight(t *testing.T) {
	daemon := NewBTCDaemon(time.NewTicker(1 * time.Second))
	height := daemon.blkHt
	t.Log("目前高度:", height)
	log.Println("目前高度:", height)
	for i := 0; i < 100; i++ {
		select {
		case <-daemon.tick.C:
			log.Println("第", i, "次")
			daemon.monitoringBtcBlockHeight()
			log.Println("current block height", daemon.blkHt)
			log.Println("old block height", height)
			if daemon.blkHt > height+2 {
				log.Println("I am coming!!!")
				return
			}
		}
	}
}

//测试push
//
//  1) 往btcDaemon中推入交易参数
//  2) local池达到限高后会转移到global池，local计数的参数会直接进入exch通道消费
func TestBtcDaemon_push(t *testing.T) {
	daemon := NewBTCDaemon(time.NewTicker(5 * time.Second))
	wg := new(sync.WaitGroup)
	count := 30
	wg.Add(count)
	go func() {
		for i := float64(0); i < float64(count); i++ {
			time.Sleep(5 * time.Second)
			transfer, _ := btcCoin.FloatToCoinAmout(i)
			fee, _ := btcCoin.StringToCoinAmout("0.0001")
			txrchan := make(chan *TxResult)
			daemon.push("输入地址", "输出地址", transfer, fee, txrchan)
			wg.Done()
		}
	}()
	//这里只打印localPoolCount个（写测试时订的是10）
	go func(d *btcDaemon) {
		i := 0
		for ch := range d.exch {
			i++
			log.Println("第", i, "个", ch.addf, ch.transfer.String())
		}
	}(daemon)

	go func(d *btcDaemon) {
		for {
			select {
			case <-d.tick.C:
				//从0开始累加，到localPoolCount个后不在递增，转移至global缓存
				log.Println("打印检测,local-", d.lpl.size, "-", d.lpl.deadline)
				//local未达到localPoolCount时，这里都为0 ，等d.lpl溢出（超过localPoolCount）后开始累加，无上限
				log.Println("打印检测,global-", d.gpl.size)
			}

		}
	}(daemon)
	wg.Wait()
}

//假设当前全局区块高度block height
//往历史池里面添加测试数据多组,target block height（tbh）, 测试数据中的tbh的分布在block height2端都需要有
//
//全局块高度高度++
//
//循环执行测试函数
func TestBtcDaemon_fillConfmQ(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	daemon := NewBTCDaemon(time.NewTicker(1 * time.Second))
	daemon.blkHt = 100 //全局高度
	//所有假设target block height [80,120]
	for i := 0; i < 100; i++ {
		tbh := 80 + rand.Int63n(41)
		t := &txexcuting{
			targetH: tbh,
			blockH:  int64(i), //借字段一用，字段含义非当前测试的含义
		}
		daemon.hpl.txcsing = append(daemon.hpl.txcsing, t)
		daemon.hpl.size++
		log.Println("第", i, "个 ", tbh)
	}
	size := 0
	//flag :=true
	for i := 0; i < 20; i++ {
		select {
		case <-daemon.tick.C:
			//if i%2==1{
			daemon.blkHt++
			//}
			daemon.fillConfmQ()
			size = len(daemon.cq.q)

			log.Println("第", i, "组")
			log.Println("size:", size)
			log.Println("当前块高:", daemon.blkHt)
			buff := &bytes.Buffer{}
			for k, v := range daemon.cq.q {
				if v.targetH > daemon.blkHt {
					t.Fail()
					t.Error("历史池中扒出的数据含有大于于当前块高的数据")
				}
				buff.WriteString("序号:")
				buff.WriteString(strconv.FormatInt(int64(k), 10))
				buff.WriteString("第")
				buff.WriteString(strconv.FormatInt(v.blockH, 10))
				buff.WriteString("个,")
				buff.WriteString("块高:")
				buff.WriteString(strconv.FormatInt(v.targetH, 10))
				buff.WriteString("\n")
			}
			if len(daemon.hpl.txcsing) != daemon.hpl.size {
				t.Fail()
				t.Error("历史池冗余计数和实际slice长度不匹配")
			}
			log.Println(buff.String())
			log.Println("历史池剩余size", daemon.hpl.size)
			buff2 := &bytes.Buffer{}
			for k, v := range daemon.hpl.txcsing {
				if v.targetH <= daemon.blkHt {
					t.Fail()
					t.Error("历史池被扒后还有小于当前块高的数据")
				}
				buff2.WriteString("序号:")
				buff2.WriteString(strconv.FormatInt(int64(k), 10))
				buff2.WriteString("第")
				buff2.WriteString(strconv.FormatInt(v.blockH, 10))
				buff2.WriteString("个,")
				buff2.WriteString("块高:")
				buff2.WriteString(strconv.FormatInt(v.targetH, 10))
				buff2.WriteString("\n")
			}
			log.Println(buff2.String())
		}
	}
}

func TestBtcDaemon_consumeeExcuCH(t *testing.T) {

	daemon := NewBTCDaemon(time.NewTicker(5 * time.Second))
	count := float64(1)
	go func(daemon *btcDaemon) {
		for i := float64(0); i < count; i++ {
			select {
			case <-daemon.tick.C:
				//exch 中填充数据
				txch := make(chan *TxResult)
				transfer, _ := btcCoin.FloatToCoinAmout(i)
				fee, _ := btcCoin.StringToCoinAmout("0.0001")
				tching := &txexcuting{
					txcache: &txcache{
						birthday: time.Now().Local(),
						addf:     "mhAfGecTPa9eZaaNkGJcV7fmUPFi3T2Ki8",
						addt:     "mmE4PtemXdgZY5wpiFYvhMS5hfjV4R1GCD",
						transfer: transfer,
						fee:      fee,
						txrchan:  txch,
					},
				}
				daemon.exch <- tching
			}
		}

	}((daemon))
	go daemon.consumeeExcuCH()

	go func(daemon *btcDaemon) {
		for i := float64(0); i < count; i++ {
			select {
			case <-daemon.tick.C:
				buff := &bytes.Buffer{}
				for k, v := range daemon.hpl.txcsing {
					buff.WriteString("序号:")
					buff.WriteString(strconv.FormatInt(int64(k), 10))
					buff.WriteString(" txId:")
					buff.WriteString(v.txHash.String())
					buff.WriteString(" 目标块高:")
					buff.WriteString(strconv.FormatInt(v.targetH, 10))
					buff.WriteString(" 提币地址:")
					buff.WriteString(v.addf)
					buff.WriteString(" 充值地址:")
					buff.WriteString(v.addt)
					buff.WriteString(" 交易金额：")
					buff.WriteString(v.transfer.String())
					buff.WriteString("\n")
				}
				log.Println(buff.String())
			}
		}
	}(daemon)
	time.Sleep(5 * 5 * time.Second)
}
