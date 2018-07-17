package coins

import (
	"blockchainDemo/cert"
	"fmt"
	"io/ioutil"
	"net/http"

	"strconv"

	"blockchainDemo/database"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
)

type XlmService struct {
}

var (
	certXlmSrv cert.XlmCertService
	XlmSrv     XlmService
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
	for _, v := range account.Balances {
		fmt.Println(v.Balance)
	}
	return 0, nil
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

	tx, err := build.Transaction(
		build.TestNetwork,
		build.SourceAccount{addrFrom},
		build.AutoSequence{horizon.DefaultTestNetClient},
		build.Payment(
			build.Destination{addrTo},
			build.NativeAmount{amount},
		),
		build.BaseFee{uint64(100)},
	)

	if err != nil {
		return err
	}
	// Sign the transaction to prove you are actually the person sending it.
	txe, err := tx.Sign(actf.Seed)
	if err != nil {
		return err
	}

	txeB64, err := txe.Base64()
	if err != nil {
		panic(err)
	}

	// And finally, send it off to Stellar!
	resp, err := horizon.DefaultTestNetClient.SubmitTransaction(txeB64)
	if err != nil {
		return err
	}

	fmt.Println("Successful Transaction:")
	fmt.Println("Ledger:", resp.Ledger)
	fmt.Println("Hash:", resp.Hash)
	return nil
}
