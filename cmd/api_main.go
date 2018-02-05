package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spolabs/wallet-api/src/api/mobile"
)

func main() {
	fmt.Println("vim-go")
	walletDir := "./temp"
	rpcAddress := "127.0.0.1:6420"
	serv, err := mobile.NewMyService(walletDir)
	if err != nil {
		fmt.Printf("create myservice failed %v\n", err)
		return
	}
	rpc, err := mobile.NewWalletRPC(serv, rpcAddress)
	if err != nil {
		fmt.Printf("new wallet rpc failed %v\n", err)
		return
	}

	wlt, err := serv.CreateWallet("skycoin", "testwallet.wlt", "aaaaa bbbb ccdddd", "first")
	if err != nil {
		fmt.Printf("create wallet failed %v\n", err)
		return
	}

	walletFilename := wlt.GetFilename()

	addresses := wlt.GetAddresses()
	for _, addr := range addresses {
		fmt.Printf("addr: %s\n", addr.String())
	}

	newAddresses, err := serv.NewAddresses(walletFilename, 3)
	if err != nil {
		fmt.Printf("get addresses failed %v\n", err)
		return
	}

	for _, addr := range newAddresses {
		fmt.Printf("new addr: %s\n", addr.String())
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Printf("get current path failed%v\n", err)
		return
	}
	walletPath := filepath.Join(dir, "temp", walletFilename)
	balance, err := rpc.Balance(walletPath)
	if err != nil {
		fmt.Printf("get wallet balance failed %v\n", err)
		return
	}
	fmt.Printf("balance is %+v\n", balance)
}
