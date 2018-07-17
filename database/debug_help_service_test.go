package database

import (
	"fmt"
	"testing"
)

var dhSrv DHService

func TestAddAccount(t *testing.T) {
	err := dhSrv.AddAccount("debug", "prv", "pub", "seed", "addr", BTC)
	if err != nil {
		t.Error(err)
	}
}
func TestGetAccountByAddress(t *testing.T) {
	act, err := dhSrv.GetAccountByAddress("addr")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(*act)
}

func TestGetAccountByAddresses(t *testing.T) {
	addrs := [1]string{"addr"}
	acts, err := dhSrv.GetAccountByAddresses(addrs[:])
	if err != nil {
		t.Error(err)
	}
	for _, v := range acts {
		fmt.Println(*v)
	}
}

func TestAddTx(t *testing.T) {
	err := dhSrv.AddTx("name", "txid", "addr")
	if err != nil {
		t.Error(err)
	}
}
func TestGetTxByAddress(t *testing.T) {
	tx, err := dhSrv.GetTxByAddress("addr")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(*tx)
}

func Test(t *testing.T) {

}
