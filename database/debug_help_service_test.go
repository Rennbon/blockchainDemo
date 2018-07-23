package database_test

import (
	"fmt"
	db "github.com/Rennbon/blockchainDemo/database"
	"testing"
)

var dhSrv db.DHService

func TestAddAccount(t *testing.T) {
	err := dhSrv.AddAccount("debug", "prv", "pub", "seed", "addr", db.BTC)
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
	err := dhSrv.AddTx("name", "txid", []string{"addr"})
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
