package coincore

import (
	"github.com/Rennbon/blockchainDemo/utils"
	"strings"
)

type BtcCoin struct {
	CoinAmount
}

//
const btcPrec int64 = 1e8

func (*BtcCoin) GetCoinUnitName(cu CoinUnit) CoinUnitName {
	return getBtcUnitName(cu)
}

func (*BtcCoin) GetNewAmount(num string, trgt CoinUnit) *CoinAmount {

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
