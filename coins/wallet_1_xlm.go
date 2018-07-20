package coins

import (
	"blockchainDemo/cert"
	"blockchainDemo/database"
	"blockchainDemo/errors"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"blockchainDemo/config"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"log"
)

////////////////////基础设施//////////////////////////
type XlmService struct {
}

//装载btc配置
func initXlmClinet(conf *config.XlmConf) {

	switch conf.Env {
	case config.None:
		panic("Please set the btc env in config.yml!")
		break
	case config.Net:
		client = horizon.DefaultPublicNetClient
		break
	case config.TestNet:
	case config.Regtest:
		client = horizon.DefaultTestNetClient
		break
	default:
		panic("Please set the btc env in config.yml!")
		break
	}
	log.Println("coins=>xlm_wallet=>initClinet sccuess.")
}

var (
	certXlmSrv    cert.XlmCertService
	xlmSrv        XlmService
	client        *horizon.Client
	baseReserve   float64 = 0.5    //账户保证金基数
	baseFee       float64 = 0.0001 //小费基数（单位:xlm）
	baseFeeLemuns uint64  = 100    //小费 (单位：lumens)

)

//////////////////////////////////////////////////////
//生成新账号
//
func (*XlmService) GetNewAddress(account string, mode AcountRunMode) (address, accountOut string, err error) {
	key, err := certXlmSrv.GenerateSimpleKey()
	if err != nil {
		return
	}

	//----------测试网络源账号创建--------start------------
	/*resp, err := http.Get("https://friendbot.stellar.org/?addr=" + key.Address)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}*/
	//----------测试网络源账号创建--------end------------
	//----------源账号创建--------start------------
	//创建新账户必须从已有资金的账户转账生成，所以理论上生产环境要用，必须要提供一个有钱的账户作为God来来创造一切
	godSeed := "SAACHR2TWFAJKLLLC5TEYTSYPXA7AIBM6A2KZ7MQ4XEYRJEZFNOR6VOC"
	godAddress := "GBZKTZBJIMLFPUGZUNCUTJCUUREEG4W4UF74K5DRJRZISQNYQP3QOUYX"

	//源账号不要开通其他付费条目
	//源账户保底剩余 基础保证金*2
	//新账户保底创建 基础保证金*2
	comparedAmount := baseReserve*2 + baseFee + baseReserve*2
	if err = checkBalanceEnough(godAddress, comparedAmount); err != nil {
		return
	}
	//获取序列数
	num, err := client.SequenceForAccount(godAddress)
	if err != nil {
		return
	}
	/*
		Trustlines
		Offers
		Signers    新用户初始化会有一条singer
		Data entries
	*/
	amount := baseReserve * 2 //基础+Singer=2条
	amountStr := strconv.FormatFloat(amount, 'f', 8, 64)
	tx, err := build.Transaction(
		build.TestNetwork,
		build.Sequence{uint64(num) + 1}, //这里用autoSequence 失败了，公链可以在尝试下
		build.SourceAccount{godSeed},
		build.MemoText{"Create Account"}, //元数据，就是便签
		build.CreateAccount(
			build.Destination{key.Address},
			build.NativeAmount{amountStr}, //初始账号最小为0.5Lumens
		),
		build.BaseFee{baseFeeLemuns},
	)
	if err != nil {
		return
	}
	txe, err := tx.Sign(godSeed) //画押
	if err != nil {
		return
	}
	txeB64, err := txe.Base64()
	if err != nil {
		return
	}
	_, err = client.SubmitTransaction(txeB64) //提交tx
	if err != nil {
		return
	}
	//----------源账号创建--------end------------
	err = dhSrv.AddAccount(account, key.PrivKey, key.PubKey, key.Address, key.Seed, database.XLM)
	if err != nil {
		return
	}
	return key.Address, account, nil
}
func (*XlmService) GetBalanceInAddress(address string) (balance float64, err error) {
	account, err := client.LoadAccount(address)
	if err != nil {
		return 0, err
	}
	bls := float64(0)
	for _, v := range account.Balances {
		fmt.Println(v.Balance)
		curBls, err := strconv.ParseFloat(v.Balance, 64)
		if err != nil {
			return 0, err
		}
		bls += curBls
	}
	return bls, nil
}

