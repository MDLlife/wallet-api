package mobile

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/skycoin/skycoin/src/util/file"
	"github.com/stretchr/testify/assert"
)

func TestWrongPassword(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), ".wallet1000")
	defer func() {
		if tmpDir == "~" || tmpDir == file.UserHome() || tmpDir == "/Users/liuguirong" {
			panic("cannot remove dir")
		}
		err := os.RemoveAll(tmpDir)
		assert.Nil(t, err)
	}()

	var err error
	wrongLengthPassword := "1234567890" //len 10
	err = Init(tmpDir, wrongLengthPassword)
	assert.Equal(t, "password length must 16", err.Error())

	rightPassword := "12345678abcdefgh" //len 16
	err = Init(tmpDir, rightPassword)
	assert.NoError(t, err)

	originSeed := "ab 12 57 xx yy zz hh oo"
	_, err = NewWallet("spo", "wronglable", originSeed, wrongLengthPassword)
	assert.Equal(t, "password length must 16", err.Error())

	wlt, err := NewWallet("spo", "rightlable", originSeed, rightPassword)
	assert.NoError(t, err)
	assert.Equal(t, "spo_rightlable", wlt)

	assert.True(t, IsExist(wlt))
	err = Remove(wlt)
	assert.NoError(t, err)

}

func TestMobileApi(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), ".wallet1000")
	fmt.Printf("tmpDir %s\n", tmpDir)
	defer func() {
		if tmpDir == "~" || tmpDir == file.UserHome() || tmpDir == "/Users/liuguirong" {
			panic("cannot remove dir")
		}
		err := os.RemoveAll(tmpDir)
		assert.Nil(t, err)
	}()
	var wlt string
	var err error

	coinTypes := GetSupportedCoin()
	assert.Equal(t, "skycoin,spo,suncoin,shellcoin,mzcoin,aynrandcoin", coinTypes)
	password := "12345678abcdefgh"
	err = Init(tmpDir, password)
	assert.NoError(t, err)
	err = RegisterNewCoin("spo", "127.0.0.1:8620")
	assert.NoError(t, err)
	err = RegisterNewCoin("skycoin", "127.0.0.1:6420")
	assert.NoError(t, err)

	originSeed := "abcd 1234 8909 bcde xmme adbn nw we hell world then at"
	wlt = "spo_lableandseed"
	if !IsExist(wlt) {
		wlt, err = NewWallet("spo", "lableandseed", originSeed, password)
		assert.NoError(t, err)
		assert.Equal(t, "spo_lableandseed", wlt)

		// wrong password
		_, err = NewWallet("spo", "lableandseed_123", originSeed, "12345678abcdabcd")
		assert.Error(t, err)
		assert.Equal(t, "wallet password incorrect", err.Error())

		addresses, err := NewAddress(wlt, 3, password)
		expectAddrs := "{\"addresses\":[{\"address\":\"3nfw5uwWtktbNbGdx5cNF4i4GRUqp53Rtr\",\"pubkey\":\"02ba3470b34ad121ae4ac8036d76ed33b80d03c2d43aca4ad3947220053af11969\",\"seckey\":\"a03dc39c34c1f715658de0e6ffae66c02f5871b578834b3b18882a73ccc8dad9\"},{\"address\":\"2fwZKXRU9PAQ7TRxVzj2MTE9uz9gvccLEGZ\",\"pubkey\":\"032979cd01374e1160cb4da6176e95ab4b0017a409a34ab121f3f76595c6d6459d\",\"seckey\":\"091296f3cc1eef7432797d9ad4f05fd240974ffb993b3d824a11f6091d9e0df4\"},{\"address\":\"27QMsG95g3u2rFnfqoJhYF7ZFJttx1ZQYg9\",\"pubkey\":\"02620ba4c261ce12210ca791ffc234a36119e9cc14071d3cd3a5934c98c5026a7b\",\"seckey\":\"b82299cbfd8983903d4a4c01bf279702eae93a210f8abfb870dc59dd779deb77\"}]}"
		assert.NoError(t, err)
		assert.Equal(t, expectAddrs, addresses)

		// wrong password
		_, err = NewAddress(wlt, 3, "1234abcd1234abcd")
		assert.Error(t, err)
		assert.Equal(t, "wallet password incorrect", err.Error())
	}
	addresses, err := GetAddresses(wlt)
	assert.NoError(t, err)
	expectAddresses := "{\"addresses\":[\"3nfw5uwWtktbNbGdx5cNF4i4GRUqp53Rtr\",\"2fwZKXRU9PAQ7TRxVzj2MTE9uz9gvccLEGZ\",\"27QMsG95g3u2rFnfqoJhYF7ZFJttx1ZQYg9\"]}"
	assert.Equal(t, expectAddresses, addresses)

	addr := "3nfw5uwWtktbNbGdx5cNF4i4GRUqp53Rtr"
	pair, err := GetKeyPairOfAddr(wlt, addr)
	assert.NoError(t, err)
	expectPair := "{\"pubkey\":\"02ba3470b34ad121ae4ac8036d76ed33b80d03c2d43aca4ad3947220053af11969\",\"seckey\":\"a03dc39c34c1f715658de0e6ffae66c02f5871b578834b3b18882a73ccc8dad9\"}"
	assert.Equal(t, expectPair, pair)

	seed1, err := GetSeed(wlt)
	assert.NoError(t, err)
	assert.Equal(t, originSeed, seed1)
}
