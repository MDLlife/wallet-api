package mdl

import (
	skycoin "github.com/MDLlife/wallet-api/src/coin/skycoin"
	"github.com/MDLlife/wallet-api/src/wallet"
)

// Type represents MDL coin type
var Type = "mdl"

func init() {
	// Register wallet creator
	wallet.RegisterCreator(Type, func() wallet.Walleter {
		return &skycoin.Wallet{
			Wallet: wallet.Wallet{
				Type: Type,
			},
		}
	})
}

// MDL will implement coin.Gateway interface
type MDL struct {
	skycoin.Skycoin // embeded from skycoin , as all apis are the same as skycoin
}

// New creates a MDL instance.
func New(nodeAddr string) *MDL {
	return &MDL{Skycoin: skycoin.Skycoin{NodeAddress: nodeAddr}}
}

// Symbol returns the MDL symbol
func (s MDL) Symbol() string {
	return "MDL"
}

// Type returns MDL type
func (s MDL) Type() string {
	return Type
}
