package main

import (
	"fmt"

	"github.com/spolabs/wallet-api/src/api/mobile"
)

func main() {
	var wlt string
	var err error
	wlt_key := "98ccd8aa2c18bcc283acae1b"
	wlt = "spo_" + wlt_key
	walletDir := "/tmp/wallets" //use default dir ~/.wallet-family

	coinTypes := mobile.GetSupportedCoin()
	fmt.Printf("supported coin types %s\n", coinTypes)
	password := "12"
	err = mobile.Init(walletDir, password)
	if err != nil {
		fmt.Printf("init %s failed %v", "spo", err)
		return
	}
	fmt.Printf("password %s\n", password)
	err = mobile.RegisterNewCoin("spo", "182.92.180.92:8620")
	if err != nil {
		fmt.Printf("register new coin %s failed %v", "spo", err)
		return
	}
	if !mobile.IsExist(wlt) {
		fmt.Printf("wallet not exists\n")
		//newseed := "abcd 1234 8909 bcde xmme adbn nw we hell world then at"
		newseed := "greate test"
		wlt, err = mobile.NewWallet("spo", "lableandseed", newseed, password)
		if err != nil {
			fmt.Printf("---new wallet err--%v\n", err)
			if err.Error() != "spo_"+wlt+" already exist" {
				return
			}
		}
		fmt.Printf("new wlt---%s\n", wlt)
		address, err := mobile.NewAddress(wlt, 2, password)
		if err != nil {
			fmt.Printf("---new address err--%v\n", err)
			return
		}
		fmt.Printf("new addresses---%s\n", address)
	} else {
		fmt.Printf("wlt %s exists\n", wlt)
	}
	addresses, err := mobile.GetAddresses(wlt)
	if err != nil {
		fmt.Printf("---get addresses err--%v\n", err)
		return
	}
	fmt.Printf("get addresses---%s\n", addresses)
	addr := "QsyueWQWvKhqsPDjt1BrdspXVaKayTXUtr"
	pair, err := mobile.GetKeyPairOfAddr(wlt, addr, password)
	if err != nil {
		fmt.Printf("---get keypair err--%v\n", err)
		return
	}
	fmt.Printf("get key pair---%v\n", pair)
	balance, err := mobile.GetWalletBalance("spo", wlt)
	if err != nil {
		fmt.Printf("---balance err--%v\n", err)
		return
	}
	fmt.Printf("balance---%s\n", balance)
	seed := mobile.NewSeed()
	fmt.Printf("new seed %s\n", seed)

	seed1, err := mobile.GetSeed(wlt, password)
	if err != nil {
		fmt.Printf("---seed err--%v\n", err)
		return
	}
	fmt.Printf("get seed---%s\n", seed1)

	txid := "76752105025ba4a84ff0e1ebe2f4a6b1b0f4e27f39433582a5abc419a7fb60de"
	txinfo, err := mobile.GetTransactionByID("spo", txid)
	if err != nil {
		fmt.Printf("---tx err--%v\n", err)
		return
	}
	fmt.Printf("tx ---%s\n", txinfo)
	////err = mobile.Remove(wlt)
	////if err != nil {
	////fmt.Printf("---remove wlt err--%v\n", err)
	////return
	////}
	txConfirm, err := mobile.IsTransactionConfirmed("spo", txid)
	if err != nil {
		fmt.Printf("---tx err--%v\n", err)
		return
	}
	fmt.Printf("tx ---%v\n", txConfirm)

	//sposrc := "QsyueWQWvKhqsPDjt1BrdspXVaKayTXUtr"
	destAddr := "hva72jTmjEdogG4RxNb9uAmDgM1MCfSnLk"
	result, err := mobile.Send("spo", wlt, destAddr, "0.5", password)
	if err != nil {
		fmt.Printf("---send err--%v\n", err)
		return
	}
	fmt.Printf("result---%s\n", result)
}