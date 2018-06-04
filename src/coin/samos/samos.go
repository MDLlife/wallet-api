package samos

import (
	skycoin "github.com/MDLlife/wallet-api/src/coin/skycoin"
	"github.com/MDLlife/wallet-api/src/wallet"
)

// Type represents samos coin type
var Type = "samos"

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

// Samos will implement coin.Gateway interface
type Samos struct {
	skycoin.Skycoin // embeded from skycoin , as all apis are the same as skycoin
}

// New creates a samos instance.
func New(nodeAddr string) *Samos {
	return &Samos{Skycoin: skycoin.Skycoin{NodeAddress: nodeAddr}}
}

// Symbol returns the samos symbol
func (s Samos) Symbol() string {
	return "SAMOS"
}

// Type returns samos type
func (s Samos) Type() string {
	return Type
}
