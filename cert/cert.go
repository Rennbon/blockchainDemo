package main

import (
	"fmt"

	"github.com/btcsuite/btcutil/hdkeychain"
)

func main() {
	seed, err := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(seed)
}
