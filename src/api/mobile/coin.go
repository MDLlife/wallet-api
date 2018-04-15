package mobile

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"

	"strings"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util/droplet"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/spolabs/wallet-api/src/coin"
	"github.com/spolabs/wallet-api/src/coin/skycoin"
	walletex "github.com/spolabs/wallet-api/src/wallet"
)

var maxDropletDivisor uint64

func init() {
	// Compute maxDropletDivisor from precision
	maxDropletDivisor = calculateDivisor(visor.MaxDropletPrecision)
}

// code copy from skycoin visor/visor.go
func calculateDivisor(precision uint64) uint64 {
	if precision > droplet.Exponent {
		log.Panic("precision must be <= droplet.Exponent")
	}

	n := droplet.Exponent - precision
	var i uint64 = 1
	for k := uint64(0); k < n; k++ {
		i = i * 10
	}
	return i
}

// DropletPrecisionCheck checks if an amount of coins is valid given decimal place restrictions
func DropletPrecisionCheck(amount uint64) error {
	if amount%maxDropletDivisor != 0 {
		return errors.New("invalid amount, too many decimal places")
	}
	return nil
}

// Coiner coin client interface
type Coiner interface {
	Name() string
	GetBalance(addrs string) (string, error)
	ValidateAddr(addr string) error
	CreateRawTx(txIns []coin.TxIn, getKey coin.GetPrivKey, txOuts interface{}) (string, error)
	BroadcastTx(rawtx string) (string, error)
	GetTransactionByID(txid string) (string, error)
	GetNodeAddr() string
	IsTransactionConfirmed(txid string) (bool, error)
	Send(walletID, toAddr, amount, passwd string) (string, error)
}

// CoinEx implements the Coin interface.
type coinEx struct {
	name     string
	nodeAddr string
}

type sendParams struct {
	WalletID string
	ToAddr   string
	Amount   uint64
}

func newCoin(name, nodeAddr string) *coinEx {
	return &coinEx{name: name, nodeAddr: nodeAddr}
}

func (cn coinEx) Name() string {
	return cn.name
}

// GetNodeAddr returns the coin's node address
func (cn coinEx) GetNodeAddr() string {
	return cn.nodeAddr
}

func (cn coinEx) getOutputs(addrs []string) (string, error) {
	/*	curl http://127.0.0.1:8620/outputs?addrs=26pnkiNjHVmvG4vj7FcQuxGob2fZXWFiZqT
		{
		    "head_outputs": [
		        {
		            "hash": "37077d5e5662840b3fe78d1716e864e7554a03468a3844a5fe002e6012130099",
		            "block_seq": 1159,
		            "src_tx": "47f9f808ca60c6cc94686c855e7f76e09a90bf7c9b4dca297fed41f7eccdd7b8",
		            "address": "26pnkiNjHVmvG4vj7FcQuxGob2fZXWFiZqT",
		            "coins": "10.000000",
		            "hours": 4919
		        }
		    ],
		    "outgoing_outputs": [],
		    "incoming_outputs": []
		}
	*/
	addrsArgs := strings.Join(addrs, ",")
	url := fmt.Sprintf("http://%s/outputs?addrs=%s", cn.nodeAddr, addrsArgs)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	allBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(allBody), nil
}

// GetBlocks gets balance of specific addresses
func (cn coinEx) GetBlocks(start, end int) (string, error) {
	url := fmt.Sprintf("http://%s/blocks?start=%d&end=%d", cn.nodeAddr, start, end)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	allBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(allBody), nil
}

func (cn coinEx) GetTransactionByID(txid string) (string, error) {
	url := fmt.Sprintf("http://%s/transaction?txid=%s", cn.nodeAddr, txid)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	allBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(allBody), nil
}

func (cn coinEx) IsTransactionConfirmed(txid string) (bool, error) {
	txstr, err := cn.GetTransactionByID(txid)
	if err != nil {
		return false, err
	}
	tx := visor.TransactionResult{}
	err = json.Unmarshal([]byte(txstr), &tx)
	if err != nil {
		return false, err
	}
	return tx.Status.Confirmed, nil
}

// GetBalance args is address joined by "," such as "a1,a2,a3"
func (cn coinEx) GetBalance(addrs string) (string, error) {
	url := fmt.Sprintf("http://%s/balance?addrs=%s", cn.nodeAddr, addrs)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	allBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(allBody), nil

}

// ValidateAddr check if the address is validated
func (cn coinEx) ValidateAddr(address string) error {
	_, err := cipher.DecodeBase58Address(address)
	return err
}

// CreateRawTx creates raw transaction
func (cn coinEx) CreateRawTx(txIns []coin.TxIn, getKey coin.GetPrivKey, txOuts interface{}) (string, error) {
	tx := skycoin.Transaction{}
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
		out := o.(skycoin.TxOut)
		// todo decimal can be set
		if (out.Coins % 1e3) != 0 {
			return "", fmt.Errorf("%s coins must be multiple of 1e3", cn.Name())
		}
		tx.PushOutput(out.Address, out.Coins, out.Hours)
	}

	keys := make([]cipher.SecKey, len(txIns))
	for i, in := range txIns {
		s, err := getKey(in.Address)
		if err != nil {
			return "", fmt.Errorf("get private key failed:%v", err)
		}
		k, err := cipher.SecKeyFromHex(s)
		if err != nil {
			return "", fmt.Errorf("invalid private key:%v", err)
		}
		keys[i] = k
	}

	tx.SignInputs(keys)
	tx.UpdateHeader()

	d, err := tx.Serialize()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(d), nil
}

