package cert

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

type CertService struct {
}
type Key struct {
	PrivKey string
	PubKey  string
}

func (*CertService) GenerateSimplePrivateKey() (*Key, error) {

	privKey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return nil, err
	}

	privKeyWif, err := btcutil.NewWIF(privKey, &chaincfg.MainNetParams, false)
	if err != nil {
		return nil, err
	}
	pubKeySerial := privKey.PubKey().SerializeUncompressed()
	pubKey, err := btcutil.NewAddressPubKey(pubKeySerial, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}
	return &Key{PrivKey: privKeyWif.String(), PubKey: pubKey.String()}, nil
}
