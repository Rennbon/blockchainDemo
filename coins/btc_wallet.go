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
func (btc *BtcService) GetAccounts() (accounts []*Account, err error) {
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
