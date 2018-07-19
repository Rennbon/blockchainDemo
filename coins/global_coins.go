package coins

import (
	"blockchainDemo/database"
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
	GetNewAddress(string, AcountRunMode) (address, account string, err error)
	GetBalanceInAddress(string) (balance float64, err error)
	SendAddressToAddress(addrFrom, addrTo string, transfer, fee float64) error
}
