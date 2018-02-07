package skycoin

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"io/ioutil"

	logging "github.com/op/go-logging"
	"github.com/skycoin/skycoin/src/cipher"
	sky "github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util/droplet"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/spolabs/wallet-api/src/coin"
)

var (
	// HideSeckey
	HideSeckey = false
	// ServeAddr  string = "127.0.0.1:6420"
	logger = logging.MustGetLogger("exchange.skycoin")
	// Type returns the coin type
	Type = "skycoin"
)

// Skycoin skycoin gateway.
type Skycoin struct {
	NodeAddress string // skycoin node address
}

// New creates a skycoin instance.
func New(nodeAddr string) *Skycoin {
	return &Skycoin{NodeAddress: nodeAddr}
}

// GetTx get skycoin verbose transaction.
func (sky *Skycoin) GetTx(txid string) (string, error) {
	url := fmt.Sprintf("http://%s/transaction?txid=%s", sky.NodeAddress, txid)
	rsp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()
	d, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}

	if rsp.StatusCode != 200 {
		return "", errors.New(string(d))
	}

	return string(d), nil
}

// GetRawTx get raw tx by txid.
func (sky *Skycoin) GetRawTx(txid string) (string, error) {
	url := fmt.Sprintf("http://%s/rawtx?txid=%s", sky.NodeAddress, txid)
	rsp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()
	s, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}
	return strings.Trim(string(s), "\""), nil
}

// InjectTx inject skycoin transaction.
func (sky *Skycoin) InjectTx(rawtx string) (string, error) {
	return BroadcastTx(sky.NodeAddress, rawtx)
}

// GetBalance get skycoin balance of specific addresses.
func (sky *Skycoin) GetBalance(addrs []string) (string, error) {
	url := fmt.Sprintf("http://%s/balance?addrs=%s", sky.NodeAddress, strings.Join(addrs, ","))
	rsp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()
	return "", nil
}

// ValidateTxid verify the valiation of specific transaction id.
func (sky *Skycoin) ValidateTxid(txid string) bool {
	_, err := cipher.SHA256FromHex(txid)
	return err == nil
}

// CreateRawTx create skycoin raw transaction.
func (sky Skycoin) CreateRawTx(txIns []coin.TxIn, txOuts interface{}) (string, error) {
	tx := Transaction{}
	// keys := make([]cipher.SecKey, len(utxos))
	for _, in := range txIns {
		tx.PushInput(cipher.MustSHA256FromHex(in.Txid))
	}

	s := reflect.ValueOf(txOuts)
	if s.Kind() != reflect.Slice {
		return "", errors.New("error tx out type")
	}
	outs := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		outs[i] = s.Index(i).Interface()
	}

	if len(outs) > 2 {
		return "", errors.New("out address more than 2")
	}

	for _, o := range outs {
		out := o.(TxOut)
		if (out.Coins % 1e6) != 0 {
			return "", errors.New("skycoin coins must be multiple of 1e6")
		}
		tx.PushOutput(out.Address, out.Coins, out.Hours)
	}

	tx.UpdateHeader()
	d, err := tx.Serialize()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(d), nil
}

// SignRawTx sign skycoin transaction.
func (sky Skycoin) SignRawTx(rawtx string, getKey coin.GetPrivKey) (string, error) {
	// decode the rawtx
	tx := Transaction{}
	b, err := hex.DecodeString(rawtx)
	if err != nil {
		return "", err
	}
	if err := tx.Deserialize(bytes.NewBuffer(b)); err != nil {
		return "", err
	}

	// TODO: need to get the address of the in hash, then get key of those address, and sign.
	hashes := make([]string, len(tx.In))
	for i, in := range tx.In {
		hashes[i] = in.Hex()
	}

	// get utxos of thoes hashes.
	utxos, err := getUnspentOutputsByHashes(sky.NodeAddress, hashes)
	if err != nil {
		return "", err
	}

	if len(utxos) != len(hashes) {
		return "", errors.New("failed to search tx in's address")
	}

	hashAddrMap := map[string]string{}
	for _, u := range utxos {
		hashAddrMap[u.GetHash()] = u.GetAddress()
	}

	keys := make([]cipher.SecKey, len(hashes))
	for i, h := range hashes {
		key, err := getKey(hashAddrMap[h])
		if err != nil {
			return "", err
		}

		keys[i], err = cipher.SecKeyFromHex(key)
		if err != nil {
			return "", err
		}
	}

	tx.SignInputs(keys)
	tx.UpdateHeader()
	d, err := tx.Serialize()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(d), nil
}

