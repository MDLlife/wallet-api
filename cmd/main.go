package main

import (
	"fmt"

	"github.com/spolabs/wallet-api/src/api/mobile"
)

func main() {

	wltKey := "25aBbV8HvBYZrQS1"
	wltType := "skycoin"
	wlt := wltType + "_" + wltKey
	walletDir := "/tmp/wallets" //use default dir ~/.wallet-family
	newseed := "skycoin is awesome"

	coinTypes := mobile.GetSupportedCoin()
	fmt.Printf("supported coin types %s\n", coinTypes)
	password := "12"
	err := mobile.Init(walletDir, password)
	if err != nil {
		fmt.Printf("init %s failed %v", wltType, err)
		return
	}
	fmt.Printf("password %s\n", password)
	err = mobile.RegisterNewCoin(wltType, "i.spo.network:6420")
	if err != nil {
		fmt.Printf("register new coin %s failed %v", wltType, err)
		return
	}
	if !mobile.IsExist(wlt) {
		fmt.Printf("wallet not exists\n")
		lable := "lable"
		wlt, err = mobile.NewWallet(wltType, lable, newseed, password)
		if err != nil {
			fmt.Printf("---new wallet err--%v\n", err)
			if err.Error() != "_"+wlt+" already exist" {
				return
			}
		}
		fmt.Printf("new wlt---%s\n", wlt)
		address, err := mobile.NewAddress(wlt, 5, password)
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
	fmt.Printf("1 get addresses---%s\n", addresses)
	//addr := "QsyueWQWvKhqsPDjt1BrdspXVaKayTXUtr"
	//pair, err := mobile.GetKeyPairOfAddr(wlt, addr, password)
	//if err != nil {
	//fmt.Printf("---get keypair err--%v\n", err)
	//return
	//}
	//fmt.Printf("get key pair---%v\n", pair)
	balance, err := mobile.GetWalletBalance(wltType, wlt)
	if err != nil {
		fmt.Printf("---balance err--%v\n", err)
		return
	}
	fmt.Printf("balance---%s\n", balance)
	//seed := mobile.NewSeed()
	//fmt.Printf("new seed %s\n", seed)

	//seed1, err := mobile.GetSeed(wlt, password)
	//if err != nil {
	//fmt.Printf("---seed err--%v\n", err)
	//return
	//}
	//fmt.Printf("get seed---%s\n", seed1)

	//txid := "76752105025ba4a84ff0e1ebe2f4a6b1b0f4e27f39433582a5abc419a7fb60de"
	//txinfo, err := mobile.GetTransactionByID(wltType, txid)
	//if err != nil {
	//fmt.Printf("---tx err--%v\n", err)
	//return
	//}
	//fmt.Printf("tx ---%s\n", txinfo)
	////err = mobile.Remove(wlt)
	////if err != nil {
	////fmt.Printf("---remove wlt err--%v\n", err)
	////return
	////}
	//txConfirm, err := mobile.IsTransactionConfirmed(wltType, txid)
	//if err != nil {
	//fmt.Printf("---tx err--%v\n", err)
	//return
	//}
	//fmt.Printf("tx ---%v\n", txConfirm)

	////src := "QsyueWQWvKhqsPDjt1BrdspXVaKayTXUtr"
	destAddr := "2Bzxqu1xBJcTqXD2RgmWhyq8v21zt3A9jmN"
	result, err := mobile.Send(wltType, wlt, destAddr, "1", password)
	if err != nil {
		fmt.Printf("---send err--%v\n", err)
		return
	}
	fmt.Printf("result---%s\n", result)
}
