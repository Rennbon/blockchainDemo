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
func TestGetAccount(t *testing.T) {
	GetAccount("0x0d13E6594AF3E9E91d1CeBfcA3F344f7F59b4a74", "0x2A41401f94Dc5b97BCB72bF07BF839C74753554b")
}
func TestErc20Transfer(t *testing.T) {
	Erc20Transfer()
}
