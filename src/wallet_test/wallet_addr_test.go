package wallet_test

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spolabs/wallet-api/src/coin"
	"github.com/spolabs/wallet-api/src/wallet"
	"github.com/stretchr/testify/assert"
)

// set rand seed.
var _ = func() int64 {
	t := time.Now().Unix()
	rand.Seed(t)
	return t
}()

func setup(t *testing.T) (string, func(), error) {
	wltName := fmt.Sprintf(".wallet%d", rand.Int31n(100))
	teardown := func() {}
	tmpDir := filepath.Join(os.TempDir(), wltName)
	if err := os.MkdirAll(tmpDir, 0777); err != nil {
		return "", teardown, err
	}

	teardown = func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			panic(err)
		}
	}
	wallet.InitDir(tmpDir)
	return tmpDir, teardown, nil
}

func TestNewAddresses(t *testing.T) {
	wltDir, teardown, err := setup(t)
	// wltDir, _, err := setup(t)
	assert.Nil(t, err)
	defer teardown()
	password := "12345678"
	wallet.Reset()
	testData := []struct {
		Type    string
		Seed    string
		Lable   string
		Num     int
		Entries []coin.AddressEntry
	}{
		{
			Type:  "spo",
			Seed:  "sd999",
			Lable: "l1",
			Num:   2,
			Entries: []coin.AddressEntry{
				{
					Address: "NKBCkfv4NW6MkvGeKjdjc8H6CKsE98KFf3",
					Public:  "0378c76e20e4f93730e67bb469bc7186681a8c85023088b64c70930e78d4aff690",
					Secret:  "",
				},
				{
					Address: "21nBJKVdpP6cLDVDJRMJmeQvtEKrWZ7fy1J",
					Public:  "0270d2d9b6df46e1b22effee8a3dfb42f6c3fe69b4361158b6101b451f6cced51c",
					Secret:  "",
				},
			},
		},
		{
			Type:  "skycoin",
			Seed:  "sd888",
			Lable: "l2",
			Num:   2,
			Entries: []coin.AddressEntry{
				{
					Address: "fYJPkCTqdChw3sPSGUgze9nuGMNtC5DvPY",
					Public:  "02ba572a03c8471822c308e5d041aba549b35676a0ef1c737b4517eef70c32377e",
					Secret:  "",
				},
				{
					Address: "t6t7bJ9Ruxq9z44pYQT5AkEeAjGjgantU",
					Public:  "039f4b6a110a9c5c38da08a0bff133edf07472348a4dc4c9d63b178fe26807606e",
					Secret:  "",
				},
			},
		},
	}

	for _, d := range testData {
		// new wallet
		wlt, err := wallet.New(d.Type, d.Lable, d.Seed, password)
		assert.Nil(t, err)

		// has 1 address default
		for i := 0; i < d.Num-1; i++ {
			if _, err := wallet.NewAddresses(wlt.GetID(), 1, password); err != nil {
				t.Fatal(err)
			}
		}
		path := filepath.Join(wltDir, fmt.Sprintf("%s.%s", wlt.GetID(), wallet.Ext))
		cnt, err := ioutil.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		for _, e := range d.Entries {
			if !strings.Contains(string(cnt), e.Address) {
				t.Fatalf("not contains address:%s", e.Address)
			}

			if !strings.Contains(string(cnt), e.Public) {
				t.Fatalf("not cointains pubkey:%s", e.Public)
			}
		}
	}
}

