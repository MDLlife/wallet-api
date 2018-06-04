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
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/util/droplet"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"
	"github.com/MDLlife/wallet-api/src/coin"
	"github.com/MDLlife/wallet-api/src/coin/skycoin"
	walletex "github.com/MDLlife/wallet-api/src/wallet"
)

var (
	maxDropletDivisor int64
	// ErrTxnNoFee is returned if a transaction has no coinhour fee
	ErrTxnNoFee = errors.New("Transaction has zero coinhour fee")

	// ErrNoChangeAddr should not be happened
	ErrNoChangeAddr = errors.New("No change address to receive balance")
)

const (
	// BurnFactor inverse fraction of coinhours that must be burned
	BurnFactor int64 = 2
)

func init() {
	// Compute maxDropletDivisor from precision
	maxDropletDivisor = calculateDivisor(int64(visor.MaxDropletPrecision))
}

// code copy from skycoin visor/visor.go
func calculateDivisor(precision int64) int64 {
	if precision > droplet.Exponent {
		log.Panic("precision must be <= droplet.Exponent")
	}

	n := droplet.Exponent - precision
	var i int64 = 1
	for k := int64(0); k < n; k++ {
		i = i * 10
	}
	return i
}

// DropletPrecisionCheck checks if an amount of coins is valid given decimal place restrictions
func DropletPrecisionCheck(amount int64) error {
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
	//tx := visor.TransactionResult{}
	tx := daemon.TransactionResult{}
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

func getChangeAddr(uxBal []wallet.UxBalance) (string, error) {
	// this will not be happened
	if len(uxBal) == 0 {
		return "", errors.New("no change address")
	}
	chgAddr := uxBal[0].Address.String()
	_, err := cipher.DecodeBase58Address(chgAddr)
	return chgAddr, err
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

	uxBalances, err := chooseSpends(totalUtxos, p.Amount)
	if err != nil {
		return nil, nil, err
	}

	bal, hours := func(uxBal []wallet.UxBalance) (uint64, uint64) {
		var c, h uint64
		for _, u := range uxBal {
			c += u.Coins
			h += u.Hours
		}
		return c, h
	}(uxBalances)

	if hours == 0 {
		return nil, nil, ErrTxnNoFee
	}

	txIns := make([]coin.TxIn, len(uxBalances))
	for i, u := range uxBalances {
		txIns[i] = coin.TxIn{
			Txid:    u.Hash.Hex(),
			Address: u.Address.String(),
		}
	}

	var txOut []skycoin.TxOut
	chgAmt := bal - p.Amount
	haveChange := chgAmt > 0
	chgHours, addrHours := distributeSpendHours(hours, haveChange)

	if chgAmt > 0 {
		chgAddr, err := getChangeAddr(uxBalances)
		if err != nil {
			return nil, nil, ErrNoChangeAddr
		}
		txOut = append(txOut,
			cn.makeTxOut(p.ToAddr, p.Amount, addrHours),
			cn.makeTxOut(chgAddr, chgAmt, chgHours))
	} else {
		txOut = append(txOut, cn.makeTxOut(p.ToAddr, p.Amount, addrHours))
	}
	return txIns, txOut, nil
}

// requiredFee returns the coinhours fee required for an amount of hours
// The required fee is calculated as hours/BurnFactor, rounded up.
func requiredFee(hours int64) int64 {
	feeHours := hours / BurnFactor
	if hours%BurnFactor != 0 {
		feeHours++
	}

	return feeHours
}

// distributeSpendHours calculates how many coin hours to transfer to the change address and how
// many to transfer to the destination addresses.
func distributeSpendHours(inputHours uint64, haveChange bool) (uint64, uint64) {
	feeHours := requiredFee(int64(inputHours))
	remainingHours := inputHours - uint64(feeHours)

	var changeHours uint64
	if haveChange {
		// Split the remaining hours between the change output and the other outputs
		changeHours = remainingHours / 2

		// If remainingHours is an odd number, give the extra hour to the change output
		if remainingHours%2 == 1 {
			changeHours++
		}
	}

	// Distribute the remaining hours
	addrHours := remainingHours - changeHours

	return changeHours, addrHours
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

	if err := DropletPrecisionCheck(int64(amt)); err != nil {
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

func chooseSpends(uxouts visor.ReadableOutputSet, coins uint64) ([]wallet.UxBalance, error) {
	// Convert spendable unspent outputs to []wallet.UxBalance
	spendableOutputs, err := visor.ReadableOutputsToUxBalances(uxouts.SpendableOutputs())
	if err != nil {
		return nil, err
	}

	// Choose which unspent outputs to spend
	// Use the MinimizeUxOuts strategy, since this is most likely used by
	// application that may need to send frequently.
	// Using fewer UxOuts will leave more available for other transactions,
	// instead of waiting for confirmation.
	outs, err := wallet.ChooseSpendsMinimizeUxOuts(spendableOutputs, coins, 0)
	if err != nil {
		// If there is not enough balance in the spendable outputs,
		// see if there is enough balance when including incoming outputs
		if err == wallet.ErrInsufficientBalance {
			expectedOutputs, otherErr := visor.ReadableOutputsToUxBalances(uxouts.ExpectedOutputs())
			if otherErr != nil {
				return nil, otherErr
			}

			if _, otherErr := wallet.ChooseSpendsMinimizeUxOuts(expectedOutputs, coins, 0); otherErr != nil {
				return nil, err
			}

			return nil, errors.New("balance is not sufficient")
		}

		return nil, err
	}

	return outs, nil
}
