package coins

import (
	"encoding/json"
	"fmt"
	"testing"
)

var xlm XlmService

/*func TestGetNewAddress1(t *testing.T) {
	address, account, err := xlm.GetNewAddress1("", AddrMode)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("address:%s\n\raccount:%s\n\r", address, account)
}*/
/*func TestGetBalanceInAddress1(t *testing.T) {
	balance, err := xlm.GetBalanceInAddress1("GD43TZONCLLNDHA5ALVRWZKMATTOKNLLTH3XTAJN6SQK77Q3ZT44QJJV")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(balance)
}
*/
/*func TestSendAddressToAddress1(t *testing.T) {

	err := xlm.SendAddressToAddress1(
		"GBZKTZBJIMLFPUGZUNCUTJCUUREEG4W4UF74K5DRJRZISQNYQP3QOUYX",
		"GCXQIFHEJDDL7MT3DJVSGPTRSG5K4YPTF2VYFS47DSCDJBOOJSH4TNLL",
		12,
		0.0001,
	)
	if err != nil {
		t.Error(err)
	}
}*/
func TestGetPaymentsNow(t *testing.T) {
	err := xlm.GetPaymentsNow("GD43TZONCLLNDHA5ALVRWZKMATTOKNLLTH3XTAJN6SQK77Q3ZT44QJJV")
	if err != nil {
		t.Error(err)
	}
}
func TestGetAccount(t *testing.T) {
	err := xlm.GetAccount("GD43TZONCLLNDHA5ALVRWZKMATTOKNLLTH3XTAJN6SQK77Q3ZT44QJJV")
	if err != nil {
		t.Error(err)
	}
}
func TestClearAccount(t *testing.T) {
	err := xlm.ClearAccount("GBZKTZBJIMLFPUGZUNCUTJCUUREEG4W4UF74K5DRJRZISQNYQP3QOUYX", "GD43TZONCLLNDHA5ALVRWZKMATTOKNLLTH3XTAJN6SQK77Q3ZT44QJJV")
	if err != nil {
		t.Error(err)
	}
}

func TestSequenceForAccount(t *testing.T) {
	err := xlm.SequenceForAccount("GCXQIFHEJDDL7MT3DJVSGPTRSG5K4YPTF2VYFS47DSCDJBOOJSH4TNLL")
	if err != nil {
		t.Error(err)
	}
}
func TestGetTxByAddress1(t *testing.T) {
	tx, err := xlm.GetTxByAddress1("5b410a62000da9d16fbffdc0b799b219599d6a303cadc6a00db821788f44c53e")
	if err != nil {
		t.Error(err)
	}
	js, _ := json.Marshal(tx)
	fmt.Println(string(js))
}