// BroadcastTx injects transaction
func (cn coinEx) BroadcastTx(rawtx string) (string, error) {
	url := fmt.Sprintf("http://%s/injectTransaction", cn.nodeAddr)
	client := &http.Client{}
	v := struct {
		Rawtx string `json:"rawtx"`
	}{}
	v.Rawtx = rawtx
	reqBody, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(reqBody))
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	allBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(allBody), nil

}

// PrepareTx prepares the transaction info
func (cn coinEx) PrepareTx(params interface{}) ([]coin.TxIn, interface{}, error) {
	p := params.(sendParams)

	tp := strings.Split(p.WalletID, "_")[0]
	if tp != cn.name {
		return nil, nil, fmt.Errorf("invalid wallet %v", tp)
	}

	// validate address
	if err := cn.ValidateAddr(p.ToAddr); err != nil {
		return nil, nil, err
	}

	addrs, err := walletex.GetAddresses(p.WalletID)
	if err != nil {
		return nil, nil, err
	}

	totalUtxosStr, err := cn.getOutputs(addrs)
	if err != nil {
		return nil, nil, err
	}
	totalUtxos := visor.ReadableOutputSet{}
	err = json.Unmarshal([]byte(totalUtxosStr), &totalUtxos)
	if err != nil {
		return nil, nil, err
	}

	utxos, err := cn.getSufficientOutputs(totalUtxos, p.Amount)
	if err != nil {
		return nil, nil, err
	}

	bal, hours := func(utxos visor.ReadableOutputs) (uint64, uint64) {
		var c, h uint64
		for _, u := range utxos {
			coins, err := droplet.FromString(u.Coins)
			if err != nil {
				continue
			}
			c += coins
			h += u.Hours
		}
		return c, h
	}(utxos)

	txIns := make([]coin.TxIn, len(utxos))
	for i, u := range utxos {
		txIns[i] = coin.TxIn{
			Txid:    u.Hash,
			Address: u.Address,
		}
	}

	var txOut []skycoin.TxOut
	chgAmt := bal - p.Amount
	chgHours := hours / 4
	chgAddr := addrs[0]
	if chgAmt > 0 {
		txOut = append(txOut,
			cn.makeTxOut(p.ToAddr, p.Amount, chgHours/2),
			cn.makeTxOut(chgAddr, chgAmt, chgHours/2))
	} else {
		txOut = append(txOut, cn.makeTxOut(p.ToAddr, p.Amount, chgHours/2))
	}
	return txIns, txOut, nil
}

// Send sends numbers of coins to toAddr from specific wallet
func (cn *coinEx) Send(walletID, toAddr, amount, passwd string) (string, error) {
	// validate amount
	amt, err := droplet.FromString(amount)
	if err != nil {
		return "", fmt.Errorf("amount to droplets failed: %v", err)
	}
	if amt == 0 {
		return "", fmt.Errorf("can not send 0 coins")
	}
	if _, err := cipher.DecodeBase58Address(toAddr); err != nil {
		return "", err
	}

	if err := DropletPrecisionCheck(amt); err != nil {
		return "", err
	}

	params := sendParams{WalletID: walletID, ToAddr: toAddr, Amount: amt}

	txIns, txOut, err := cn.PrepareTx(params)
	if err != nil {
		return "", err
	}

	// prepare keys
	rawtx, err := cn.CreateRawTx(txIns, getPrivateKey(walletID, passwd), txOut)
	if err != nil {
		return "", fmt.Errorf("create raw transaction failed:%v", err)
	}

	txid, err := cn.BroadcastTx(rawtx)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`{"txid":%s}`, txid), nil
}

func (cn coinEx) makeTxOut(addr string, coins uint64, hours uint64) skycoin.TxOut {
	out := skycoin.TxOut{}
	out.Address = cipher.MustDecodeBase58Address(addr)
	out.Coins = coins
	out.Hours = hours
	return out
}

func (cn coinEx) getSufficientOutputs(utxos visor.ReadableOutputSet, amt uint64) (visor.ReadableOutputs, error) {
	outMap := make(map[string]visor.ReadableOutputs)
	for _, u := range utxos.SpendableOutputs() {
		outMap[u.Address] = append(outMap[u.Address], u)
	}

	allUtxos := visor.ReadableOutputs{}
	var allBal uint64
	for _, utxos := range outMap {
		allBal += func(utxos visor.ReadableOutputs) uint64 {
			var bal uint64
			for _, u := range utxos {
				coins, err := droplet.FromString(u.Coins)
				if err != nil {
					continue
				}
				if coins == 0 {
					continue
				}
				bal += coins
			}
			return bal
		}(utxos)

		allUtxos = append(allUtxos, utxos...)
		if allBal >= amt {
			return allUtxos, nil
		}
	}

	return nil, errors.New("insufficient balance")
}
