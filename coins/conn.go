package coins

import (
	"blockchainDemo/config"
	"log"

	"github.com/btcsuite/btcd/rpcclient"
)

var keys []string = []string{"BtcConf"}
var btcConn *rpcclient.ConnConfig

func init() {
	initConfig()
	initClinet()
}
func initConfig() {
	conf, err := config.LoadConfig()
	if err != nil {
		panic("wallet init LoadConfig panic.")
	}
	err = config.CheckConfig(conf, keys)
	if err != nil {
		panic("wallet init CheckConfig panic.")
	}
	btcConn = &rpcclient.ConnConfig{
		Host:         conf.BtcConf.IP + ":" + conf.BtcConf.Port,
		User:         conf.BtcConf.User,
		Pass:         conf.BtcConf.Passwd,
		HTTPPostMode: true,
		DisableTLS:   true,
	}
	log.Println("coins=>conn=>initConfig sccuess.")
}
