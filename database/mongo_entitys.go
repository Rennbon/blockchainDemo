package database

import (
	"gopkg.in/mgo.v2/bson"
)

type Account struct {
	Id     bson.ObjectId `bson:"_id,omitempty"`
	PrvKey string        `bson:"prvkey"`
	PubKey string        `bson:"pubkey"`
	Name   string        `bson:"nam"`
}
