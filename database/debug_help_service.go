package database

import "gopkg.in/mgo.v2/bson"

//改服务只为帮助调试blockchain，更方便的获取参数，不适合实际场景
type DHService struct {
}

/* 添加account
name:用户名
prvkey:私钥
pubkey:公钥 */
func (*DHService) AddAccount(name, prvKey, pubKey, address string) error {
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
address：地址 */
func (*DHService) GetAccountByAddress(address string) (accounts *Account, err error) {
	session, col := accountProvider()
	defer session.Close()
	query := bson.M{
		"addr": address,
	}
	var model *Account
	err = col.Find(query).One(&model)
	if err != nil {
		return nil, err
	}
	return model, nil
}

/* 根据address组获取对应的account信息
addresses：地址组 */
func (*DHService) GetAccountByAddresses(addresses []string) (accounts []*Account, err error) {
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

/* 添加tx备份，方便调试
name:用户名
txid:交易id
address:地址 */
func (*DHService) AddTx(name, txId, address string) error {
	session, col := txProvider()
	defer session.Close()
	tx := &Tx{
		Name:    name,
		Address: address,
		TxId:    txId,
	}
	err := col.Insert(tx)
	if err != nil {
		return err
	}
	return nil
}

/* 根据address组获取对应的tx信息
address：地址 */
func (*DHService) GetTxByAddress(address string) (tx *Tx, err error) {
	session, col := txProvider()
	defer session.Close()
	query := bson.M{
		"addr": address,
	}
	var model *Tx
	err = col.Find(query).One(&model)
	if err != nil {
		return nil, err
	}
	return model, nil
}
