package coins

import "strconv"

type XmlCoin struct {
}

func (c *XmlCoin) FloatToCoinAmout(f float64) (*CoinAmount, error) {
	return c.praseCoinAmount(strconv.FormatFloat(f, 'f', 8, 64))
}
func (c *XmlCoin) StringToCoinAmout(num string) (ca *CoinAmount, err error) {
	err = regutil.CanPraseBigFloat(num)
	if err != nil {
		return
	}
	return c.praseCoinAmount(num)
}
func (c *XmlCoin) ConvertAmountPrec(ca *CoinAmount, trgt CoinUnit) (caout *CoinAmount, err error) {
	return convertCoinUnit(ca, trgt, c.GetUnitPrec)
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
		cup.UnitName = "BXLM"
		cup.Prec = 14
		return
	case CoinMega:
		cup.UnitName = "MXLM"
		cup.Prec = 11
		return
	case CoinKilo:
		cup.UnitName = "KXLM"
		cup.Prec = 9
		return
	case CoinOrdinary:
		cup.UnitName = "XLM"
		cup.Prec = 6
		return
	case CoinMilli:
		cup.UnitName = "KLumens"
		cup.Prec = 3
		return
	case CoinMicro:
		cup.UnitName = "Lumens"
		cup.Prec = 0
		return
	default:
		return
	}
}
func (c *XmlCoin) praseCoinAmount(num string) (ca *CoinAmount, err error) {
	return splitStrToNum(num, CoinOrdinary, c.GetUnitPrec, CoinMicro)
}
