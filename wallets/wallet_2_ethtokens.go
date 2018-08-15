package wallets

import (
	"context"
	"fmt"
	"github.com/Rennbon/blockchainDemo/wallets/token"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"io/ioutil"
	"log"
	"math/big"
)

var (
	tokenClient *ethclient.Client
)

func init() {
	client, err := ethclient.Dial("/Users/rennbon/geth/mychain/geth.ipc")
	if err != nil {
		panic(err)
	}
	tokenClient = client
}

type EthTokensService struct {
}

func (c *EthTokensService) GetBalance(address string) {

	balance, err := tokenClient.BalanceAt(context.TODO(), common.HexToAddress(address), nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(balance)
}

func (c *EthTokensService) GetAccount() *keystore.KeyStore {
	ks := keystore.NewKeyStore(
		"/Users/rennbon/geth/mychain/keystore",
		keystore.LightScryptN,
		keystore.LightScryptP)
	return ks
}
func (c *EthTokensService) Transfer() {

	ks := c.GetAccount()
	fromAccDef := accounts.Account{
		Address: common.HexToAddress("0x3bb953729848873c2f6da94d8273e8c33654f7d8"),
	}

	toAccDef := accounts.Account{
		Address: common.HexToAddress("0x0640287c23c3c3c59388f73cf904ec9277887820"),
	}
	KEYJSON_FILEDIR := `/Users/rennbon/geth/mychain/keystore/UTC--2018-08-14T09-58-16.332150769Z--3bb953729848873c2f6da94d8273e8c33654f7d8`
	// 打开账户私钥文件
	keyJson, readErr := ioutil.ReadFile(KEYJSON_FILEDIR)
	if readErr != nil {
		fmt.Println("key json read error:")
		panic(readErr)
	}

	// 解析私钥文件
	keyWrapper, keyErr := keystore.DecryptKey(keyJson, "qwe123456")
	if keyErr != nil {
		fmt.Println("key decrypt error:")
		panic(keyErr)
	}
	// 查找将给定的帐户解析为密钥库中的唯一条目:找到签名的账户
	signAcc, err := ks.Find(fromAccDef)
	if err != nil {
		fmt.Println("account keystore find error:")
		panic(err)
	}
	fmt.Printf("account found: signAcc.addr=%s; signAcc.url=%s\n", signAcc.Address.String(), signAcc.URL)
	fmt.Println()
	// 解鎖簽名的賬户
	errUnlock := ks.Unlock(signAcc, "qwe123456")
	if errUnlock != nil {
		fmt.Println("account unlock error:")
		panic(err)
	}

	nonce, _ := tokenClient.NonceAt(context.TODO(), keyWrapper.Address, nil)
	//(nonce uint64, to common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte)
	// 建立交易
	tx := types.NewTransaction(
		nonce,
		toAccDef.Address,
		big.NewInt(1e18),
		3231744,
		big.NewInt(18000000000),
		common.FromHex("哈哈哈"),
	)
	balance, _ := tokenClient.BalanceAt(context.TODO(), fromAccDef.Address, nil)
	cost := tx.Cost()
	fmt.Println(balance.String(), cost.String())
	if balance.Cmp(cost) < 0 {
		return
	}

	fmt.Printf("key extracted: addr=%s", keyWrapper.Address.String())

	signTx, err := types.SignTx(tx, types.HomesteadSigner{}, keyWrapper.PrivateKey)
	if err != nil {
		panic(err)
	}
	txerr := tokenClient.SendTransaction(context.TODO(), signTx)
	if txerr != nil {
		panic(txerr)
	}
}

func GetAccount(contract string, wallet string) {
	address := common.HexToAddress(wallet)

	token, err := token.NewTokenCaller(common.HexToAddress(contract), tokenClient)
	if err != nil {
		log.Printf("Failed to instantiate a Token contract: %v\n", err)
		panic(err)
	}

	block, err := tokenClient.BlockByNumber(context.TODO(), nil)
	if err != nil {
		log.Printf("Failed to get current block number: %v\n", err)
		panic(err)
	}

	decimals, err := token.Decimals(nil)
	if err != nil {
		log.Printf("Failed to get decimals from contract: %v \n", contract)
		panic(err)
	}

	ethBalance, err := tokenClient.BalanceAt(context.TODO(), address, nil)
	if err != nil {
		log.Printf("Failed to get ethereum balance from address: %v \n", address)
	}

	balance, err := token.BalanceOf(nil, address)
	if err != nil {
		log.Printf("Failed to get balance from contract: %v %v\n", contract, err)
	}

	symbol, err := token.Symbol(nil)
	if err != nil {
		log.Printf("Failed to get symbol from contract: %v \n", contract)

	}

	name, err := token.Name(nil)
	if err != nil {
		log.Printf("Failed to retrieve token name from contract: %v | %v\n", contract, err)

	}

	log.Println(block, decimals, ethBalance, balance, symbol, name)
	return
}

func Erc20Transfer() {

	//首先导入上面生成的账户密钥（json）和密码
	KEYJSON_FILEDIR := `/Users/rennbon/geth/mychain/keystore/UTC--2018-08-15T02-11-45.184398213Z--2a41401f94dc5b97bcb72bf07bf839c74753554b`
	// 打开账户私钥文件
	keyJson, readErr := ioutil.ReadFile(KEYJSON_FILEDIR)
	if readErr != nil {
		fmt.Println("key json read error:")
		panic(readErr)
	}

	// 解析私钥文件
	keyWrapper, keyErr := keystore.DecryptKey(keyJson, "qwe123456")
	if keyErr != nil {
		fmt.Println("key decrypt error:")
		panic(keyErr)
	}
	auth := bind.NewKeyedTransactor(keyWrapper.PrivateKey)

	from := common.HexToAddress("0x2A41401f94Dc5b97BCB72bF07BF839C74753554b")
	to := common.HexToAddress("0xAf67c0bC065D8D22cFDc2d8971062Ff8c10e3154")
	contract := common.HexToAddress("0x0d13E6594AF3E9E91d1CeBfcA3F344f7F59b4a74")
	//来源地址
	tkn, err := token.NewToken(contract, tokenClient)
	if err != nil {
		panic(err)
	}
	balance, err := tkn.BalanceOf(nil, from)
	if err != nil {
		panic(err)
	}
	log.Println(balance)
	//每个代币都会有相应的位数，例如eos是18位，那么我们转账的时候，需要在金额后面加18个0
	decimal, err := tkn.Decimals(nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(decimal)

	// TO ADDRESS
	tx, err := tkn.Transfer(auth, to, big.NewInt(1e8))
	if nil != err {
		fmt.Printf("err: %v \n", err)
		return
	}
	log.Println(tx.Hash().String())

}
