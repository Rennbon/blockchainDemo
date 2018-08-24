package certs

import (

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"encoding/hex"
)

type BtcCertService struct {
}

func (*BtcCertService) GenerateSimpleKey() (*Key, error) {

	privKey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return nil, err
	}
	privKey.Serialize()
	//这边的compress会影响到后期Pubkey address的解析方案，解析与此处的方式不一致会导致签名验证不通过
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
func (*BtcCertService) GetNewAddress(pubKey string) (address string, err error) {
	pubKeyByte, err := hex.DecodeString(pubKey)
	if err != nil {
		return
	}


	addrspub, err := btcutil.NewAddressPubKey(pubKeyByte, &chaincfg.RegressionNetParams)
	if err != nil {
		return
	}
	return addrspub.EncodeAddress(), nil

	return
}
