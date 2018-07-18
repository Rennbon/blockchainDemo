package coins

import (
	"blockchainDemo/cert"
	"blockchainDemo/errors"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"blockchainDemo/database"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"time"
)

type XlmService struct {
}

var (
	certXlmSrv cert.XlmCertService
	xlmSrv     XlmService
)

func (*XlmService) GetNewAddress1(account string, mode AcountRunMode) (address, accountOut string, err error) {
	// pair is the pair that was generated from previous example, or create a pair based on
	// existing keys.
	key, err := certXlmSrv.GenerateSimpleKey()
	if err != nil {
		return "", "", err
	}
	resp, err := http.Get("https://friendbot.stellar.org/?addr=" + key.Address)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	err = dhSrv.AddAccount(account, key.PrivKey, key.PubKey, key.Address, key.Seed, database.XLM)
	if err != nil {
		return "", "", err
	}
	fmt.Println(string(body))
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
	//获取账号总金额
	bls, err := xlmSrv.GetBalanceInAddress1(addrFrom)
	if err != nil {
		return err
	}

	amount := strconv.FormatFloat(transfer, 'f', 6, 64)
	baseFee := float64(100)

	totalTran := transfer + baseFee
	if totalTran > bls {
		return errors.ERR_NOT_ENOUGH_COIN
	}
	//简单转给同一个
	tx, err := build.Transaction(
		build.TestNetwork,
		build.SourceAccount{addrFrom}, //lumens（代币名称）当前主人的地址
		build.AutoSequence{horizon.DefaultTestNetClient},
		build.MemoText{"Just do it"}, //元数据，就是便签
		build.Payment(
			build.Destination{addrTo},  // lumens（代币名称）下个主人的地址
			build.NativeAmount{amount}, //官方payments用string主要防止精度丢失
		),
		/*build.Payment(
			build.Destination{addrTo},  // lumens（代币名称）下个主人的地址
			build.NativeAmount{amount}, //官方payments用string主要防止精度丢失
		),*/
		build.BaseFee{uint64(baseFee)}, //小费，不能100都不给，多笔payment to other则相应的100*N，这个是固定的，和btc什么的有区别
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
		panic(err)
	}

	// And finally, send it off to Stellar!
	resp, err := horizon.DefaultTestNetClient.SubmitTransaction(txeB64) //用签名画押
	if err != nil {
		return err
	}

	fmt.Println("Successful Transaction:")
	fmt.Println("Ledger:", resp.Ledger)
	fmt.Println("Hash:", resp.Hash)
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
