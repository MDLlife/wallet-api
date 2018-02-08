package mobile

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spolabs/wallet-api/src/coin"
	"github.com/spolabs/wallet-api/src/wallet"
)

//go:generate gomobile bind -target=ios github.com/spolabs/wallet-api/src/api/mobile

var coinMap map[string]Coiner

// Init initialize wallet dir and coin manager.
func Init(walletDir string) {
	wallet.InitDir(walletDir)
	coinMap = make(map[string]Coiner)
}

// RegisterNewCoin register a new coin to wallet
// the server address is consisted of ip and port, eg: 127.0.0.1:6420
func RegisterNewCoin(coinType, serverAddr string) error {
	if _, ok := coinMap[coinType]; ok {
		return fmt.Errorf("coin %s already registed", coinType)
	}

	coinMap[coinType] = newCoin(coinType, serverAddr)
	return nil
}

// NewWallet create a new wallet base on the wallet type and seed
func NewWallet(coinType, lable, seed string) (string, error) {
	wlt, err := wallet.New(coinType, lable, seed)
	if err != nil {
		return "", err
	}
	return wlt.GetID(), nil
}

// IsExist wallet exists or not
func IsExist(walletID string) bool {
	return wallet.IsExist(walletID)
}

// IsContain wallet contains address (format "a1,a2,a3") or not
func IsContain(walletID string, addrs string) (bool, error) {
	addresses := strings.Split(addrs, ",")
	return wallet.IsContain(walletID, addresses)
}

// NewAddress generate address in specific wallet.
func NewAddress(walletID string, num int) (string, error) {
	es, err := wallet.NewAddresses(walletID, num)
	if err != nil {
		return "", err
	}
	var res = struct {
		Entries []coin.AddressEntry `json:"addresses"`
	}{
		es,
	}
	d, err := json.Marshal(res)
	if err != nil {
		return "", err
	}

	return string(d), nil
}

// GetAddresses return all addresses in the wallet.
// returns {"addresses":["jvzYqvdZs17i67cxZ5R8zGE4446JGPVYyz","FNhfaxwWgDVfuXdn2kUoMkxpDFGvqoSPzq","5spraVxAAkFC9j1cpMEdMu7CoV3iHRG7pG"]}
func GetAddresses(walletID string) (string, error) {
	addrs, err := wallet.GetAddresses(walletID)
	if err != nil {
		return "", err
	}
	var res = struct {
		Addresses []string `json:"addresses"`
	}{
		addrs,
	}

	d, err := json.Marshal(res)
	if err != nil {
		return "", err
	}

	return string(d), nil
}

// Remove delete wallet.
func Remove(walletID string) error {
	return wallet.Remove(walletID)
}

// GetKeyPairOfAddr get pubkey and seckey pair of address in specific wallet.
func GetKeyPairOfAddr(walletID string, addr string) (string, error) {
	p, s, err := wallet.GetKeypair(walletID, addr)
	if err != nil {
		return "", err
	}
	var res = struct {
		Pubkey string `json:"pubkey"`
		Seckey string `json:"seckey"`
	}{
		p,
		s,
	}

	d, err := json.Marshal(res)
	if err != nil {
		return "", err
	}
	return string(d), nil
}

// GetBalance return balance of a specific address.
// returns {"balance":"70.000000"}
func GetBalance(coinType string, address string) (string, error) {
	coin, ok := coinMap[coinType]
	if !ok {
		return "", fmt.Errorf("%s is not supported", coinType)
	}

	if err := coin.ValidateAddr(address); err != nil {
		return "", err
	}

	bal, err := coin.GetBalance(address)
	if err != nil {
		return "", err
	}

	var res = struct {
		Balance string `json:"balance"`
	}{
		bal,
	}

	d, err := json.Marshal(res)
	if err != nil {
		return "", err
	}
	return string(d), nil
}

// GetWalletBalance return balance of wallet.
func GetWalletBalance(coinType string, wltID string) (string, error) {
	coin, ok := coinMap[coinType]
	if !ok {
		return "", fmt.Errorf("%s is not supported", coinType)
	}

	addrs, err := wallet.GetAddresses(wltID)
	if err != nil {
		return "", err
	}

	bal, err := coin.GetBalance(strings.Join(addrs, ","))
	if err != nil {
		return "", err
	}
	var res = struct {
		Balance string `json:"balance"`
	}{
		bal,
	}

	d, err := json.Marshal(res)
	if err != nil {
		return "", err
	}
	return string(d), nil
}

// Send send coins, support bitcoin and all coins in skycoin ledger
func Send(coinType, wid, toAddr, amount string) (string, error) {
	coin, ok := coinMap[coinType]
	if !ok {
		return "", fmt.Errorf("%s is not supported", coinType)
	}

	return coin.Send(wid, toAddr, amount)
}

// GetTransactionByID gets transaction verbose info by id
func GetTransactionByID(coinType, txid string) (string, error) {
	coin, ok := coinMap[coinType]
	if !ok {
		return "", fmt.Errorf("%s is not supported", coinType)
	}

	return coin.GetTransactionByID(txid)
}

// IsTransactionConfirmed gets transaction verbose info by id
func IsTransactionConfirmed(coinType, txid string) (bool, error) {
	coin, ok := coinMap[coinType]
	if !ok {
		return false, fmt.Errorf("%s is not supported", coinType)
	}

	return coin.IsTransactionConfirmed(txid)
}

// ValidateAddress validate the address
func ValidateAddress(coinType, addr string) (bool, error) {
	coin, ok := coinMap[coinType]
	if !ok {
		return false, fmt.Errorf("%s is not supported", coinType)
	}

	if err := coin.ValidateAddr(addr); err != nil {
		return false, err
	}

	return true, nil
}

// NewSeed generates mnemonic seed
func NewSeed() string {
	return wallet.NewSeed()
}

// GetSeed returun wallet seed
func GetSeed(walletID string) (string, error) {
	return wallet.GetSeed(walletID)
}

func getPrivateKey(walletID string) coin.GetPrivKey {
	return func(addr string) (string, error) {
		_, s, err := wallet.GetKeypair(walletID, addr)
		return s, err
	}
}
