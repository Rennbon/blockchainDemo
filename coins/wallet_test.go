package coins

import (
	"fmt"
	"testing"
)

func TestInit(t *testing.T) {
	initConifg()
	fmt.Println(btcConn)
	fmt.Println(btcCli)
}
