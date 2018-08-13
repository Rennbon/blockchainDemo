package certs

import (
	"github.com/ethereum/go-ethereum/crypto"
)

type EthCertService struct {
}

func (*EthCertService) GenerateSimpleKey() (*Key, error) {

	key, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	address := crypto.PubkeyToAddress(key.PublicKey)
	return &Key{PrivKey: key.D.String(), PubKey: key.PublicKey.X.String(), Address: address.String()}, nil
}
func (*EthCertService) GetNewAddress(pubKey string) (address string, err error) {
	return
}
