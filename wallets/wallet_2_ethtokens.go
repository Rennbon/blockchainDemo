package wallets

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

var (
	tokenClient *ethclient.Client
)

type EthTokensService struct {
}

func (c *EthTokensService) GetBalance(address string) {

	rpcDial, err := rpc.Dial("http://127.0.0.1:8545")
	if err != nil {
		panic(err)
	}
	tokenClient := ethclient.NewClient(rpcDial)

	if err != nil {
		panic(err.Error())
	}
	balance, err := tokenClient.BalanceAt(context.TODO(), common.HexToAddress(address), nil)

	fmt.Println(balance)
}
