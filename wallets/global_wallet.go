package wallets

import (
	"github.com/Rennbon/blockchainDemo/coins"
	"github.com/Rennbon/blockchainDemo/database"
)

type AcountRunMode int8

const (
	NoneMode AcountRunMode = iota //什么都不导入
	PrvMode                       //导入私钥
	PubMode                       //导入公钥
	AddrMode                      //导入地址
)

var dhSrv database.DHService

//基础接口约定，各类coin接入必须实现的接口
type Walleter interface {
	/*	//将代币转为最精力度的整数
		// 已法币为例
		// 假设人民币以分为最小单位
		// 那么 10元8毛5分：10.85 出来后是 1085 分
		ConvertAmountToCoinAmount(amount string) *CoinAmount*/
	//生成新地址并导入到共链
	GetNewAddress(string, AcountRunMode) (address, accountOut string, err error)
	//获取指定地址的余额
	GetBalanceInAddress(string) (balance *coins.CoinAmount, err error)
	//账户转钱到账户
	//
	SendAddressToAddress(addrFrom, addrTo string, transfer, fee float64) (txId string, err error)
	//检测交易状态（交易是否被确认）
	//txId:交易id
	CheckTxMergerStatus(txId string) error
	//检测地址是否有效（在公链中存在）
	CheckAddressExists(string) error
}
