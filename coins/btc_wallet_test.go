package coins

import (
	"fmt"
	"testing"
)

var btc BtcService

func TestGetAccounts(t *testing.T) {
	acts, _ := btc.GetAccounts()
	for _, v := range acts {
		fmt.Println(*v)
	}
}
