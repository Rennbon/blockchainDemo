package database

import (
	"gopkg.in/mgo.v2/bson"
)

type Account struct {
	Id      bson.ObjectId `bson:"_id,omitempty"`
	PrvKey  string        `bson:"prvkey"`
	PubKey  string        `bson:"pubkey"`
	Name    string        `bson:"nam"`
	Address string        `bson:"addr"`
	Seed    string        `bson:"seed"`
}

type Tx struct {
	Id       bson.ObjectId `bson:"_id,omitempty"`
	TxId     string        `bson:"txid"`
	AddressF string        `bson:"addr4"`
	AddressT []string      `bson:"addr2"`
}
