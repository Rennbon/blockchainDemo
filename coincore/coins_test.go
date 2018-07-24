package coincore_test

import (
	"fmt"
	"github.com/Rennbon/blockchainDemo/coincore"
	"math/big"
	"reflect"
	"testing"
)

type CoinsHandler struct {
	coincore.CoinAmounter
	TypeName string
}

func (ch *CoinsHandler) LoadService(g coincore.CoinAmounter) error {
	if g != nil {
		ch.CoinAmounter = g
	}
	typ := reflect.TypeOf(g)
	ch.TypeName = typ.String()
	return nil
}


var btc *coincore.BtcCoin
var handler CoinsHandler



func TestCoinAmount_String(t *testing.T) {
	bg := big.NewInt(1000)
	amount := &coincore.CoinAmount{bg, 0.00000004, "元", coincore.CoinMicro}
	fmt.Println(amount.String(true))
}
//测试用例模板
func TestBtcCoin_GetNewAmount(t *testing.T) {
	handler.LoadService(btc)
	switch handler.TypeName {
	case "*coincore.BtcCoin":
		btresult:=&coincore.CoinAmount{
			big.NewInt(996123812),
			0.123123123,
			"BTC",
			coincore.CoinOrdinary,

		}
		ca,err:= handler.GetNewOrdinaryAmount("996123812.123123123")
		if err!=nil{
			t.Error(err)
			t.Fail()
		}
		if btresult.String(true) !=ca.String(true){
			t.Error("生成值错误")
			t.Fail()
		}
		t.Log(ca)
		break
	case "*coincore.XlmCoin":
		break
	}
}


//测试用例模板
func Test(t *testing.T) {
	handler.LoadService(btc)
	switch handler.TypeName {
	case "*coincore.BtcCoin":
		break
	case "*coincore.XlmCoin":
		break
	}
}
