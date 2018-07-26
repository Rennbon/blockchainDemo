package coins_test

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/Rennbon/blockchainDemo/coins"
)

type CoinsHandler struct {
	coins.CoinAmounter
	*coins.CoinAmount
	TypeName string
}

////////////////////测试用实体/////////////////////////////

//////////////////////////////////////////////////

func (ch *CoinsHandler) LoadService(g coins.CoinAmounter) error {
	if g != nil {
		ch.CoinAmounter = g
	}
	typ := reflect.TypeOf(g)
	ch.TypeName = typ.String()
	ch.CoinAmount = &coins.CoinAmount{
		IntPart:      big.NewInt(996123812),
		DecPart:      0.123123123,
		CoinUnit:     coins.CoinOrdinary,
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
	t.Log(handler.CoinAmount.String())
}

//测试用例模板
func Test_GetNewAmount(t *testing.T) {
	handler.LoadService(xlm)
	ca, err := handler.StringToCoinAmout("996123812.123123123")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if handler.CoinAmount.String() != ca.String() {
		t.Error("生成值错误")
		t.Fail()
	}
	t.Log(ca)
}

//测试用例模板
func Test_ConvertAmountPrec(t *testing.T) {
	handler.LoadService(xlm)
	caout, err := handler.ConvertAmountPrec(handler.CoinAmount, coins.CoinMicro)
	if err != nil {
		t.Error(err)
		t.Fail()
	} else {
		t.Log("\r\n原始:", handler.CoinAmount.String(), "\r\n小数点精度prec:", handler.CoinAmount.Prec, "\r\n单位:", handler.CoinAmount.UnitName)
		t.Log("\r\n转变:", caout.String(), "\r\n小数点精度prec:", caout.Prec, "\r\n单位:", caout.UnitName)
	}

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
