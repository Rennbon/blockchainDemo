package coins

import (
	"strconv"
)

type BtcCoin struct {
	*coinAmount
}

func (c *BtcCoin) FloatToCoinAmout(f float64) (CoinAmounter, error) {
	return c.praseCoinAmount(strconv.FormatFloat(f, 'f', 8, 64))
}
func (c *BtcCoin) StringToCoinAmout(num string) (ca CoinAmounter, err error) {
	err = regutil.CanPraseBigFloat(num)
	if err != nil {
		return
	}
	return c.praseCoinAmount(num)
}
func (c *BtcCoin) GetOrginCoinUnit() CoinUnit {
	return CoinBox
}

func (*BtcCoin) GetUnitPrec(cu CoinUnit) (cup *CoinUnitPrec) {
	cup = &CoinUnitPrec{
		coinUnit: cu,
	}
	switch cu {
	case CoinBilli:
		cup.unitName = "BBTC"
		cup.prec = 17

		return
	case CoinMega:
		cup.unitName = "MBTC"
		cup.prec = 14
		return
	case CoinKilo:
		cup.unitName = "KBTC"
		cup.prec = 11
		return
	case CoinOrdinary:
		cup.unitName = "BTC"
		cup.prec = 8
		return
	case CoinMilli:
		cup.unitName = "mBTC"
		cup.prec = 5
		return
	case CoinMicro:
		cup.unitName = "μBTC"
		cup.prec = 2
		return
	case CoinBox:
		cup.unitName = "Satoshi"
		cup.prec = 0
		return
	default:
		return
	}
}
func (c *BtcCoin) praseCoinAmount(num string) (ca CoinAmounter, err error) {
	return stringToAmount(num, CoinOrdinary, c.GetUnitPrec, c.GetOrginCoinUnit())
}