func TestGetAddresses(t *testing.T) {
	_, teardown, err := setup(t)
	assert.Nil(t, err)
	defer teardown()
	password := "12345678"
	wallet.Reset()
	testData := []struct {
		Type    string
		Seed    string
		Lable   string
		Num     int
		Entries []coin.AddressEntry
	}{
		{
			Type:  "spo",
			Seed:  "sd999",
			Lable: "l1",
			Num:   2,
			Entries: []coin.AddressEntry{
				{
					Address: "NKBCkfv4NW6MkvGeKjdjc8H6CKsE98KFf3",
					Public:  "0378c76e20e4f93730e67bb469bc7186681a8c85023088b64c70930e78d4aff690",
					Secret:  "",
				},
				{
					Address: "21nBJKVdpP6cLDVDJRMJmeQvtEKrWZ7fy1J",
					Public:  "0270d2d9b6df46e1b22effee8a3dfb42f6c3fe69b4361158b6101b451f6cced51c",
					Secret:  "",
				},
			},
		},
		{
			Type:  "skycoin",
			Seed:  "sd888",
			Lable: "l2",
			Num:   2,
			Entries: []coin.AddressEntry{
				{
					Address: "fYJPkCTqdChw3sPSGUgze9nuGMNtC5DvPY",
					Public:  "02ba572a03c8471822c308e5d041aba549b35676a0ef1c737b4517eef70c32377e",
					Secret:  "",
				},
				{
					Address: "t6t7bJ9Ruxq9z44pYQT5AkEeAjGjgantU",
					Public:  "039f4b6a110a9c5c38da08a0bff133edf07472348a4dc4c9d63b178fe26807606e",
					Secret:  "",
				},
			},
		},
	}

	for _, d := range testData {
		// new wallet
		wlt, err := wallet.New(d.Type, d.Lable, d.Seed, password)
		if err != nil {
			t.Fatal(err)
		}

		if _, err := wallet.NewAddresses(wlt.GetID(), d.Num-1, password); err != nil {
			t.Fatal(err)
		}

		addrs, err := wallet.GetAddresses(wlt.GetID())
		if err != nil {
			t.Fatal(err)
		}

		for _, e := range d.Entries {
			find := func(addr string) bool {
				for _, a := range addrs {
					if a == addr {
						return true
					}
				}
				return false
			}
			if !find(e.Address) {
				t.Fatal("GetAddresses failed")
			}
		}
	}
}

func TestGetKeypair(t *testing.T) {
	_, teardown, err := setup(t)
	assert.Nil(t, err)
	defer teardown()
	password := "12345678"
	wallet.Reset()
	testData := []struct {
		Type    string
		Seed    string
		Lable   string
		Num     int
		Entries []coin.AddressEntry
	}{
		{
			Type:  "spo",
			Seed:  "sd999",
			Lable: "l1",
			Num:   2,
			Entries: []coin.AddressEntry{
				{
					Address: "NKBCkfv4NW6MkvGeKjdjc8H6CKsE98KFf3",
					Public:  "0378c76e20e4f93730e67bb469bc7186681a8c85023088b64c70930e78d4aff690",
					Secret:  "ddfb2c79a00f39e006a3b51af09c019502a35ecddbd698c5abfd4f619eacf149",
				},
				{
					Address: "21nBJKVdpP6cLDVDJRMJmeQvtEKrWZ7fy1J",
					Public:  "0270d2d9b6df46e1b22effee8a3dfb42f6c3fe69b4361158b6101b451f6cced51c",
					Secret:  "578fb0706b56ea42416a6b28e1f509eb40d8becd7e3ae3fc4c395998d189b4e5",
				},
			},
		},
		{
			Type:  "skycoin",
			Seed:  "sd888",
			Lable: "l2",
			Num:   2,
			Entries: []coin.AddressEntry{
				{
					Address: "fYJPkCTqdChw3sPSGUgze9nuGMNtC5DvPY",
					Public:  "02ba572a03c8471822c308e5d041aba549b35676a0ef1c737b4517eef70c32377e",
					Secret:  "2f4aacc72a6d192e04ec540328689588caf4167d71904bdb870a4a2cee7f29c8",
				},
				{
					Address: "t6t7bJ9Ruxq9z44pYQT5AkEeAjGjgantU",
					Public:  "039f4b6a110a9c5c38da08a0bff133edf07472348a4dc4c9d63b178fe26807606e",
					Secret:  "b720d3c0f67f3c91e23805237f182e78121b90890f483133cc46f9d91232cf4c",
				},
			},
		},
	}

	for _, d := range testData {
		// new wallet
		wlt, err := wallet.New(d.Type, d.Lable, d.Seed, password)
		if err != nil {
			t.Fatal(err)
		}

		if _, err := wallet.NewAddresses(wlt.GetID(), d.Num-1, password); err != nil {
			t.Fatal(err)
		}

		for _, e := range d.Entries {
			p, s, err := wallet.GetKeypair(wlt.GetID(), e.Address, password)
			if err != nil {
				t.Fatal(err)
			}
			if p != e.Public || s != e.Secret {
				t.Fatal("get key pair failed")
			}
		}
	}
}
