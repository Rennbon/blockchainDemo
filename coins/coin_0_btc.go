package coins

import "strconv"

type BtcCoin struct {
	Orgin CoinUnit
}

func (c *BtcCoin) FloatToCoinAmout(f float64) (*CoinAmount, error) {
	return c.praseCoinAmount(strconv.FormatFloat(f, 'f', 8, 64))
}
func (c *BtcCoin) StringToCoinAmout(num string) (ca *CoinAmount, err error) {
	err = regutil.CanPraseBigFloat(num)
	if err != nil {
		return
	}
	return c.praseCoinAmount(num)
}
func (c *BtcCoin) ConvertAmountPrec(ca *CoinAmount, trgt CoinUnit) (caout *CoinAmount, err error) {
	return convertCoinUnit(ca, trgt, c.GetUnitPrec)
}

func (*BtcCoin) GetUnitPrec(cu CoinUnit) (cup *CoinUnitPrec) {
	cup = &CoinUnitPrec{}
	switch cu {
	case CoinBilli:
		cup.UnitName = "BBTC"
		cup.Prec = 17
		return
	case CoinMega:
		cup.UnitName = "MBTC"
		cup.Prec = 14
		return
	case CoinKilo:
		cup.UnitName = "KBTC"
		cup.Prec = 11
		return
	case CoinOrdinary:
		cup.UnitName = "BTC"
		cup.Prec = 8
		return
	case CoinMilli:
		cup.UnitName = "mBTC"
		cup.Prec = 5
		return
	case CoinMicro:
		cup.UnitName = "μBTC"
		cup.Prec = 2
		return
	case CoinBox:
		cup.UnitName = "Satoshi"
		cup.Prec = 0
		return
	default:
		return
	}
}
func (c *BtcCoin) praseCoinAmount(num string) (ca *CoinAmount, err error) {
	return splitStrToNum(num, CoinOrdinary, c.GetUnitPrec, CoinBox)
}
