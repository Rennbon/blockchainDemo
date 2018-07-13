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
}

type Tx struct {
	Id      bson.ObjectId `bson:"_id,omitempty"`
	TxId    string        `bson:"txid"`
	Address string        `bson:"addr"`
	Name    string        `bson:"nam"`
}
