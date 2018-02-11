package wallet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/spolabs/spo/src/util/encrypt"
	"github.com/spolabs/wallet-api/src/coin"
)

// Wallet wallet struct
type Wallet struct {
	ID             string              `json:"id"`                // wallet id
	InitSeed       string              `json:"init_seed"`         // Init seed, used to recover the wallet.
	Seed           string              `json:"seed"`              // used to track the latset seed
	Lable          string              `json:"lable"`             // lable
	AddressEntries []coin.AddressEntry `json:"entries,omitempty"` // address entries.
	Type           string              `json:"type"`              // wallet type
}

// DiskWallet wallet struct for disk store
type DiskWallet struct {
	ID             string             `json:"id"`                // wallet id
	InitSeed       string             `json:"init_seed"`         // Init seed, used to recover the wallet.
	Seed           string             `json:"seed"`              // used to track the latset seed
	Lable          string             `json:"lable"`             // lable
	AddressEntries []DiskAddressEntry `json:"entries,omitempty"` // address entries.
	Type           string             `json:"type"`              // wallet type
}

// DiskAddressEntry represents the wallet address
type DiskAddressEntry struct {
	Address string `json:"address"`
	Public  string `json:"pubkey"`
	Secret  string `json:"seckey"`
}

func (dw DiskWallet) toWallet(pwd string) (*Wallet, error) {
	wlt := &Wallet{}
	wlt.ID = dw.ID
	wlt.Lable = dw.Lable
	wlt.Type = dw.Type
	var err error
	wlt.Seed, err = encrypt.Decrypt([]byte(pwd), dw.Seed)
	if err != nil {
		return nil, err
	}
	wlt.InitSeed, err = encrypt.Decrypt([]byte(pwd), dw.InitSeed)
	if err != nil {
		return nil, err
	}
	wlt.AddressEntries, err = recoverAddressEntry(dw.AddressEntries, pwd)
	if err != nil {
		return nil, err
	}

	return wlt, err
}

func (wlt Wallet) toDiskWallet(pwd string) (*DiskWallet, error) {
	dw := &DiskWallet{}
	dw.ID = wlt.ID
	dw.Lable = wlt.Lable
	dw.Type = wlt.Type
	var err error
	dw.InitSeed, err = encrypt.Encrypt([]byte(pwd), wlt.InitSeed)
	if err != nil {
		return nil, err
	}
	dw.Seed, err = encrypt.Encrypt([]byte(pwd), wlt.Seed)
	if err != nil {
		return nil, err
	}
	dw.AddressEntries, err = toDiskAddressEntry(wlt.AddressEntries, pwd)
	if err != nil {
		return nil, err
	}
	return dw, nil
}

func toDiskAddressEntry(entrys []coin.AddressEntry, pwd string) ([]DiskAddressEntry, error) {
	diskAddresses := []DiskAddressEntry{}
	for _, entry := range entrys {
		addr := DiskAddressEntry{}
		addr.Address = entry.Address
		addr.Public = entry.Public
		var err error
		addr.Secret, err = encrypt.Encrypt([]byte(pwd), entry.Secret)
		if err != nil {
			return []DiskAddressEntry{}, nil
		}
		diskAddresses = append(diskAddresses, addr)
	}
	return diskAddresses, nil
}

// GetID return wallet id.
func (wlt Wallet) GetID() string {
	return wlt.ID
}

// SetID set wallet id
func (wlt *Wallet) SetID(id string) {
	wlt.ID = id
}

// SetSeed initialize the wallet seed.
func (wlt *Wallet) SetSeed(seed string) {
	wlt.InitSeed = seed
	wlt.Seed = seed
}

// SetLable set wallet lable.
func (wlt *Wallet) SetLable(lable string) {
	wlt.Lable = lable
}

// GetAddresses return all addresses in wallet.
func (wlt *Wallet) GetAddresses() []string {
	addrs := []string{}
	for _, e := range wlt.AddressEntries {
		addrs = append(addrs, e.Address)
	}
	return addrs
}

// GetKeypair get pub/sec key pair of specific address
func (wlt Wallet) GetKeypair(addr string) (string, string, error) {
	for _, e := range wlt.AddressEntries {
		if e.Address == addr {
			return e.Public, e.Secret, nil
		}
	}
	return "", "", fmt.Errorf("%s addr does not exist in wallet", addr)
}

// Save save the wallet
func (wlt *Wallet) Save(w io.Writer, pwd string) error {
	diskWlt, err := wlt.toDiskWallet(pwd)
	if err != nil {
		return err
	}
	d, err := json.MarshalIndent(diskWlt, "", "    ")
	if err != nil {
		return err
	}
	_, err = io.Copy(w, bytes.NewBuffer(d))
	return err
}

// Load load wallet from reader.
func (wlt *Wallet) Load(r io.Reader, pwd string) error {
	dw := &DiskWallet{}
	err := json.NewDecoder(r).Decode(dw)
	if err != nil {
		return err
	}
	wlt1, err := dw.toWallet(pwd)
	if err != nil {
		return err
	}
	wlt.ID = wlt1.ID
	wlt.Type = wlt1.Type
	wlt.Seed = wlt1.Seed
	wlt.InitSeed = wlt1.InitSeed
	wlt.AddressEntries = wlt1.AddressEntries
	wlt.Lable = wlt1.Lable
	return nil
}

func recoverAddressEntry(entrys []DiskAddressEntry, pwd string) ([]coin.AddressEntry, error) {
	addresses := []coin.AddressEntry{}
	for _, entry := range entrys {
		addr := coin.AddressEntry{}
		addr.Address = entry.Address
		addr.Public = entry.Public
		var err error
		addr.Secret, err = encrypt.Decrypt([]byte(pwd), entry.Secret)
		if err != nil {
			return []coin.AddressEntry{}, nil
		}
		addresses = append(addresses, addr)
	}
	return addresses, nil
}

// GetType returns the wallet type
func (wlt *Wallet) GetType() string {
	return wlt.Type
}

// GetSeed returns the wallet seed
func (wlt *Wallet) GetSeed() string {
	return wlt.InitSeed
}

// Copy return the copy of self, for thread safe.
func (wlt Wallet) Copy() Wallet {
	return Wallet{
		ID:             wlt.ID,
		InitSeed:       wlt.InitSeed,
		Seed:           wlt.Seed,
		Lable:          wlt.Lable,
		AddressEntries: wlt.AddressEntries,
	}
}
