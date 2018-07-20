package coins

import (
	"blockchainDemo/database"
	"reflect"
)

type AcountRunMode int

const (
	_AcountRunMode = iota
	NoneMode       //什么都不导入
	PrvMode        //导入私钥
	PubMode        //导入公钥
	AddrMode       //导入地址
)

var dhSrv database.DHService

type Coiner interface {
	GetNewAddress(string, AcountRunMode) (address, accountOut string, err error)
	GetBalanceInAddress(string) (balance float64, err error)
	SendAddressToAddress(addrFrom, addrTo string, transfer, fee float64) (txId string, err error)
	CheckTxMergerStatus(string) error
	CheckAddressExists(string) error
}

type CoinHandler struct {
	Coiner
	TypeName string
}

func (ch *CoinHandler) LoadService(g Coiner) error {
	if g != nil {
		ch.Coiner = g
	}
	typ := reflect.TypeOf(g)
	ch.TypeName = typ.String()
	return nil
}