// Utxo unspent outputs interface
type Utxo interface {
	GetHash() string
	GetSrcTx() string
	GetAddress() string
	GetCoins() uint64
	GetHours() uint64
}

// SkyUtxo skycoin utxo struct
type SkyUtxo struct {
	visor.ReadableOutput
}

// TxOut transaction output filed
type TxOut struct {
	sky.TransactionOutput
}

// GetHash returns utxo hash
func (su SkyUtxo) GetHash() string {
	return su.Hash
}

// GetSrcTx returns source transaction
func (su SkyUtxo) GetSrcTx() string {
	return su.SourceTransaction
}

// GetAddress returns output address
func (su SkyUtxo) GetAddress() string {
	return su.Address
}

// GetCoins returns coins in output
func (su SkyUtxo) GetCoins() uint64 {
	i, err := droplet.FromString(su.Coins)
	if err != nil {
		panic(err)
	}

	return i
}

// GetHours returns coin hours
func (su SkyUtxo) GetHours() uint64 {
	return su.Hours
}

// MakeUtxoOutput generates transaction output base on the addr, amount and hours.
func MakeUtxoOutput(addr string, amount uint64, hours uint64) TxOut {
	uo := TxOut{}
	uo.Address = cipher.MustDecodeBase58Address(addr)
	uo.Coins = amount
	uo.Hours = hours
	return uo
}

// VerifyAmount check if the amout is validated.
func VerifyAmount(amt uint64) error {
	if (amt % 1e6) != 0 {
		return errors.New("Transaction amount must be multiple of 1e6 ")
	}
	return nil
}

// GenerateAddresses generate addresses.
func GenerateAddresses(seed []byte, num int) (string, []coin.AddressEntry) {
	sd, seckeys := cipher.GenerateDeterministicKeyPairsSeed(seed, num)
	entries := make([]coin.AddressEntry, num)
	for i, sec := range seckeys {
		pub := cipher.PubKeyFromSecKey(sec)
		entries[i].Address = cipher.AddressFromPubKey(pub).String()
		entries[i].Public = pub.Hex()
		if !HideSeckey {
			entries[i].Secret = sec.Hex()
		}
	}
	return fmt.Sprintf("%2x", sd), entries
}

// GetUnspentOutputs return the unspent outputs
func GetUnspentOutputs(nodeAddr string, addrs []string) ([]Utxo, error) {
	var url string
	if len(addrs) == 0 {
		return []Utxo{}, nil
	}

	addrParam := strings.Join(addrs, ",")
	url = fmt.Sprintf("http://%s/outputs?addrs=%s", nodeAddr, addrParam)

	rsp, err := http.Get(url)
	if err != nil {
		return []Utxo{}, errors.New("get outputs failed")
	}
	defer rsp.Body.Close()
	outputSet := visor.ReadableOutputSet{}
	if err := json.NewDecoder(rsp.Body).Decode(&outputSet); err != nil {
		return []Utxo{}, err
	}

	spendableOuts := outputSet.SpendableOutputs()
	ux := make([]Utxo, len(spendableOuts))
	for i, u := range spendableOuts {
		ux[i] = SkyUtxo{u}
	}
	return ux, nil
}

func getUnspentOutputsByHashes(nodeAddr string, hashes []string) ([]Utxo, error) {
	if len(hashes) == 0 {
		return []Utxo{}, nil
	}

	url := fmt.Sprintf("http://%s/outputs?hashes=%s", nodeAddr, strings.Join(hashes, ","))
	rsp, err := http.Get(url)
	if err != nil {
		return []Utxo{}, err
	}
	defer rsp.Body.Close()
	outSet := visor.ReadableOutputSet{}
	if err := json.NewDecoder(rsp.Body).Decode(&outSet); err != nil {
		return []Utxo{}, err
	}

	ux := make([]Utxo, len(outSet.HeadOutputs))
	for i, u := range outSet.HeadOutputs {
		ux[i] = SkyUtxo{u}
	}
	return ux, nil
}

// Symbol returns skycoin sybmol
func (sky *Skycoin) Symbol() string {
	return "SKY"
}

// Type returns skycoin type name
func (sky *Skycoin) Type() string {
	return Type
}
