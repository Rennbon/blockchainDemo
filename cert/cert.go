package cert

import (
	"encoding/hex"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

type CertService struct {
}
type Key struct {
	PrivKey string
	PubKey  string
	Address string
}

func (*CertService) GenerateSimpleKey() (*Key, error) {

	privKey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return nil, err
	}
	privKey.Serialize()
	privKeyWif, err := btcutil.NewWIF(privKey, &chaincfg.RegressionNetParams, false)
	if err != nil {
		return nil, err
	}
	pubKeySerial := privKey.PubKey().SerializeUncompressed()
	pubKey, err := btcutil.NewAddressPubKey(pubKeySerial, &chaincfg.RegressionNetParams)
	if err != nil {
		return nil, err
	}
	addr := pubKey.EncodeAddress()
	return &Key{PrivKey: privKeyWif.String(), PubKey: pubKey.String(), Address: addr}, nil
}
func (*CertService) GetNewAddress(pubKey string) (string, error) {
	pubKeyByte, err := hex.DecodeString(pubKey)
	if err != nil {
		return "", err
	}
	addrspub, err := btcutil.NewAddressPubKey(pubKeyByte, &chaincfg.RegressionNetParams)
	if err != nil {
		return "", err
	}
	return addrspub.EncodeAddress(), nil
}
