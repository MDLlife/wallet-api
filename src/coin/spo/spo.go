package spo

import (
	skycoin "github.com/spolabs/wallet-api/src/coin/skycoin"
	"github.com/spolabs/wallet-api/src/wallet"
)

// Type represents spo coin type
var Type = "spo"

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

// Spo will implement coin.Gateway interface
type Spo struct {
	skycoin.Skycoin // embeded from skycoin , as all apis are the same as skycoin
}

// New creates a spo instance.
func New(nodeAddr string) *Spo {
	return &Spo{Skycoin: skycoin.Skycoin{NodeAddress: nodeAddr}}
}

// Symbol returns the spo symbol
func (s Spo) Symbol() string {
	return "SPO"
}

// Type returns spo type
func (s Spo) Type() string {
	return Type
}
