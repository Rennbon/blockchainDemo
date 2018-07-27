package wallets_test

import (
	"github.com/Rennbon/blockchainDemo/coins"
	"github.com/Rennbon/blockchainDemo/wallets"
	"reflect"
	"strconv"
	"testing"
	"time"
)

type CoinHandler struct {
	wallets.Walleter
	coins.BtcCoin
	TypeName string
}

func (ch *CoinHandler) LoadService(g wallets.Walleter) error {
	if g != nil {
		ch.Walleter = g

	}
	typ := reflect.TypeOf(g)
	ch.TypeName = typ.String()
	return nil
}

var (
	btc        *wallets.BtcService
	btcCoin    coins.BtcCoin
	btcStrName = "*wallets.BtcService"
	xlm        *wallets.XlmService
	xlmCoin    coins.XmlCoin
	xlmStrName = "*wallets.XlmService"
	handler    CoinHandler
)

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
	case btcStrName:
		address, account, err = handler.GetNewAddress("Test"+strconv.FormatInt(time.Now().Unix(), 10), wallets.AddrMode)
		break
	case xlmStrName:
		address, account, err = handler.GetNewAddress("", wallets.AddrMode)
	}
	if err != nil {
		t.Error(err)
	}
	t.Logf("address:%s\n\raccount:%s\n\r", address, account)
}

//测试获取账户余额
func TestGetBalanceInAddress(t *testing.T) {
	handler.LoadService(xlm)
	var (
		balance coins.CoinAmounter
		err     error
	)
	switch handler.TypeName {
	case btcStrName:
		balance, err = handler.GetBalanceInAddress("n3ZT36odeAbur87bTdR6JGCtnWaquGgFZ2")
		break
	case xlmStrName:
		balance, err = handler.GetBalanceInAddress("GD43TZONCLLNDHA5ALVRWZKMATTOKNLLTH3XTAJN6SQK77Q3ZT44QJJV")
		break
	}
	if err != nil {
		t.Error(err)
	}
	t.Log(balance.String())
}

//测试账号到账号
func TestSendAddressToAddress(t *testing.T) {
	handler.LoadService(xlm)
	var (
		txId          string
		err           error
		transfer, fee coins.CoinAmounter
	)

	switch handler.TypeName {
	case btcStrName:
		transfer, _ = btcCoin.StringToCoinAmout("10")
		fee, _ = btcCoin.StringToCoinAmout("0.0001")
		txId, err = handler.SendAddressToAddress("mhAfGecTPa9eZaaNkGJcV7fmUPFi3T2Ki8", "n3ZT36odeAbur87bTdR6JGCtnWaquGgFZ2", transfer, fee)
		break
	case xlmStrName:
		transfer, _ = xlmCoin.StringToCoinAmout("10")
		fee, _ = xlmCoin.StringToCoinAmout("0.0001")
		txId, err = handler.SendAddressToAddress("n4UYCTwXvJ7ijCC9ERGr7qYAuJbiLjUcwT", "mvY3JLZNZrvRewbgMZwvj9CHUJWtQeZjff", transfer, transfer)
		break
	}

	if err != nil {
		t.Error(err)
	}
	t.Log(txId)
}

//测试交易状态
func TestCheckTxMergerStatus(t *testing.T) {
	handler.LoadService(btc)
	var (
		err error
	)
	switch handler.TypeName {
	case btcStrName:
		err = handler.CheckTxMergerStatus("88af33f7e455751ae130746f7a7bd6538fc8b791aa67692dba33ab450ada9c92")
		break
	case xlmStrName:
		err = handler.CheckTxMergerStatus("5b410a62000da9d16fbffdc0b799b219599d6a303cadc6a00db821788f44c53e")
		break
	}
	if err != nil {
		t.Error(err)
	}
}

//测试账号是否存在
func TestCheckAddressExists(t *testing.T) {
	handler.LoadService(btc)
	var (
		err error
	)
	switch handler.TypeName {
	case btcStrName:
		err = handler.CheckAddressExists("046c9bbd1c67db7a99bb45a98c592ec89bffe65174ddd130395d632cb428f7423c3cc4de7d623bc4da321451ddede0e39e8bec0105103268e609cb175ea2fedf91")
		break
	case xlmStrName:
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
	case btcStrName:
		break
	case xlmStrName:
		break
	}
}
