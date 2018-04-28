package integration

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/spolabs/wallet-api/src/api/mobile"
)

func TestAllFunction(t *testing.T) {
	wltType := "spo"
	wltKey := "QsyueWQWvKhqsPDj"
	wlt := wltType + "_" + wltKey
	password := "12"
	walletDir := "/tmp/wallets" //use default dir ~/.wallet-family
	newseed := "greate test"

	coinTypes := mobile.GetSupportedCoin()
	assert.Equal(t, coinTypes, "skycoin,spo,suncoin,shellcoin,mzcoin,aynrandcoin")
	err := mobile.Init(walletDir, password)
	assert.NoError(t, err)
	err = mobile.RegisterNewCoin("spo", "182.92.180.92:8620")
	assert.NoError(t, err)
	if !mobile.IsExist(wlt) {
		fmt.Printf("wallet not exists\n")
		lable := "lable"
		wlt, err = mobile.NewWallet(wltType, lable, newseed, password)
		assert.NoError(t, err)
		assert.Equal(t, wlt, wltType+"_"+wltKey)
		address, err := mobile.NewAddress(wlt, 2, password)
		assert.NoError(t, err)
		expectAddress := "{\"addresses\":[{\"address\":\"2LZJg6FzLsVuHVLNU7EqZzbJTZw2fatj7qe\",\"pubkey\":\"02145bd5b695ebbca98e08b5fc3bafa4dbdb08a97ce392025c5a0699868307db0b\",\"seckey\":\"\"},{\"address\":\"bGxidkUzEetshF7oFb7G9vY3FMo8k7fnNF\",\"pubkey\":\"0253f7e925098bbf07b5315355f3569ef7ac3f9ab5e533adf80d14cac5c0cae803\",\"seckey\":\"\"}]}"
		assert.Equal(t, expectAddress, address)
	}
	expectAddress := "{\"addresses\":[\"QsyueWQWvKhqsPDjt1BrdspXVaKayTXUtr\",\"2LZJg6FzLsVuHVLNU7EqZzbJTZw2fatj7qe\",\"bGxidkUzEetshF7oFb7G9vY3FMo8k7fnNF\"]}"
	addresses, err := mobile.GetAddresses(wlt)
	assert.NoError(t, err)
	assert.Equal(t, expectAddress, addresses)
	addr := "QsyueWQWvKhqsPDjt1BrdspXVaKayTXUtr"
	pair, err := mobile.GetKeyPairOfAddr(wlt, addr, password)
	assert.NoError(t, err)
	expectPair := "{\"pubkey\":\"02581c05040376a9acf7d854f4f1f6fd698b60e0c30309674d20dbd8fbbfc9975c\",\"seckey\":\"5727658c4df4c86a537ddb1564de7d9b363b38664ae2119009a483c0183f4eb1\"}"
	assert.Equal(t, expectPair, pair)

	balance, err := mobile.GetWalletBalance(wltType, wlt)
	assert.NoError(t, err)
	//expectBalance := "{"balance":"0.100000","hours":995}"
	assert.NotEmpty(t, balance)
	assert.Contains(t, balance, "balance")

	seed := mobile.NewSeed()
	assert.NotEmpty(t, seed)

	seed1, err := mobile.GetSeed(wlt, password)
	assert.NoError(t, err)
	assert.Equal(t, newseed, seed1)

	txid := "76752105025ba4a84ff0e1ebe2f4a6b1b0f4e27f39433582a5abc419a7fb60de"
	txinfo, err := mobile.GetTransactionByID("spo", txid)
	assert.NoError(t, err)
	assert.NotEmpty(t, txinfo)
	assert.Contains(t, txinfo, "9a8f82802581e9360405e6b93a18d6f5162392ad5a07d9f110bcdd380c97b814")
	txConfirm, err := mobile.IsTransactionConfirmed(wltType, txid)
	assert.NoError(t, err)
	assert.Equal(t, true, txConfirm)

	//sposrc := "QsyueWQWvKhqsPDjt1BrdspXVaKayTXUtr"
	destAddr := "hva72jTmjEdogG4RxNb9uAmDgM1MCfSnLk"
	_, err = mobile.Send(wltType, wlt, destAddr, "0.0", password)
	assert.EqualError(t, err, "can not send 0 coins")
	wrongDestAddr := "hva72jTmjEdogG4RxNb9uAmDgM1MCfSnLn"
	_, err = mobile.Send(wltType, wlt, wrongDestAddr, "0.1", password)
	assert.EqualError(t, err, "Invalid checksum")
	_, err = mobile.Send(wltType, wlt, destAddr, "0.1234", password)
	assert.EqualError(t, err, "invalid amount, too many decimal places")
	_, err = mobile.Send(wltType, wlt, destAddr, "0.5", password)
	assert.EqualError(t, err, "insufficient balance")
}
