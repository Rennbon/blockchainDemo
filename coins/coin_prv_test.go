package coins

import (
	"math/big"
	"testing"
)

var simpleca = &CoinAmount{
	big.NewInt(996123812),
	0.12312312,
	CoinOrdinary,
	&CoinUnitPrec{
		8,
		"BTC",
	},
}
var orgnca = &CoinAmount{big.NewInt(12345), 0.6789, CoinOrdinary, &CoinUnitPrec{8, "BTC"}}

/**
目标
	CoinUnit     const    	    name 	prec gap
	CoinBilli    CoinUnit = 9   BBTC 	17   9
	CoinMega     CoinUnit = 6   MBTC 	14   6
	CoinKilo     CoinUnit = 3   KBTC 	11   3
	CoinOrdinary CoinUnit = 0   BTC  	8    0
	CoinMilli    CoinUnit = -3  mBTC 	5   -3
	CoinMicro    CoinUnit = -6  μBTC 	2   -6
	CoinBox      CoinUnit = -8  Satoshi 0 	-8
*/
var slcTarget = []*CoinAmount{
	{big.NewInt(0), 0.0000123456789, CoinBilli, &CoinUnitPrec{17, "BBTC"}},
	{big.NewInt(0), 0.0123456789, CoinMega, &CoinUnitPrec{14, "MBTC"}},
	{big.NewInt(12), 0.3456789, CoinKilo, &CoinUnitPrec{11, "KBTC"}},
	{big.NewInt(12345), 0.6789, CoinOrdinary, &CoinUnitPrec{8, "BTC"}},
	{big.NewInt(12345678), 0.9, CoinMilli, &CoinUnitPrec{5, "mBTC"}},
	{big.NewInt(12345678900), 0, CoinMicro, &CoinUnitPrec{2, "μBTC"}},
	{big.NewInt(1234567890000), 0, CoinBox, &CoinUnitPrec{0, "Satoshi"}},
}

var gup = &BtcCoin{}

func Benchmark_ConvertcoinUnit(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for k, v := range slcTarget {
			if k == 4 {
				continue
			}
			convertCoinUnit(orgnca, v.CoinUnit, gup.GetUnitPrec)
		}
	}
}

func Benchmark_ConvertcoinUnit1(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for k, v := range slcTarget {
			if k == 4 {
				continue
			}
			convertCoinUnit1(orgnca, v.CoinUnit, gup.GetUnitPrec)
		}
	}
}

func TestConvertCoinUnit(t *testing.T) {
	for k, v := range slcTarget {
		if k == 4 {
			continue
		}
		caout, err := convertCoinUnit(orgnca, v.CoinUnit, gup.GetUnitPrec)
		if err != nil {
			t.Error(err)
		} else if caout.String() != v.String() {
			t.Errorf("测试失败,下标：%d\r\n原文参数:%s\r\n预期结果:%s\r\n实际结果:%s\r\n", k, orgnca.String(), v.String(), caout.String())
			t.Fail()
		} else {
			t.Logf("测试成功,下标：%d\r\n原文参数:%s\r\n预期结果:%s\r\n实际结果:%s\r\n", k, orgnca.String(), v.String(), caout.String())
		}
	}
}
func TestConvertCoinUnit1(t *testing.T) {
	for k, v := range slcTarget {
		if k == 4 {
			continue
		}
		caout, err := convertCoinUnit1(orgnca, v.CoinUnit, gup.GetUnitPrec)
		if err != nil {
			t.Error(err)
		} else if caout.String() != v.String() {
			t.Errorf("测试失败,下标：%d\r\n原文参数:%s\r\n预期结果:%s\r\n实际结果:%s\r\n", k, orgnca.String(), v.String(), caout.String())
			t.Fail()
		} else {
			t.Logf("测试成功,下标：%d\r\n原文参数:%s\r\n预期结果:%s\r\n实际结果:%s\r\n", k, orgnca.String(), v.String(), caout.String())
		}
	}

}
