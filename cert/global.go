package cert

import (
	"reflect"
)

type Generater interface {
	GenerateSimpleKey() (*Key, error)
	GetNewAddress(string) (string, error)
}

type CertHandler struct {
	Generater
	TypeName string
}

func (ch *CertHandler) LoadService(g Generater) error {
	if g != nil {
		ch.Generater = g
	}
	typ := reflect.TypeOf(g)
	ch.TypeName = typ.String()
	return nil
}
