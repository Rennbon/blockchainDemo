package coins

import (
	"github.com/Rennbon/blockchainDemo/utils"
)

type BtcCoin struct {

}

var regSer utils.RegUtil

const btcPrec int64 = 1e8


func (b *BtcCoin) NewCoinAmout(num string) (ca *CoinAmount, err error) {
	err = regSer.CanPraseBigFloat(num)
	if err != nil {
		return
	}
	return splitStrToNum(num, CoinOrdinary, getBtcUnitPrec)
}
func (*BtcCoin) ConvertAmountPrec(ca *CoinAmount, trgt CoinUnit) (caout *CoinAmount, err error) {
	return ConvertcoinUnit(ca, trgt, getBtcUnitPrec)
}


func getBtcUnitPrec(cu CoinUnit) (cup *CoinUnitPrec) {
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
