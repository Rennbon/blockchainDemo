package certs

import (
	"github.com/stellar/go/keypair"
)

type XlmCertService struct {
}

//	生成种子和address
func (*XlmCertService) GenerateSimpleKey() (*Key, error) {
	pair, err := keypair.Random()

	if err != nil {
		return nil, err
	}
	//只有address和seed有用
	key := &Key{Address: pair.Address(), Seed: pair.Seed()}
	return key, nil
}
func (*XlmCertService) GetNewAddress(seed string) (address string, err error) {
	kp, err := keypair.Parse(seed)
	if err != nil {
		return
	}
	return kp.Address(), nil
}
