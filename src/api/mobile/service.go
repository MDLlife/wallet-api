package mobile

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/wallet"
)

type MyService struct {
	serv *wallet.Service
}

func NewMyService(walletDir string) (*MyService, error) {
	var err error
	ms := MyService{}
	ms.serv, err = wallet.NewService(walletDir)
	if err != nil {
		return nil, err
	}
	return &ms, nil
}

func (ms *MyService) GetWallets() wallet.Wallets {
	return ms.serv.GetWallets()
}

func (ms *MyService) CreateWallet(coinType, wltName, seed, lable string) (wallet.Wallet, error) {
	return ms.serv.CreateWallet(wltName, wallet.Options{Seed: seed, Label: lable, Coin: wallet.CoinType(coinType)})
}

func (ms *MyService) GetWallet(wltName string) (wallet.Wallet, error) {
	return ms.GetWallet(wltName)
}

func (ms *MyService) NewAddresses(wltName string, num uint64) ([]cipher.Address, error) {
	return ms.serv.NewAddresses(wltName, num)
}

func (ms *MyService) GetAddresses(wltName string) ([]cipher.Address, error) {
	return ms.serv.GetAddresses(wltName)
}

func (ms *MyService) ReloadWallets() error {
	return ms.serv.ReloadWallets()
}

func (ms *MyService) GetWalletsReadable() []*wallet.ReadableWallet {
	return ms.serv.GetWalletsReadable()
}

func (ms *MyService) UpdateWalletLabel(wltID, lable string) error {
	return ms.serv.UpdateWalletLabel(wltID, lable)
}
