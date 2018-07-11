package coins

import (
	"btcDemo/cert"
	"btcDemo/database"
	"btcDemo/errors"
	"log"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

type BtcService struct {
	client *rpcclient.Client
}

var (
	certSrv cert.CertService
	btcSrv  BtcService
	actSrv  database.AccountService
)

func initClinet() {
	cli, err := rpcclient.New(btcConn, nil)
	if err != nil {
		panic("btc rpcclient error.")
	}
	btcSrv.client = cli
	log.Println("coins=>btc_wallet=>initClinet sccuess.")
}

/*
*获取新的地址
*account:账户名
 */
func (*BtcService) GetNewAddress(account string) (address, accountOut string, err error) {
	key, err := certSrv.GenerateSimpleKey()
	if err != nil {
		return "", "", err
	}
	if err = actSrv.AddAccount(account, key.PrivKey, key.PubKey, key.Address); err != nil {
		return "", "", err
	}
	/* if account, err = btcSrv.AddAddressToWallet(key.PubKey, account); err != nil {
		return "", "", err
	} */
	/* if account, err = btcSrv.AddPubkeyToWallet(key.PubKey, account); err != nil {
		return "", "", err
	} */
	if account, err = btcSrv.AddPrvkeyToWallet(key.PrivKey, account); err != nil {
		return "", "", err
	}
	return key.Address, account, nil
}

func (*BtcService) AddPrvkeyToWallet(prvkey, accoutIn string) (accountOut string, err error) {
	wif, err := btcutil.DecodeWIF(prvkey)
	if err != nil {
		return "", err
	}
	if err = btcSrv.client.ImportPrivKeyLabel(wif, accoutIn); err != nil {
		return "", err
	}
	return accoutIn, nil
}

/* 将publickey对应的address添加到链中，
pubKey 公钥
account 地址自定义名称 */
func (*BtcService) AddPubkeyToWallet(pubKey, accountIn string) (accountOut string, err error) {
	//验证地址是否已存在
	address, err := btcSrv.CheckAddressExisted(pubKey)
	if err != nil {
		return "", err
	}
	if err = btcSrv.client.ImportPubKey(pubKey); err != nil {
		return "", err
	}
	//修改名字 忽略错误
	if err = btcSrv.client.SetAccount(address, accountIn); err != nil {
		return "", nil
	}
	return accountIn, nil
}

/* 将publickey对应的address添加到链中，
pubKey 公钥
account 地址自定义名称 */
func (*BtcService) AddAddressToWallet(pubKey, accountIn string) (accountOut string, err error) {
	//验证地址是否已存在
	address, err := btcSrv.CheckAddressExisted(pubKey)
	if err != nil {
		return "", err
	}
	if err = btcSrv.client.ImportAddress(address.EncodeAddress()); err != nil {
		return "", err
	}
	//修改名字 忽略错误
	if btcSrv.client.SetAccount(address, accountIn) != nil {
		return "", nil
	}
	return accountIn, nil
}

/* 验证publickey对应的地址是否已存在于链中
pubkey 公钥 */
func (*BtcService) CheckAddressExisted(pubKey string) (btcutil.Address, error) {
	address, err := btcutil.DecodeAddress(pubKey, &chaincfg.RegressionNetParams)
	addrValid, err := btcSrv.client.ValidateAddress(address)
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
	accs, err := btcSrv.client.ListAccounts()
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
