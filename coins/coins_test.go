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

func TestCoinAmount_String(t *testing.T) {
	bg := big.NewInt(1000)
	amount := &coins.CoinAmount{bg, 0.00000004, coins.CoinMicro, &coins.CoinUnitPrec{}}

	t.Log(amount.String())
}

//测试用例模板
func TestBtcCoin_GetNewAmount(t *testing.T) {
	handler.LoadService(btc)
	switch handler.TypeName {
	case btcSerName:
		btresult := &coins.CoinAmount{
			big.NewInt(996123812),
			0.123123123,
			coins.CoinOrdinary,
			&coins.CoinUnitPrec{
				8,
				"BTC",
			},
		}
		ca, err := handler.NewCoinAmout("996123812.123123123")
		if err != nil {
			t.Error(err)
			t.Fail()
		}
		if btresult.String() != ca.String() {
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
func TestBtcCoin_ConvertAmountPrec(t *testing.T) {
	handler.LoadService(btc)
	switch handler.TypeName {
	case btcSerName:
		ca := &coins.CoinAmount{
			big.NewInt(996123812),
			0.123123123,
			coins.CoinOrdinary,
			&coins.CoinUnitPrec{
				8,
				"BTC",
			},
		}
		caout, err := handler.ConvertAmountPrec(ca, coins.CoinMega)
		if err != nil {
			t.Error(err)
			t.Fail()
		} else {
			t.Log(caout)
		}

		break
	case "*coins.XlmCoin":
		break
	}
}
func BenchmarkBtcCoin_ConvertAmountPrec(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ { //use b.N for looping
		btresult := &coins.CoinAmount{
			big.NewInt(996123812),
			0.123123123,
			coins.CoinOrdinary,
			&coins.CoinUnitPrec{
				8,
				"BTC",
			},
		}
		coins.ConvertcoinUnit1(btresult, coins.CoinBox, btc.GetUnitPrec)
	}
}

func TestBtcCoin_GetNewOrdinaryAmount(t *testing.T) {
	btresult := &coins.CoinAmount{
		big.NewInt(996123812),
		0.123123123,
		coins.CoinOrdinary,
		&coins.CoinUnitPrec{
			8,
			"BTC",
		},
	}
	caout, err := coins.ConvertcoinUnit1(btresult, coins.CoinBox, btc.GetUnitPrec)
	if err != nil {
		t.Error(err)
		t.Fail()
	} else {
		t.Log(caout)
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
