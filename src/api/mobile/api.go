package mobile

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/skycoin/skycoin/src/util/droplet"
	skywallet "github.com/skycoin/skycoin/src/wallet"
	"github.com/MDLlife/wallet-api/src/coin"
	"github.com/MDLlife/wallet-api/src/coin/aynrandcoin"
	"github.com/MDLlife/wallet-api/src/coin/mzcoin"
	"github.com/MDLlife/wallet-api/src/coin/samos"
	"github.com/MDLlife/wallet-api/src/coin/shellcoin"
	"github.com/MDLlife/wallet-api/src/coin/skycoin"
	"github.com/MDLlife/wallet-api/src/coin/spo"
	"github.com/MDLlife/wallet-api/src/coin/mdl"
	"github.com/MDLlife/wallet-api/src/coin/suncoin"
	"github.com/MDLlife/wallet-api/src/wallet"
)

//go:generate gomobile bind -target=ios github.com/spolabs/wallet-api/src/api/mobile

var coinMap map[string]Coiner

// Init initialize wallet dir and coin manager. must gave password
func Init(walletDir, passwd string) error {
	wallet.InitDir(walletDir)
	if err := LoadWallet(passwd); err != nil {
		return err
	}
	coinMap = make(map[string]Coiner)
	return nil
}

// LoadWallet Load wallet already exists
func LoadWallet(passwd string) error {
	if len(passwd) == 0 {
		return errors.New("password cannot empty")
	}
	return wallet.LoadWallet(passwd)
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

// GetSupportedCoin return supported coins, joined by ","
func GetSupportedCoin() string {
	coinTypes := []string{skycoin.Type, samos.Type, spo.Type, suncoin.Type, shellcoin.Type, mzcoin.Type, aynrandcoin.Type, mdl.Type}
	return strings.Join(coinTypes, ",")
}

// NewWallet create a new wallet base on the wallet type and seed
func NewWallet(coinType, lable, seed, passwd string) (string, error) {
	if len(passwd) == 0 {
		return "", errors.New("password cannot empty")
	}
	wlt, err := wallet.New(coinType, lable, seed, passwd)
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
func NewAddress(walletID string, num int, passwd string) (string, error) {
	if len(passwd) == 0 {
		return "", errors.New("password cannot empty")
	}
	es, err := wallet.NewAddresses(walletID, num, passwd)
	if err != nil {
		return "", err
	}
	tempes := []coin.AddressEntry{}
	for _, ee := range es {
		ee.Secret = ""
		tempes = append(tempes, ee)
	}
	var res = struct {
		Entries []coin.AddressEntry `json:"addresses"`
	}{
		tempes,
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
func GetKeyPairOfAddr(walletID, addr, passwd string) (string, error) {
	p, s, err := wallet.GetKeypair(walletID, addr, passwd)
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
// returns {"balance":"70.000000", "hours": "32001"}
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
	b := skywallet.BalancePair{}
	err = json.Unmarshal([]byte(bal), &b)
	if err != nil {
		return "", err
	}
	coins, err := droplet.ToString(b.Confirmed.Coins)
	if err != nil {
		return "", err
	}

	hours := int64(b.Confirmed.Hours)

	var res = struct {
		Balance string `json:"balance"`
		Hours   int64  `json:"hours"`
	}{
		coins,
		hours,
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
	b := skywallet.BalancePair{}
	err = json.Unmarshal([]byte(bal), &b)
	if err != nil {
		return "", err
	}
	coins, err := droplet.ToString(b.Confirmed.Coins)
	if err != nil {
		return "", err
	}

	hours := int64(b.Confirmed.Hours)
	var res = struct {
		Balance string `json:"balance"`
		Hours   int64  `json:"hours"`
	}{
		coins,
		hours,
	}

	d, err := json.Marshal(res)
	if err != nil {
		return "", err
	}
	return string(d), nil
}

// Send send coins, support bitcoin and all coins in skycoin ledger
func Send(coinType, wid, toAddr, amount, passwd string) (string, error) {
	coin, ok := coinMap[coinType]
	if !ok {
		return "", fmt.Errorf("%s is not supported", coinType)
	}

	return coin.Send(wid, toAddr, amount, passwd)
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
func GetSeed(walletID, passwd string) (string, error) {
	return wallet.GetSeed(walletID, passwd)
}

func getPrivateKey(walletID, passwd string) coin.GetPrivKey {
	return func(addr string) (string, error) {
		_, s, err := wallet.GetKeypair(walletID, addr, passwd)
		return s, err
	}
}
