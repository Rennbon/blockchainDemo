package coins

import "strconv"

type XmlCoin struct{}

func (c *XmlCoin) FloatToCoinAmout(f float64) (CoinAmounter, error) {
	return c.praseCoinAmount(strconv.FormatFloat(f, 'f', 6, 64))
}
func (c *XmlCoin) StringToCoinAmout(num string) (ca CoinAmounter, err error) {
	err = regutil.CanPraseBigFloat(num)
	if err != nil {
		return
	}
	return c.praseCoinAmount(num)
}
func (c *XmlCoin) GetOrginCoinUnit() CoinUnit {
	return CoinMicro
}

func (c *XmlCoin) praseCoinAmount(num string) (ca CoinAmounter, err error) {
	return stringToAmount(num, CoinOrdinary, c.GetUnitPrec, c.GetOrginCoinUnit())
}

/*
	0.000001
	baseFee       float64 = 0.0001 //小费基数（单位:xlm）
	baseFeeLemuns uint64  = 100    //小费 (单位：lumens)
*/
func (*XmlCoin) GetUnitPrec(cu CoinUnit) (cup *CoinUnitPrec) {
	cup = &CoinUnitPrec{}
	switch cu {
	case CoinBilli:
		cup.unitName = "BXLM"
		cup.prec = 15
		return
	case CoinMega:
		cup.unitName = "MXLM"
		cup.prec = 12
		return
	case CoinKilo:
		cup.unitName = "KXLM"
		cup.prec = 9
		return
	case CoinOrdinary:
		cup.unitName = "XLM"
		cup.prec = 7
		return
	case CoinMilli:
		cup.unitName = "KLumens"
		cup.prec = 3
		return
	case CoinMicro:
		cup.unitName = "Lumens"
		cup.prec = 0
		return
	default:
		return
	}
}
