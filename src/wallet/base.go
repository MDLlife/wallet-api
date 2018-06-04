package wallet

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/MDLlife/wallet-api/src/coin"
	"github.com/MDLlife/wallet-api/src/util/encrypt"
)

const (
	metaEncrypted = "encrypted" // whether the wallet is encrypted
	metaVersion   = "version"   // wallet version
	metaInitSeed  = "init_seed" // wallet seed
	metaSeed      = "seed"      // seed for generating next address
)

// Wallet wallet struct
type Wallet struct {
	Version        string              `json:"version"`           // version
	ID             string              `json:"id"`                // wallet id
	InitSeed       string              `json:"init_seed"`         // Init seed, used to recover the wallet.
	Seed           string              `json:"seed"`              // used to track the latset seed
	Lable          string              `json:"lable"`             // lable
	AddressEntries []coin.AddressEntry `json:"entries,omitempty"` // address entries.
	Type           string              `json:"type"`              // wallet type
	Tm             string              `json:"tm"`
	WalletType     string              `json:"wallet_type"`
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

// SetConstant set wallet version and type.
func (wlt *Wallet) SetConstant() {
	wlt.Version = WalletVersion
	wlt.WalletType = WalletType
}

// SetTime set wallet created time.
func (wlt *Wallet) SetTime(tm string) {
	wlt.Tm = tm
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
func (wlt *Wallet) GetKeypair(addr, passwd string) (string, string, error) {
	if err := wlt.Decryption(passwd); err != nil {
		return "", "", err
	}

	defer wlt.erase()

	for _, e := range wlt.AddressEntries {
		if e.Address == addr {
			return e.Public, e.Secret, nil
		}
	}
	return "", "", fmt.Errorf("%s addr does not exist in wallet", addr)
}

// Save save the wallet
func (wlt *Wallet) Save(w io.Writer, passwd string) error {
	if err := wlt.Encryption(passwd); err != nil {
		return err
	}
	d, err := json.MarshalIndent(wlt, "", "    ")
	if err != nil {
		return err
	}
	_, err = io.Copy(w, bytes.NewBuffer(d))
	return err
}

// Encryption encryption seed and private key.
func (wlt *Wallet) Encryption(passwd string) error {
	if wlt.Seed == "" || wlt.InitSeed == "" {
		return errors.New("empty seed")
	}
	// temporary map
	metaMap := make(map[string]string)
	metaMap[metaSeed] = wlt.Seed
	metaMap[metaInitSeed] = wlt.InitSeed
	for _, entry := range wlt.AddressEntries {
		metaMap[entry.Address] = entry.Secret
	}

	// delete critical info
	defer func() {
		for k := range metaMap {
			metaMap[k] = ""
			delete(metaMap, k)
		}
	}()

	secretsBinary, err := json.Marshal(metaMap)
	if err != nil {
		return err
	}

	sb, err := encrypt.Encrypt([]byte(passwd), string(secretsBinary))
	if err != nil {
		return err
	}
	wlt.Secrets = string(sb)

	wlt.erase()

	return nil
}

// Load load wallet from reader.
func (wlt *Wallet) Load(r io.Reader, passwd string) error {
	err := json.NewDecoder(r).Decode(wlt)
	if err != nil {
		return err
	}
	return wlt.Decryption(passwd)
}

// Validate validates the wallet
func (wlt *Wallet) Validate() error {
	if wlt.ID == "" {
		return errors.New("wallet id not set")
	}

	if wlt.Seed == "" {
		return errors.New("seed field not set")
	}

	if wlt.Type == "" {
		return errors.New("type field not set")
	}

	if wlt.WalletType != "deterministic" {
		return errors.New("wallet type invalid")
	}

	return nil
}

// IsPasswordCorrect check password correct or not.
func (wlt *Wallet) IsPasswordCorrect(passwd string) (err error) {
	// first
	if wlt.Secrets == "" {
		return nil
	}
	_, err = encrypt.Decrypt([]byte(passwd), wlt.Secrets)
	return
}

// Decryption decryption wallet recover seed and private key.
func (wlt *Wallet) Decryption(passwd string) error {
	if wlt.Secrets == "" {
		return errors.New("secrets is empty")
	}
	metaMapB, err := encrypt.Decrypt([]byte(passwd), wlt.Secrets)
	if err != nil {
		return err
	}
	metaMap := make(map[string]string)
	err = json.Unmarshal([]byte(metaMapB), &metaMap)
	if err != nil {
		return err
	}
	defer func() {
		for k := range metaMap {
			metaMap[k] = ""
			delete(metaMap, k)
		}
	}()

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
	for i, entry := range wlt.AddressEntries {
		secret, ok := metaMap[entry.Address]
		if !ok {
			return fmt.Errorf("address %s no secret", entry.Address)
		}
		wlt.AddressEntries[i].Secret = secret
	}
	return nil
}

func (wlt *Wallet) String() string {
	s1 := fmt.Sprintf("\t[seed:%s\n\tinit_seed:%s\n", wlt.Seed, wlt.InitSeed)
	s1 += fmt.Sprintf("\taddress number:%d\n", len(wlt.AddressEntries))
	for _, entry := range wlt.AddressEntries {
		s1 += fmt.Sprintf("\t\taddress:%s\n", entry.Address)
		s1 += fmt.Sprintf("\t\tpubkey:%s\n", entry.Public)
		s1 += fmt.Sprintf("\t\tseckey:%s\n", entry.Secret)
		s1 += "\n"
	}
	return s1 + "\t]\n"
}

// erase erase critical field such as seed, private key
func (wlt *Wallet) erase() {
	wlt.Seed = ""
	wlt.InitSeed = ""
	for i := range wlt.AddressEntries {
		wlt.AddressEntries[i].Secret = ""
	}
}

// GetType returns the wallet type
func (wlt *Wallet) GetType() string {
	return wlt.Type
}

// GetSeed returns the wallet seed
func (wlt *Wallet) GetSeed(passwd string) string {
	if err := wlt.Decryption(passwd); err != nil {
		return ""
	}
	defer wlt.erase()
	return wlt.InitSeed
}

// Copy return the copy of self, for thread safe.
func (wlt *Wallet) Copy() Wallet {
	return Wallet{
		ID:             wlt.ID,
		Lable:          wlt.Lable,
		AddressEntries: wlt.AddressEntries,
		Tm:             wlt.Tm,
		WalletType:     wlt.WalletType,
		Version:        wlt.Version,
		Type:           wlt.Type,
		Secrets:        wlt.Secrets,
	}
}
