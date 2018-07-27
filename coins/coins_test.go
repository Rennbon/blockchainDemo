package coins_test

import (
	"reflect"
	"testing"

	"github.com/Rennbon/blockchainDemo/coins"
	"math/big"
)

type CoinsHandler struct {
	coins.CoinAmounter
	*coins.CoinAmount
	TypeName string
}

////////////////////测试用实体/////////////////////////////

//////////////////////////////////////////////////
var (
	amountInt    int64 = 111111111111
	amountString       = "111111111111"
)

func (ch *CoinsHandler) LoadService(g coins.CoinAmounter) error {
	if g != nil {
		ch.CoinAmounter = g
	}
	typ := reflect.TypeOf(g)
	ch.TypeName = typ.String()
	ch.CoinAmount = &coins.CoinAmount{
		Amount:       big.NewInt(amountInt),
		CoinUnitPrec: g.GetUnitPrec(coins.CoinOrdinary),
	}
	return nil
}

var (
	btc        *coins.BtcCoin
	btcSerName = "*coins.BtcCoin"
	xlm        *coins.XmlCoin
	xlmSerName = "*coins.XmlCoin"
	handler    CoinsHandler
)

//prec会约束DecPart的float精度
func TestCoinAmount_String(t *testing.T) {
	handler.LoadService(xlm)
	t.Log(handler.CoinAmount.String(handler.GetOrginCoinUnit))
}

func TestCoinAmount_Add(t *testing.T) {
	handler.LoadService(xlm)
	ca, err := handler.StringToCoinAmout(amountString)
	if err != nil {
		t.Error(err)
		t.Log("请转到Test_StringToCoinAmout调试")
		t.Fail()
	}
	t.Log(ca.Amount.String())
	t.Log(handler.Amount.String())
	ca.Add(handler.CoinAmount)
	t.Log(ca.String(handler.GetOrginCoinUnit))
}

//测试用例模板
func Test_StringToCoinAmout(t *testing.T) {
	handler.LoadService(xlm)
	ca, err := handler.StringToCoinAmout(amountString)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	str1 := ca.String(handler.GetOrginCoinUnit)

	t.Log(str1)
}

//测试用例模板
func Test(t *testing.T) {
	handler.LoadService(btc)
	switch handler.TypeName {
	case btcSerName:
		break
	case xlmSerName:
		break
	}
}
