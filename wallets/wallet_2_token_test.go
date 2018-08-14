package wallets

import "testing"

var etoken EthTokensService

func TestEthTokensService_GetBalance(t *testing.T) {
	etoken.GetBalance("0x3bb953729848873c2f6da94d8273e8c33654f7d8")
}
func TestEthTokensService_GetKey(t *testing.T) {
	etoken.GetAccount()
}

func TestEthTokensService_Transfer(t *testing.T) {
	etoken.Transfer()
}
