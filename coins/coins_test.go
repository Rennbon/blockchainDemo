package coins_test

import (
	"reflect"
	"testing"

	"github.com/Rennbon/blockchainDemo/coins"
	//"math/big"
)

type CoinsHandler struct {
	coins.DistributionCoiner

	TypeName string
}

////////////////////测试用实体/////////////////////////////

//////////////////////////////////////////////////
var (
	amountInt    = float64(111111.111111)
	amountString = "111111.111111"
)

func (ch *CoinsHandler) LoadService(g coins.DistributionCoiner) error {
	if g != nil {
		ch.DistributionCoiner = g
	}
	typ := reflect.TypeOf(g)
	ch.TypeName = typ.String()

	return nil
}

var (
	btc        *coins.BtcCoin
	btcSerName = "*coins.BtcCoin"
	//xlm        *coins.XmlCoin
	xlmSerName = "*coins.XmlCoin"
	handler    CoinsHandler
)

//prec会约束DecPart的float精度
func TestCoinAmount_String(t *testing.T) {
	handler.LoadService(btc)
	//t.Log(handler.coinAmount.String(handler.GetOrginCoinUnit))
}

func TestCoinAmount_Add(t *testing.T) {
	handler.LoadService(btc)
	/*ca, err := handler.StringToCoinAmout(amountString)
	if err != nil {
		t.Error(err)
		t.Log("请转到Test_StringToCoinAmout调试")
		t.Fail()
	}*/
	/*	t.Log(ca.Amount.String())
		t.Log(handler.Amount.String())
		ca.Add(handler.CoinAmount)
		t.Log(ca.String(handler.GetOrginCoinUnit))*/
}

//测试用例模板
func Test_StringToCoinAmout(t *testing.T) {
	handler.LoadService(btc)
	ca, err := handler.StringToCoinAmout(amountString)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	str1 := ca.String()
	t.Log(str1)
}
func Test_FloatToCoinAmout(t *testing.T) {
	handler.LoadService(btc)
	ca, err := handler.FloatToCoinAmout(float64(amountInt))
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	str1 := ca.String()
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
