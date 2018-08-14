package wallets

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"io/ioutil"
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
	/*	//用私钥签署交易签名
		signature, signatureErr := crypto.Sign(tx.Hash().Bytes(), keyWrapper.PrivateKey)
		if signatureErr != nil {
			fmt.Println("signature create error:")
			panic(signatureErr)
		}

		signedTx, signErr := tx.WithSignature(signer, signature)
		if signErr != nil {
			fmt.Println("signer with signature error:")
			panic(signErr)
		}

		//发送交易到网络
		txErr := tokenClient.SendTransaction(context.TODO(), signedTx)

		if txErr != nil {
			fmt.Println("send tx error:")
			panic(txErr)
		}*/
}
