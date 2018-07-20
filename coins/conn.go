package coins

import (
	"blockchainDemo/config"
	"log"
)

var (
	//验证基础配置是否有效
	keys = []string{"BtcConf"}
	conf *config.Config
)

func init() {
	initConfig()
	initBtcClinet(&conf.BtcConf)
	initXlmClinet(&conf.XlmConf)
}
func initConfig() {
	conftemp, err := config.LoadConfig()
	if err != nil {
		panic("wallet init LoadConfig panic.")
	}
	err = config.CheckConfig(conftemp, keys)
	if err != nil {
		panic("wallet init CheckConfig panic.")
	}
	conf = conftemp
	log.Println("coins=>conn=>initConfig sccuess.")
}
