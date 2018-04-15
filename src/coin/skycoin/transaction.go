package skycoin

import (
	"io"
	"io/ioutil"

	"github.com/skycoin/skycoin/src/cipher/encoder"
	sky "github.com/skycoin/skycoin/src/coin"
)

type Transaction struct {
	sky.Transaction
}

func (tx *Transaction) Serialize() ([]byte, error) {
	return tx.Transaction.Serialize(), nil
}

func (tx *Transaction) Deserialize(r io.Reader) error {
	d, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	if err := encoder.DeserializeRaw(d, tx); err != nil {
		return err
	}
	return nil
}
