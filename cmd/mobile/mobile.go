package main

import (
	"fmt"

	"github.com/spolabs/wallet-api/src/api/mobile"
)

func main() {
	cfg := mobile.NewConfig()
	cfg.ServerAddr = "127.0.0.1:8620"
	mobile.Init(cfg)
	var wlt string
	var err error
	wlt = "skycoin_aaabbbcccddd"
	if !mobile.IsExist(wlt) {
		wlt, err = mobile.NewWallet("skycoin", "aaabbbcccddd")
		if err != nil {
			fmt.Printf("---err--%v\n", err)
			if err.Error() != "skycoin_aaabbbcccddd already exist" {
				return
			}
		}
		fmt.Printf("wlt---%s\n", wlt)
		address, err := mobile.NewAddress(wlt, 3)
		if err != nil {
			fmt.Printf("---err--%v\n", err)
			return
		}
		fmt.Printf("address---%s\n", address)
	} else {
		fmt.Printf("wlt %s exists\n", wlt)
	}
	addresses, err := mobile.GetAddresses(wlt)
	if err != nil {
		fmt.Printf("---err--%v\n", err)
		return
	}
	fmt.Printf("addresses---%s\n", addresses)
	balance, err := mobile.GetWalletBalance("skycoin", wlt)
	if err != nil {
		fmt.Printf("---balance err--%v\n", err)
		return
	}
	fmt.Printf("balance---%s\n", balance)
	seed := mobile.NewSeed()
	fmt.Printf("seed %s\n", seed)

	txid := "76752105025ba4a84ff0e1ebe2f4a6b1b0f4e27f39433582a5abc419a7fb60de"
	txinfo, err := mobile.GetTransactionByID("skycoin", txid)
	if err != nil {
		fmt.Printf("---tx err--%v\n", err)
		return
	}
	fmt.Printf("tx ---%s\n", txinfo)
	txConfirm, err := mobile.CoinfirmTransaction("skycoin", txid)
	if err != nil {
		fmt.Printf("---tx err--%v\n", err)
		return
	}
	fmt.Printf("tx ---%v\n", txConfirm)

	//destAddr := "wfdokE6kMCfn4UuJhQEe5FNkZwnPzW3kCQ"
	//result, err := mobile.Send("skycoin", wlt, destAddr, "10", nil)
	//if err != nil {
	//fmt.Printf("---send err--%v\n", err)
	//return
	//}
	//fmt.Printf("result---%s\n", result)
}
