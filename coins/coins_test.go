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
	amountString = "12345.111111"
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
	xlm        *coins.XmlCoin
	xlmSerName = "*coins.XmlCoin"
	handler    CoinsHandler
)

func TestCoinAmount_Add(t *testing.T) {
	handler.LoadService(xlm)
	am1, _ := handler.StringToCoinAmout("100")
	am2, _ := handler.StringToCoinAmout("200")
	am3 := am1
	am1.Add(am2, am2)
	t.Log(am1.String(), am2.String(), am3.String())
}

//测试用例模板
func Test_StringToCoinAmout(t *testing.T) {
	handler.LoadService(xlm)
	var (
		str  string
		trgt string
	)
	switch handler.TypeName {
	case btcSerName:
		str = "12345.67890"
		trgt = "1234567890000"
		break
	case xlmSerName:
		str = "12345.67890"
		trgt = "12345678900"
		break
	}
	ca, err := handler.StringToCoinAmout(str)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if trgt != ca.Val().String() {
		t.Fail()
	}
	t.Logf("\r\n原数据:%s\r\n目标数:%s\r\n实际数:%s", str, trgt, ca.Val().String())
}
func Test_FloatToCoinAmout(t *testing.T) {
	handler.LoadService(xlm)
	var (
		str  string
		trgt string
		flt  float64
	)
	switch handler.TypeName {
	case btcSerName:
		str = "12345.67890"
		flt = float64(12345.67890)
		trgt = "1234567890000"
		break
	case xlmSerName:
		str = "12345.67890"
		flt = float64(12345.67890)
		trgt = "12345678900"
		break
	}
	ca, err := handler.FloatToCoinAmout(flt)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if trgt != ca.Val().String() {
		t.Fail()
	}
	t.Logf("\r\n原数据:%s\r\n目标数:%s\r\n实际数:%s", str, trgt, ca.Val().String())
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
