package coincore_test

import (
	"fmt"
	"github.com/Rennbon/blockchainDemo/coincore"
	"math/big"
	"reflect"
	"testing"
)

type CoinHandler struct {
	coincore.CoinAmounter
	TypeName string
}

func (ch *CoinHandler) LoadService(g coincore.CoinAmounter) error {
	if g != nil {
		ch.CoinAmounter = g
	}
	typ := reflect.TypeOf(g)
	ch.TypeName = typ.String()
	return nil
}

func TestCoinAmount_String(t *testing.T) {
	bg := big.NewInt(1000)
	amount := &coincore.CoinAmount{bg, 0.00000004, "元", coincore.CoinMicro}
	fmt.Println(amount.String(true))
}
