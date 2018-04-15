package skycoin

import (
	"fmt"

	logging "github.com/op/go-logging"
	"github.com/skycoin/skycoin/src/cipher"
	sky "github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/spolabs/wallet-api/src/coin"
)

var (
	// ServeAddr  string = "127.0.0.1:6420"
	logger = logging.MustGetLogger("exchange.skycoin")
	// Type returns the coin type
	Type = "skycoin"
)

// Skycoin skycoin gateway.
type Skycoin struct {
	NodeAddress string // skycoin node address
}

// New creates a skycoin instance.
func New(nodeAddr string) *Skycoin {
	return &Skycoin{NodeAddress: nodeAddr}
}

// SkyUtxo skycoin utxo struct
type SkyUtxo struct {
	visor.ReadableOutput
}

// TxOut transaction output filed
type TxOut struct {
	sky.TransactionOutput
}

// GenerateAddresses generate addresses.
func GenerateAddresses(seed []byte, num int) (string, []coin.AddressEntry) {
	sd, seckeys := cipher.GenerateDeterministicKeyPairsSeed(seed, num)
	entries := make([]coin.AddressEntry, num)
	for i, sec := range seckeys {
		pub := cipher.PubKeyFromSecKey(sec)
		entries[i].Address = cipher.AddressFromPubKey(pub).String()
		entries[i].Public = pub.Hex()
		entries[i].Secret = sec.Hex()
	}
	return fmt.Sprintf("%2x", sd), entries
}

// Symbol returns skycoin sybmol
func (sky *Skycoin) Symbol() string {
	return "SKY"
}

// Type returns skycoin type name
func (sky *Skycoin) Type() string {
	return Type
}
