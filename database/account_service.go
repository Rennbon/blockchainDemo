package database

type AccountService struct {
}

/* 添加account
name:用户名
prvkey:私钥
pubkey:公钥 */
func (*AccountService) AddAccount(name, prvKey, pubKey string) error {
	session, col := accountProvider()
	defer session.Close()
	acc := &Account{
		Name:   name,
		PrvKey: prvKey,
		PubKey: pubKey,
	}
	err := col.Insert(acc)
	if err != nil {
		return err
	}
	return nil
}
