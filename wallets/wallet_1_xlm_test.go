package wallets

import (
	"testing"
)

var xlm XlmService

func TestGetPaymentsNow(t *testing.T) {
	err := xlm.GetPaymentsNow("GD43TZONCLLNDHA5ALVRWZKMATTOKNLLTH3XTAJN6SQK77Q3ZT44QJJV")
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
	err := xlm.sequenceForAccount("GCXQIFHEJDDL7MT3DJVSGPTRSG5K4YPTF2VYFS47DSCDJBOOJSH4TNLL")
	if err != nil {
		t.Error(err)
	}
}
