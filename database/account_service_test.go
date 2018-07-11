package database

import (
	"fmt"
	"testing"
)

var act AccountService

func TestAddAccount(t *testing.T) {
	err := act.AddAccount("debug", "prv", "pub", "addr")
	if err != nil {
		t.Error(err)
	}
}

func TestGetAccountByAddresses(t *testing.T) {
	addrs := [1]string{"addr"}
	acts, err := act.GetAccountByAddresses(addrs[:])
	if err != nil {
		t.Error(err)
	}
	for _, v := range acts {
		fmt.Println(*v)
	}
}
