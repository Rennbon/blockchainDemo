package coins

import (
	"blockchainDemo/cert"
	"context"
	"fmt"
	"strconv"

	"time"

	"blockchainDemo/database"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"io/ioutil"
	"net/http"
)

type XlmService struct {
}

var (
	certXlmSrv cert.XlmCertService
	xlmSrv     XlmService
)

func (*XlmService) GetNewAddress1(account string, mode AcountRunMode) (address, accountOut string, err error) {
	key, err := certXlmSrv.GenerateSimpleKey()
	if err != nil {
		return
	}
	//----------测试网络创建--------start------------
	resp, err := http.Get("https://friendbot.stellar.org/?addr=" + key.Address)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	//----------测试网络创建--------end------------
	//----------真实网络创建--------start------------
	//创建新账户必须从已有资金的账户转账生成，所以理论上生产环境要用，必须要提供一个有钱的账户作为God来来创造一切
	/*	seed := "SAACHR2TWFAJKLLLC5TEYTSYPXA7AIBM6A2KZ7MQ4XEYRJEZFNOR6VOC"
		tx, err := build.Transaction(
			build.TestNetwork,
			build.Sequence{1},
			build.SourceAccount{seed},
			build.MemoText{"Create Account"}, //元数据，就是便签
			build.CreateAccount(
				build.Destination{key.Address},
				build.NativeAmount{"50"},
			),
		)
		if err != nil {
			return
		}
		// Sign the transaction to prove you are actually the person sending it.
		txe, err := tx.Sign(seed) //签名需要用seed
		if err != nil {
			return
		}
		txeB64, err := txe.Base64()
		if err != nil {
			return
		}
		rsp, err := horizon.DefaultTestNetClient.SubmitTransaction(txeB64) //用签名画押
		if err != nil {
			return
		}
		fmt.Println(rsp)*/
	//----------真实网络创建--------end------------
	err = dhSrv.AddAccount(account, key.PrivKey, key.PubKey, key.Address, key.Seed, database.XLM)
	if err != nil {
		return
	}
	return key.Address, account, nil
}
func (*XlmService) GetBalanceInAddress1(address string) (balance float64, err error) {
	account, err := horizon.DefaultTestNetClient.LoadAccount(address)
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
func (*XlmService) SendAddressToAddress1(addrFrom, addrTo string, transfer, fee float64) error {
	//数据库获取prv pub key等信息，便于调试--------START------
	actf, err := dhSrv.GetAccountByAddress(addrFrom)
	if err != nil {
		return err
	}
	//----------------------------------------END-----------
	//验证地址是否有效
	if _, err := horizon.DefaultTestNetClient.LoadAccount(addrTo); err != nil {
		return nil
	}

	amount := strconv.FormatFloat(transfer, 'f', 6, 64)
	//获取账号总金额
	/*	bls, err := xlmSrv.GetBalanceInAddress1(addrFrom)
		if err != nil {
			return err
		}


		baseFee := float64(100)

		totalTran := transfer + baseFee
		if totalTran > bls {
			return errors.ERR_NOT_ENOUGH_COIN
		}*/
	//小费是自己扣的，不需要这边实现，金额总数也不需要验证，当然可以验证
	tx, err := build.Transaction(
		build.TestNetwork,
		build.SourceAccount{addrFrom},                    //lumens（代币名称）当前主人的地址
		build.AutoSequence{horizon.DefaultTestNetClient}, //选择网络
		build.MemoText{"Just do it"},                     //元数据，就是便签
		build.Payment(
			build.Destination{addrTo},  // lumens（代币名称）下个主人的地址
			build.NativeAmount{amount}, //官方payments用string主要防止精度丢失
		),
	)

	if err != nil {
		return err
	}
	// Sign the transaction to prove you are actually the person sending it.
	txe, err := tx.Sign(actf.Seed) //签名需要用seed
	if err != nil {
		return err
	}

	txeB64, err := txe.Base64()
	if err != nil {
		return err
	}

	// And finally, send it off to Stellar!
	_, err = horizon.DefaultTestNetClient.SubmitTransaction(txeB64) //用签名画押
	if err != nil {
		return err
	}
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
		time.Sleep(60 * time.Second)
		cancel()
	}()

	err := horizon.DefaultTestNetClient.StreamPayments(ctx, address, &cursor, func(payment horizon.Payment) {
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

//获取所有账户信息
func (*XlmService) GetAccount(address string) error {

	account, err := horizon.DefaultTestNetClient.LoadAccount(address)
	if err != nil {
		return err
	}
	fmt.Println(account)
	return nil
}
func (*XlmService) GetAllApi(address string) error {
	hd, err := horizon.DefaultTestNetClient.HomeDomainForAccount(address)
	fmt.Println(hd, err)

	/*	horizon.DefaultTestNetClient.LoadAccountMergeAmount(&horizon.Payment{})
		horizon.DefaultTestNetClient.lo*/

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
	_, err = horizon.DefaultTestNetClient.SubmitTransaction(txeB64) //用签名画押
	if err != nil {
		return
	}
	return
}
