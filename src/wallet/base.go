package wallet

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/spolabs/wallet-api/src/coin"
	"github.com/spolabs/wallet-api/src/util/encrypt"
)

const (
	metaEncrypted = "encrypted" // whether the wallet is encrypted
	metaVersion   = "version"   // wallet version
	metaInitSeed  = "init_seed" // wallet seed
	metaSeed      = "seed"      // seed for generating next address
)

// Wallet wallet struct
type Wallet struct {
	ID             string              `json:"id"`                // wallet id
	InitSeed       string              `json:"-"`                 // Init seed, used to recover the wallet.
	Seed           string              `json:"-"`                 // used to track the latset seed
	Lable          string              `json:"lable"`             // lable
	AddressEntries []coin.AddressEntry `json:"-"`                 // address entries.
	StoreEntries   []coin.AddressEntry `json:"entries,omitempty"` // address entries.
	Type           string              `json:"type"`              // wallet type
	Secrets        string              `json:"secrets"`
}

// GetID return wallet id.
func (wlt *Wallet) GetID() string {
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
func (wlt *Wallet) GetKeypair(addr string) (string, string, error) {
	for _, e := range wlt.AddressEntries {
		if e.Address == addr {
			return e.Public, e.Secret, nil
		}
	}
	return "", "", fmt.Errorf("%s addr does not exist in wallet", addr)
}

// Save save the wallet
func (wlt *Wallet) Save(w io.Writer, pwd string) error {
	metaMap := make(map[string]string)
	metaMap[metaSeed] = wlt.InitSeed
	metaMap[metaInitSeed] = wlt.InitSeed
	newEntryies := []coin.AddressEntry{}
	for _, entry := range wlt.AddressEntries {
		metaMap[entry.Address] = entry.Secret
		newEntryies = append(newEntryies, entry)
	}
	// secret set empty
	wlt.StoreEntries = []coin.AddressEntry{}
	for _, entry := range newEntryies {
		entry.Secret = ""
		wlt.StoreEntries = append(wlt.StoreEntries, entry)
	}

	secretsBinary, err := json.Marshal(metaMap)
	if err != nil {
		return err
	}

	sb, err := encrypt.Encrypt([]byte(pwd), string(secretsBinary))
	if err != nil {
		return err
	}
	wlt.Secrets = string(sb)
	d, err := json.MarshalIndent(wlt, "", "    ")
	if err != nil {
		return err
	}
	_, err = io.Copy(w, bytes.NewBuffer(d))
	return err
}

// Load load wallet from reader.
func (wlt *Wallet) Load(r io.Reader, pwd string) error {
	err := json.NewDecoder(r).Decode(wlt)
	if err != nil {
		return err
	}
	metaMapB, err := encrypt.Decrypt([]byte(pwd), wlt.Secrets)
	if err != nil {
		return err
	}
	metaMap := make(map[string]string)
	err = json.Unmarshal([]byte(metaMapB), &metaMap)
	if err != nil {
		return err
	}

	seed, ok := metaMap[metaSeed]
	if !ok {
		return errors.New("no seed")
	}
	initSeed, ok := metaMap[metaInitSeed]
	if !ok {
		return errors.New("no init seed")
	}
	wlt.Seed = seed
	wlt.InitSeed = initSeed
	for _, entry := range wlt.StoreEntries {
		secret, ok := metaMap[entry.Address]
		if !ok {
			return fmt.Errorf("address %s no secret", entry.Address)
		}
		newEntry := entry
		newEntry.Secret = secret
		wlt.AddressEntries = append(wlt.AddressEntries, newEntry)
	}
	return nil
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
func (wlt *Wallet) Copy() Wallet {
	return Wallet{
		ID:             wlt.ID,
		Lable:          wlt.Lable,
		AddressEntries: wlt.StoreEntries,
		Type:           wlt.Type,
	}
}
