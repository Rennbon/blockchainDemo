package wallets

import "testing"

var etoken EthTokensService

func TestEthTokensService_GetBalance(t *testing.T) {
	etoken.GetBalance("0x911ba3baFb43798BF4443a0BA93f3470Ab10E1c5")
}
