package coins

import (
	"btcDemo/cert"
	"btcDemo/database"
	"btcDemo/errors"
	"btcd/txscript"
	"btcd/wire"
	"btcutil"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
)

type BtcService struct {
	client *rpcclient.Client
}

var (
	certService cert.CertService
	btcCli      BtcService
	actService  database.AccountService
)

func initClinet() {
	cli, err := rpcclient.New(btcConn, nil)
	if err != nil {
		panic("btc rpcclient error.")
	}
	btcCli.client = cli
	return
}

/*
*获取新的地址
*account:账户名
 */
func (*BtcService) GetNewAddress(account string) (string, error) {
	key, err := certService.GenerateSimpleKey()
	if err != nil {
		return "", err
	}
	err = actService.AddAccount(account, key.PrivKey, key.PubKey)
	if err != nil {
		return "", err
	}
	return key.Address, nil
}

func (*BtcService) AddAddressToChain(pubKey, account string) error {
	address, err := btcCli.CheckAddressExisted(pubKey)
	if err != nil {
		return err
	}
	if err = btcCli.client.ImportAddress(address.EncodeAddress()); err != nil {
		return err
	}
	return nil
}
func (*BtcService) CheckAddressExisted(pubKey string) (btcutil.Address, error) {
	address, err := btcutil.DecodeAddress(pubKey, &chaincfg.RegressionNetParams)
	addrValid, err := btcCli.client.ValidateAddress(address)
	if err != nil {
		return nil, err
	}
	if addrValid.IsWatchOnly {
		return address, errors.ERR_DATA_EXISTS
	}
	return address, nil
}

/*
*获取所有account
 */
func (*BtcService) GetAccounts() (accounts []*Account, err error) {
	accs, err := btcCli.client.ListAccounts()
	if err != nil {
		return nil, err
	}
	for k, v := range accs {
		accounts = append(accounts, &Account{
			Amount: v.ToBTC(),
			Name:   k,
			Unit:   "BTC",
		})
	}
	return accounts, nil
}
func (*BtcService) GetBalanceInAddress(address string) (balance int64, err error) {
	return
}
func (*BtcService) SendBtcToAddress(addrFrom, addrTo string, amount, fee int64) error {

	// 1. 构造输出
	outputs := []*wire.TxOut{}

	// 输出1, 给form
	addrf, err := btcutil.DecodeAddress(addrFrom, &chaincfg.RegressionNetParams)
	if err != nil {
		return err
	}
	pkScriptf, err := txscript.PayToAddrScript(addrf)
	if err != nil {
		return err
	}
	outputs = append(outputs, wire.NewTxOut(amount, pkScriptf))
	//输出2，给To
	addrt, err := btcutil.DecodeAddress(addrTo, &chaincfg.RegressionNetParams)
	if err != nil {
		return err
	}

	pkScriptt, err := txscript.PayToAddrScript(addrt)
	if err != nil {
		return err
	}
	outputs = append(outputs, wire.NewTxOut(amount, pkScriptt))

	return nil
}
