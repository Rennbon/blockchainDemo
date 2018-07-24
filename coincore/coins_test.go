package coincore_test

import (
	"github.com/Rennbon/blockchainDemo/coincore"
	"reflect"
	"testing"
)

type CoinHandler struct {
	coincore.Coiner
	TypeName string
}

func (ch *CoinHandler) LoadService(g coincore.Coiner) error {
	if g != nil {
		ch.Coiner = g
	}
	typ := reflect.TypeOf(g)
	ch.TypeName = typ.String()
	return nil
}

func TesTString(t *testing.T) {

}
