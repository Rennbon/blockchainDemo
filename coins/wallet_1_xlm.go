package coins

import (
	"blockchainDemo/cert"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/stellar/go/clients/horizon"
)

type XlmService struct {
}

var (
	certXlmSrv cert.XlmCertService
	XlmSrv     XlmService
)

func (*XlmService) GetNewAddress1(account string, mode AcountRunMode) (address, accountOut string, err error) {
	// pair is the pair that was generated from previous example, or create a pair based on
	// existing keys.
	key, err := certXlmSrv.GenerateSimpleKey()
	if err != nil {
		return "", "", err
	}
	resp, err := http.Get("https://friendbot.stellar.org/?addr=" + key.Address)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	err = dhSrv.AddAccount(account, key.PrivKey, key.PubKey, key.Address, key.Seed)
	if err != nil {
		return "", "", err
	}
	fmt.Println(string(body))
	return key.Address, account, nil
}
func (*XlmService) GetBalanceInAddress1(address string) (balance float64, err error) {
	account, err := horizon.DefaultTestNetClient.LoadAccount(address)
	if err != nil {
		return 0, err
	}
	for _, v := range account.Balances {
		fmt.Println(v.Balance)
	}
	return 0, nil
}
