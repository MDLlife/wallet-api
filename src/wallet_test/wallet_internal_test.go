package wallet_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/spolabs/wallet-api/src/coin/skycoin"
	_ "github.com/spolabs/wallet-api/src/coin/spo"
	"github.com/spolabs/wallet-api/src/wallet"
	"github.com/stretchr/testify/assert"
)

func TestNewWallet(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), ".wallet1000")
	wallet.Reset()
	wallet.InitDir(tmpDir)

	defer func() {
		err := os.RemoveAll(tmpDir)
		assert.Nil(t, err)
	}()

	// create wallets.
	password := "walletpasswd"
	testData := []struct {
		Type   string
		Seed   string
		Lable  string
		Passwd string
	}{
		{"skycoin", "seed1", "l1", password},
		{"skycoin", "seed2", "l2", password},
		{"skycoin", "seed3", "l3", password},
		{"spo", "seed1", "l1", password},
		{"spo", "seed2", "l2", password},
		{"spo", "seed3", "l3", password},
		{"spo", "seed4", "l4", password},
	}

	for _, d := range testData {
		wlt, err := wallet.New(d.Type, d.Lable, d.Seed, d.Passwd)
		if err != nil {
			fmt.Println(d.Type, " ", d.Lable, " ", d.Seed, " ", d.Passwd)
			t.Error(err)
			return
		}
		// seed field has set to empty
		assert.NotNil(t, wlt.Validate())

		assert.Equal(t, wlt.GetID(), wallet.MakeWltID(d.Type, d.Seed))
		assert.Equal(t, true, wallet.IsExist(wlt.GetID()))
		assert.Equal(t, wlt.GetType(), d.Type)
		// default address number 1
		assert.Equal(t, len(wlt.GetAddresses()), 1)
		assert.Equal(t, wlt.GetSeed(password), d.Seed)

		walletFile := filepath.Join(tmpDir, (wlt.GetID() + "." + wallet.Ext))
		if _, err := os.Stat(walletFile); os.IsNotExist(err) {
			t.Error("create wallet failed")
			return
		}
	}

	err := wallet.LoadWallet(password)
	assert.NoError(t, err)
	dir := wallet.GetWalletDir()
	assert.Equal(t, dir, tmpDir)

}
