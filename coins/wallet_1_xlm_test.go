package coins

import (
	"fmt"
	"testing"
)

var xlm XlmService

func TestGetNewAddress1(t *testing.T) {
	address, account, err := xlm.GetNewAddress1("", AddrMode)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("address:%s\n\raccount:%s\n\r", address, account)
}
func TestGetBalanceInAddress1(t *testing.T) {
	balance, err := xlm.GetBalanceInAddress1("GBZKTZBJIMLFPUGZUNCUTJCUUREEG4W4UF74K5DRJRZISQNYQP3QOUYX")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(balance)
}
