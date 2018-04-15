package wallet

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/skycoin/skycoin/src/util/file"
	"github.com/spolabs/wallet-api/src/coin"
	bip39 "github.com/tyler-smith/go-bip39"
)

// Walleter interface, new wallet type can be supported if it fullfills this interface.
type Walleter interface {
	GetID() string                                          // get wallet id.
	SetID(id string)                                        // set wallet id.
	SetSeed(seed string)                                    // init the wallet seed.
	SetLable(lable string)                                  // set the wallet lable.
	SetConstant()                                           // set constant such as version, type
	GetType() string                                        // get the wallet coin type.
	SetTime(tm string)                                      // set the wallet created time.
	GetSeed(passwd string) string                           // get the wallet seed.
	Validate() error                                        // Validate wallet fields
	IsPasswordCorrect(passwd string) error                  // check password correct or not.
	Decryption(passwd string) error                         // decryption secrets for new address.
	NewAddresses(num int) ([]coin.AddressEntry, error)      // generate new addresses.
	GetAddresses() []string                                 // get all addresses in the wallet.
	GetKeypair(addr, passwd string) (string, string, error) // get pub/sec key pair of specific address
	Save(w io.Writer, passwd string) error                  // save the wallet.
	Load(r io.Reader, passwd string) error                  // load wallet from reader.
	Copy() Walleter                                         // copy of self, for thread safe.
}

// wltDir default wallet dir, wallet file name sturct: $type_$lable.wlt.
// example: spo_lable.wlt, skycoin_lable.wlt.
var wltDir = filepath.Join(file.UserHome(), ".wallet-family")

// Version represents the current wallet version
var WalletVersion = "0.1"
var WalletType = "deterministic"

// Ext wallet file extension name
var Ext = "wlt"

// Creator wallet creator.
type Creator func() Walleter

var gWalletCreators = make(map[string]Creator)

// RegisterCreator when new type wallet need to be supported,
// the wallet must provide Creator, and use this function to register the creator into system.
func RegisterCreator(tp string, ctor Creator) error {
	if _, ok := gWalletCreators[tp]; ok {
		return fmt.Errorf("%s wallet already registered", tp)
	}
	gWalletCreators[tp] = ctor
	return nil
}

// InitDir initialize the wallet file storage dir,
// load wallets if exist.
func InitDir(path string) {
	if path == "" {
		path = wltDir
	} else {
		wltDir = path
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		//create the dir.
		if err := os.MkdirAll(path, 0777); err != nil {
			panic(err)
		}
	}

}

// LoadWallet load wallet from disk
func LoadWallet(passwd string) error {
	// load wallets.
	return gWallets.mustLoad(passwd)
}

// GetWalletDir return the current wallet dir.
func GetWalletDir() string {
	return wltDir
}

// New create wallet base on seed and coin type.
func New(tp, lable, seed, passwd string) (Walleter, error) {
	if gWallets.verifyPassword(passwd) != nil {
		return nil, fmt.Errorf("wallet password incorrect")
	}
	newWlt, ok := gWalletCreators[tp]
	if !ok {
		return nil, fmt.Errorf("%s wallet not regestered", tp)
	}

	// create wallet base on the wallet creator.
	wlt := newWlt()

	wlt.SetConstant()

	wlt.SetTime(fmt.Sprintf("%v", time.Now().Unix()))

	wlt.SetLable(lable)

	if seed == "" {
		seed = NewSeed()
	}

	wlt.SetID(MakeWltID(tp, seed))

	wlt.SetSeed(seed)

	// Validate the wallet
	if err := wlt.Validate(); err != nil {
		return nil, err
	}

	if err := gWallets.add(wlt, passwd); err != nil {
		return nil, err
	}

	// generate 1 address default
	if _, err := gWallets.newAddresses(wlt.GetID(), 1, passwd); err != nil {
		return nil, err
	}

	return wlt.Copy(), nil
}

// NewSeed generates mnemonic seed
func NewSeed() string {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		panic(err)
	}

	sd, err := bip39.NewMnemonic(entropy)
	if err != nil {
		panic(err)
	}
	return sd
}

// IsExist check if the wallet is already exist.
func IsExist(id string) bool {
	return gWallets.isExist(id)
}

// MakeWltID make wallet id base on coin type and md5(seed)[0:12]
func MakeWltID(cp, seed string) string {
	md5Value := md5.Sum([]byte(seed))
	return fmt.Sprintf("%s_%x", cp, md5Value[0:12])
}

// NewAddresses create address
func NewAddresses(id string, num int, passwd string) ([]coin.AddressEntry, error) {
	if gWallets.verifyPassword(passwd) != nil {
		return nil, fmt.Errorf("wallet password incorrect")
	}
	return gWallets.newAddresses(id, num, passwd)
}

// GetAddresses get all addresses in specific wallet.
func GetAddresses(id string) ([]string, error) {
	return gWallets.getAddresses(id)
}

// GetSeed get seed in specific wallet.
func GetSeed(id, passwd string) (string, error) {
	return gWallets.getSeed(id, passwd)
}

// IsContain check if the addresses are int the wallet.
func IsContain(id string, addrs []string) (bool, error) {
	return gWallets.isContain(id, addrs)
}

// GetKeypair get pub/sec key pair of specific addresse in wallet.
func GetKeypair(id, addr, passwd string) (string, string, error) {
	return gWallets.getKeypair(id, addr, passwd)
}

// Remove remove wallet of specific id.
func Remove(id string) error {
	return gWallets.remove(id)
}
