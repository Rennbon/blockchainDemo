package coins

import (
	"github.com/Rennbon/blockchainDemo/utils"
)

type BtcCoin struct {
}

var regSer utils.RegUtil

//
const btcPrec int64 = 1e8

func (*BtcCoin) GetBtcUnitName(cu CoinUnit) CoinUnitName {
	return getBtcUnitName(cu)
}

func (*BtcCoin) GetNewOrdinaryAmount(num string) (ca *CoinAmount, err error) {
	err = regSer.CanPraseBigFloat(num)
	if err != nil {
		return
	}
	return splitStrToNum(num, CoinOrdinary, getBtcUnitName)
}
func (*BtcCoin) ConvertAmountPrec(ca *CoinAmount, trgt CoinUnit) (caout *CoinAmount, err error) {
	return ConvertcoinUnit(ca, trgt, getBtcUnitName)
}
func getBtcUnitName(cu CoinUnit) CoinUnitName {
	switch cu {
	case CoinBilli:
		return "BBTC"
	case CoinMega:
		return "MBTC"
	case CoinKilo:
		return "KBTC"
	case CoinOrdinary:
		return "BTC"
	case CoinMilli:
		return "mBTC"
	case CoinMicro:
		return "μBTC"
	case CoinBox:
		return "Satoshi"
	default:
		return ""
	}
}
