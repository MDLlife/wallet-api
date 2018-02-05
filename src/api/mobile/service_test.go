package mobile

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ = func() int64 {
	t := time.Now().Unix()
	rand.Seed(t)
	return t
}()

func TestNewWallet(t *testing.T) {
	testCases := []struct {
		name        string
		wltName     string
		coinType    string
		seed        string
		lable       string
		expectWltID string
		expectErr   error
	}{
		{
			"create skycoin wallet",
			"firstwlt",
			"skycoin",
			"abc",
			"first",
			"firstwlt",
			nil,
		},
		{
			"create unknow wallet",
			"firstwlt",
			"unknow",
			"abcde",
			"second",
			"",
			errors.New("wallet name would conflict with existing wallet, renaming"),
		},
	}

	wltName := fmt.Sprintf(".wallet%d", rand.Int31n(100))
	tmpDir := filepath.Join(os.TempDir(), wltName)
	if err := os.MkdirAll(tmpDir, 0777); err != nil {
		return
	}
	fmt.Printf("dir:%+v\n", tmpDir)

	teardown := func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			panic(err)
		}
	}

	defer teardown()

	ms, err := NewMyService(tmpDir)
	require.NoError(t, err)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			wlt, err := ms.CreateWallet(tc.coinType, tc.wltName, tc.seed, tc.lable)
			assert.Equal(t, tc.expectErr, err)
			assert.Equal(t, tc.expectWltID, wlt.GetID())
		})
	}
}