//转账
//addrForm来源地址，addrTo去向地址
//transfer 转账金额
//fee 小费
func (*XlmService) SendAddressToAddress(addrFrom, addrTo string, transfer, fee float64) (txId string, err error) {
	//数据库获取prv pub key等信息，便于调试--------START------
	actf, err := dhSrv.GetAccountByAddress(addrFrom)
	if err != nil {
		return
	}
	//----------------------------------------END-----------
	//验证地址是否有效
	if _, err = client.LoadAccount(addrTo); err != nil {
		return
	}
	//100 stroops (0.00001 XLM).
	//The base fee (currently 100 stroops) is used in transaction fees.
	//sumfee = num of operations × base fee
	//The base reserve (currently 0.5 XLM) is used in minimum account balances.
	//(2 + n) × base reserve = 2.5 XLM.
	amount := strconv.FormatFloat(transfer, 'f', 8, 64)
	//验证金额总数
	comparedAmount := transfer + baseFee + baseReserve*2*2
	if err = checkBalanceEnough(addrFrom, comparedAmount); err != nil {
		return
	}
	//小费是自己扣的，不需要这边实现，金额总数也不需要验证，当然可以验证
	tx, err := build.Transaction(
		build.TestNetwork,
		build.SourceAccount{addrFrom}, //lumens（代币名称）当前主人的地址
		build.AutoSequence{client},    //sequence序列号自动
		build.MemoText{"Just do it"},  //元数据，就是便签
		build.Payment(
			build.Destination{addrTo},  // lumens（代币名称）下个主人的地址
			build.NativeAmount{amount}, //官方payments用string主要防止精度丢失
		),
		//build.BaseFee{baseFeeLemuns},//小费不写也会扣，只要钱够
	)

	if err != nil {
		return
	}
	// Sign the transaction to prove you are actually the person sending it.
	txe, err := tx.Sign(actf.Seed) //签名需要用seed
	if err != nil {
		return
	}

	txeB64, err := txe.Base64()
	if err != nil {
		return
	}

	// And finally, send it off to Stellar!
	resp, err := client.SubmitTransaction(txeB64) //提交tx
	if err != nil {
		return
	}
	//存储到数据库，方便检验
	dhSrv.AddTx(resp.Hash, addrFrom, []string{addrTo})
	return resp.Hash, nil
}

func (*XlmService) CheckTxMergerStatus(txId string) error {
	_, err := client.LoadTransaction(txId)
	if err != nil {
		return err
	}
	return nil
}

//获取所有账户信息
func (*XlmService) CheckAddressExists(address string) error {

	account, err := client.LoadAccount(address)
	if err != nil {
		return err
	}
	js, err := json.Marshal(account)
	fmt.Println(string(js))
	return nil
}

//查询当前地址最近一笔交易情况
//获取时间较长，暂时未知返回总是： only expected 1 event, got: 0
//这个方法因为响应时间问题，如果要对接最好限流
func (*XlmService) GetPaymentsNow(address string) error {
	cursor := horizon.Cursor("now")
	fmt.Println("Waiting for a payment...")
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		// Stop streaming after 60 seconds.
		time.Sleep(120 * time.Second)
		cancel()
	}()

	err := client.StreamPayments(ctx, address, &cursor, func(payment horizon.Payment) {
		fmt.Println("Payment type", payment.Type)
		fmt.Println("Payment Paging Token", payment.PagingToken)
		fmt.Println("Payment From", payment.From)
		fmt.Println("Payment To", payment.To)
		fmt.Println("Payment Asset Type", payment.AssetType)
		fmt.Println("Payment Asset Code", payment.AssetCode)
		fmt.Println("Payment Asset Issuer", payment.AssetIssuer)
		fmt.Println("Payment Amount", payment.Amount)
		fmt.Println("Payment Memo Type", payment.Memo.Type)
		fmt.Println("Payment Memo", payment.Memo.Value)
	})

	if err != nil {
		return err
	}
	return nil
}

func (*XlmService) GetAllApi(address string) error {
	hd, err := client.HomeDomainForAccount(address)
	fmt.Println(hd, err)

	/*	client.LoadAccountMergeAmount(&horizon.Payment{})
		client.lo*/

	return nil
}

func (*XlmService) ClearAccount(from, to string) (err error) {
	//数据库获取prv pub key等信息，便于调试--------START------
	actf, err := dhSrv.GetAccountByAddress(from)
	if err != nil {
		return
	}
	//----------------------------------------END-----------
	tx, err := build.Transaction(
		build.SourceAccount{from},
		build.Sequence{1},
		build.TestNetwork,
		build.AccountMerge(
			build.Destination{to},
		),
	)
	if err != nil {
		return
	}

	txe, err := tx.Sign(actf.Seed)
	if err != nil {
		return
	}

	txeB64, err := txe.Base64()

	if err != nil {
		return
	}

	// And finally, send it off to Stellar!
	_, err = client.SubmitTransaction(txeB64) //提交交易
	if err != nil {
		return
	}
	return
}

///////////////////////////////////////////////////内部方法////////////////////////////////////////////////////////
//验证balance是否足够创建tx并成功
//sourceAddress 付款地址
//comparedAmount 目标金额
func checkBalanceEnough(sourceAddress string, comparedAmount float64) error {
	balance, err := xlmSrv.GetBalanceInAddress(sourceAddress)
	if err != nil {
		return err
	}
	if balance < comparedAmount {
		return errors.ERR_NOT_ENOUGH_COIN
	}
	return nil
}

//获取账户序列数
func sequenceForAccount(account string) error {
	num, err := client.SequenceForAccount(account)
	if err != nil {
		return err
	}
	fmt.Println(num)
	return nil
}

/////////////////////////////////////////////////////just for test///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (*XlmService) Other() {
	homeDomain, err := client.HomeDomainForAccount("GBZKTZBJIMLFPUGZUNCUTJCUUREEG4W4UF74K5DRJRZISQNYQP3QOUYX")
	fmt.Println(homeDomain, err)
	offerset, err := client.LoadAccountOffers("GBZKTZBJIMLFPUGZUNCUTJCUUREEG4W4UF74K5DRJRZISQNYQP3QOUYX")
	fmt.Println(offerset, err)
	client.LoadOperation("")

}
