package coins

import (
	"github.com/btcsuite/btcd/rpcclient"
)

type BtcService struct {
	Client *rpcclient.Client
}

var btcCli BtcService

func initClinet() {
	cli, err := rpcclient.New(btcConn, nil)
	if err != nil {
		panic("btc rpcclient error.")
	}
	btcCli.Client = cli
	return
}

/*
*获取新的地址
*account:账户名
 */
func (*BtcService) GetNewAddress(account string) (string, error) {
	return "", nil
}

/*
*获取所有account
 */
func (*BtcService) GetAccounts() (accounts []*Account, err error) {
	accs, err := btcCli.Client.ListAccounts()
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
