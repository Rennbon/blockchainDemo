package coins_test

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/Rennbon/blockchainDemo/coins"
)

type CoinsHandler struct {
	coins.CoinAmounter
	TypeName string
}

////////////////////测试用实体/////////////////////////////

var simpleca = &coins.CoinAmount{
	big.NewInt(996123812),
	0.123123123,
	coins.CoinOrdinary,
	&coins.CoinUnitPrec{
		8,
		"BTC",
	},
}

//////////////////////////////////////////////////

func (ch *CoinsHandler) LoadService(g coins.CoinAmounter) error {
	if g != nil {
		ch.CoinAmounter = g
	}
	typ := reflect.TypeOf(g)
	ch.TypeName = typ.String()
	return nil
}

var (
	btc        *coins.BtcCoin
	btcSerName = "*coins.BtcCoin"
	handler    CoinsHandler
)

//prec会约束DecPart的float精度
func TestCoinAmount_String(t *testing.T) {
	t.Log(simpleca.String())
}

//测试用例模板
func Test_GetNewAmount(t *testing.T) {
	handler.LoadService(btc)
	switch handler.TypeName {
	case btcSerName:
		ca, err := handler.NewCoinAmout("996123812.123123123")
		if err != nil {
			t.Error(err)
			t.Fail()
		}
		if simpleca.String() != ca.String() {
			t.Error("生成值错误")
			t.Fail()
		}
		t.Log(ca)
		break
	case "*coins.XlmCoin":
		break
	}
}

//测试用例模板
func Test_ConvertAmountPrec(t *testing.T) {
	handler.LoadService(btc)
	switch handler.TypeName {
	case btcSerName:
		caout, err := handler.ConvertAmountPrec(simpleca, coins.CoinMicro)
		if err != nil {
			t.Error(err)
			t.Fail()
		} else {
			t.Log("\r\n原始:", simpleca.String(), "\r\n小数点精度prec:", simpleca.Prec, "\r\n单位:", simpleca.UnitName)
			t.Log("\r\n转变:", caout.String(), "\r\n小数点精度prec:", caout.Prec, "\r\n单位:", caout.UnitName)
		}

		break
	case "*coins.XlmCoin":
		break
	}
}

//测试用例模板
func Test(t *testing.T) {
	handler.LoadService(btc)
	switch handler.TypeName {
	case btcSerName:
		break
	case "*coins.XlmCoin":
		break
	}
}
