package mzcoin

import (
	skycoin "github.com/spolabs/wallet-api/src/coin/skycoin"
	"github.com/spolabs/wallet-api/src/wallet"
)

// Type represents mzcoin coin type
var Type = "mzcoin"

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

// Mzcoin will implement coin.Gateway interface
type Mzcoin struct {
	skycoin.Skycoin // embeded from skycoin , as all apis are the same as skycoin
}

// New creates a mzcoin instance.
func New(nodeAddr string) *Mzcoin {
	return &Mzcoin{Skycoin: skycoin.Skycoin{NodeAddress: nodeAddr}}
}

// Symbol returns the mzcoin symbol
func (mz Mzcoin) Symbol() string {
	return "MZC"
}

// Type returns mzcoin type
func (mz Mzcoin) Type() string {
	return Type
}
