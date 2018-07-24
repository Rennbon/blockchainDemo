package coins_test

import (
	"fmt"
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


var btc *coins.BtcCoin
var handler CoinsHandler



func TestCoinAmount_String(t *testing.T) {
	bg := big.NewInt(1000)
	amount := &coins.CoinAmount{bg, 0.00000004, "元", coins.CoinMicro}
	fmt.Println(amount.String(true))
}
//测试用例模板
func TestBtcCoin_GetNewAmount(t *testing.T) {
	handler.LoadService(btc)
	switch handler.TypeName {
	case "*coincore.BtcCoin":
		btresult:=&coins.CoinAmount{
			big.NewInt(996123812),
			0.123123123,
			"BTC",
			coins.CoinOrdinary,

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
