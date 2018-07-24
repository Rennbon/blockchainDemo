package certs

type Key struct {
	PrivKey string
	PubKey  string
	Address string
	Seed    string
}
type Generater interface {
	//生成公私钥seed等，按需生成
	GenerateSimpleKey() (*Key, error)
	//根据指定公私钥或者seed生成地址address
	GetNewAddress(string) (address string, err error)
}
