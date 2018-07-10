package database

import (
	"btcDemo/config"
	"log"
	"strings"
	"time"

	mgo "gopkg.in/mgo.v2"
)

const (
	col_account string = "account"
)

var (
	db_blockChain string = "blockChain" //会根据配置文件动态变化值
	blockChain    *mgo.Session
	keys          []string = []string{"Mongo"}
)

func init() {
	initMongo()
}
func initMongo() {
	conf, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	var keys []string
	err = config.CheckConfig(conf, keys)
	if err != nil {
		panic(err)
	}
	loadBlockChainSession(conf)
	log.Println("database=>mongo=>initMongo sucess")
}
func loadBlockChainSession(c *config.Config) error {
	mongoDBInfo := &mgo.DialInfo{
		Addrs:     strings.Split(c.Mongo.Addr, ","),
		Timeout:   c.Mongo.Timeout * time.Second,
		PoolLimit: c.Mongo.PoolLimit,
		Database:  c.Mongo.Database,
	}
	session, err := mgo.DialWithInfo(mongoDBInfo)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	blockChain = session
	db_blockChain = mongoDBInfo.Database
	return nil
}

func accountProvider() (*mgo.Session, *mgo.Collection) {
	session := blockChain.Clone()
	col := session.DB(db_blockChain).C(col_account)
	return session, col
}
