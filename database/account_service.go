package database

import "gopkg.in/mgo.v2/bson"

type AccountService struct {
}

/* 添加account
name:用户名
prvkey:私钥
pubkey:公钥 */
func (*AccountService) AddAccount(name, prvKey, pubKey, address string) error {
	session, col := accountProvider()
	defer session.Close()
	acc := &Account{
		Name:    name,
		PrvKey:  prvKey,
		PubKey:  pubKey,
		Address: address,
	}
	err := col.Insert(acc)
	if err != nil {
		return err
	}
	return nil
}

/* 根据address组获取对应的account信息
addresses：地址组 */
func (*AccountService) GetAccountByAddresses(addresses []string) (accounts []*Account, err error) {
	session, col := accountProvider()
	defer session.Close()
	query := bson.M{
		"addr": bson.M{
			"$in": addresses,
		},
	}
	var acts []*Account
	err = col.Find(query).All(&acts)
	if err != nil {
		return nil, err
	}
	return acts, nil
}
