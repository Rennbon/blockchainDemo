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
CheckAddressExists(string) error

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

//测试获取账户余额
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

//测试账号到账号
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

//测试交易状态
func TestCheckTxMergerStatus(t *testing.T) {
	handler.LoadService(btc)
	var (
		err error
	)
	switch handler.TypeName {
	case "*coins.BtcService":
		err = handler.CheckTxMergerStatus("7f11a56ce356281ff5244ae57804da370c3cb0b685367088d10bf67be0a93f59")
		break
	case "*coins.XlmService":
		err = handler.CheckTxMergerStatus("5b410a62000da9d16fbffdc0b799b219599d6a303cadc6a00db821788f44c53e")
		break
	}
	if err != nil {
		t.Error(err)
	}
}

//测试账号是否存在
func TestCheckAddressExists(t *testing.T) {
	handler.LoadService(xlm)
	var (
		err error
	)
	switch handler.TypeName {
	case "*coins.BtcService":
		err = handler.CheckAddressExists("046c9bbd1c67db7a99bb45a98c592ec89bffe65174ddd130395d632cb428f7423c3cc4de7d623bc4da321451ddede0e39e8bec0105103268e609cb175ea2fedf91")
		break
	case "*coins.XlmService":
		err = handler.CheckAddressExists("GD43TZONCLLNDHA5ALVRWZKMATTOKNLLTH3XTAJN6SQK77Q3ZT44QJJV")
		break
	}
	if err != nil {
		t.Error(err)
	}
}

//测试用例模板
func Test(t *testing.T) {
	handler.LoadService(btc)
	switch handler.TypeName {
	case "*coins.BtcService":
		break
	case "*coins.XlmService":
		break
	}
}
