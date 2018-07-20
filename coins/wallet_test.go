package coins_test

import (
	"blockchainDemo/coins"
	"fmt"
	"strconv"
	"testing"
	"time"
)

var btc *coins.BtcService
var xlm *coins.XlmService
var handler coins.CoinHandler

/* 相关接口
GetNewAddress(string, AcountRunMode) (address, accountOut string, err error)
GetBalanceInAddress(string) (balance float64, err error)
SendAddressToAddress(addrFrom, addrTo string, transfer, fee float64) (txId string, err error)
CheckTxStatus(string) error

//全局测试代码相关
handler.LoadService(btc)
switch handler.TypeName {
	case "*coins.BtcService":
		break
case "*coins.XlmService":
		break
}
*/

//获取新地址，同事数据库会存储key以便调试
func TestGetNewAddress(t *testing.T) {
	handler.LoadService(btc)
	var (
		address string
		account string
		err     error
	)
	switch handler.TypeName {
	case "*coins.BtcService":
		address, account, err = handler.GetNewAddress("Test"+strconv.FormatInt(time.Now().Unix(), 10), coins.AddrMode)
		break
	case "*coins.XlmService":
		address, account, err = handler.GetNewAddress("", coins.AddrMode)
	}
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("address:%s\n\raccount:%s\n\r", address, account)
}

func TestGetBalanceInAddress(t *testing.T) {
	handler.LoadService(btc)
	var (
		balance float64
		err     error
	)
	switch handler.TypeName {
	case "*coins.BtcService":
		balance, err = handler.GetBalanceInAddress("mhAfGecTPa9eZaaNkGJcV7fmUPFi3T2Ki8")
		break
	case "*coins.XlmService":
		balance, err = handler.GetBalanceInAddress("GD43TZONCLLNDHA5ALVRWZKMATTOKNLLTH3XTAJN6SQK77Q3ZT44QJJV")
		break
	}
	if err != nil {
		t.Error(err)
	}
	bal := strconv.FormatFloat(balance, 'f', 8, 64)
	fmt.Println(bal)
}

func TestSendAddressToAddress(t *testing.T) {
	handler.LoadService(btc)
	var (
		txId string
		err  error
	)
	switch handler.TypeName {
	case "*coins.BtcService":
		txId, err = handler.SendAddressToAddress("n4UYCTwXvJ7ijCC9ERGr7qYAuJbiLjUcwT", "mvY3JLZNZrvRewbgMZwvj9CHUJWtQeZjff", 10, 0.0001)
		break
	case "*coins.XlmService":
		txId, err = handler.SendAddressToAddress("n4UYCTwXvJ7ijCC9ERGr7qYAuJbiLjUcwT", "mvY3JLZNZrvRewbgMZwvj9CHUJWtQeZjff", 10, 0.0001)
		break
	}

	if err != nil {
		t.Error(err)
	}
	fmt.Println(txId)
}
