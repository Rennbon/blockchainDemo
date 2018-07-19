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
	balance, err := xlm.GetBalanceInAddress1("GDJ22GN5AIOL63PCEZM7MJFKX2IYVVCVDO73HTBAKHPRGRPFZBMOQTR4")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(balance)
}

func TestSendAddressToAddress1(t *testing.T) {

	err := xlm.SendAddressToAddress1(
		"GBZKTZBJIMLFPUGZUNCUTJCUUREEG4W4UF74K5DRJRZISQNYQP3QOUYX",
		"GD43TZONCLLNDHA5ALVRWZKMATTOKNLLTH3XTAJN6SQK77Q3ZT44QJJV",
		1,
		0.0001,
	)
	if err != nil {
		t.Error(err)
	}
}
func TestGetPaymentsNow(t *testing.T) {
	err := xlm.GetPaymentsNow("GC2BKLYOOYPDEFJKLKY6FNNRQMGFLVHJKQRGNSSRRGSMPGF32LHCQVGF")
	if err != nil {
		t.Error(err)
	}
}
func TestGetAccount(t *testing.T) {
	err := xlm.GetAccount("GC2BKLYOOYPDEFJKLKY6FNNRQMGFLVHJKQRGNSSRRGSMPGF32LHCQVGF")
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
	err := xlm.SequenceForAccount("GBZKTZBJIMLFPUGZUNCUTJCUUREEG4W4UF74K5DRJRZISQNYQP3QOUYX")
	if err != nil {
		t.Error(err)
	}
}
