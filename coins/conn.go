package coins

import (
	"blockchainDemo/config"
	"log"
)

var (
	//验证基础配置是否有效
	keys = []string{"BtcConf", "XlmConf"}
	conf *config.Config
)

//注入先写死
//实际场景，最好依赖注入，维护线程池，动态config
//然后能实现动态更新配置并同步到conn,及一些配置的变量
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
