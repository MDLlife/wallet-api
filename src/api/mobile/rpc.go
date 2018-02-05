package mobile

import (
	"errors"
	"fmt"

	"github.com/skycoin/skycoin/src/api/cli"
	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/spaco/spo/src/visor"
)

// RPCError wraps errors from the skycoin CLI/RPC library
type RPCError struct {
	error
}

// NewRPCError wraps an err with RPCError
func NewRPCError(err error) RPCError {
	return RPCError{err}
}

// RPC provides methods for sending coins
type WalletRPC struct {
	WalletControl *MyService
	rpcClient     *webrpc.Client
}

// NewRPC creates RPC instance
func NewWalletRPC(walletServ *MyService, rpcAddr string) (*WalletRPC, error) {
	rpcClient := &webrpc.Client{
		Addr: rpcAddr,
	}

	return &WalletRPC{
		rpcClient:     rpcClient,
		WalletControl: walletServ,
	}, nil
}

// CreateTransaction creates a raw Skycoin transaction offline, that can be broadcast later
func (c *WalletRPC) CreateTransaction(walletFile string, sendAmountList []cli.SendAmount) (*coin.Transaction, error) {
	if err := validateSendAmount(sendAmountList); err != nil {
		return nil, err
	}

	wlt, err := c.WalletControl.GetWallet(walletFile)
	if err != nil {
		return nil, err
	}
	changeAddr := wlt.Entries[0].Address.String()
	txn, err := cli.CreateRawTxFromWallet(c.rpcClient, walletFile, changeAddr, sendAmountList)
	if err != nil {
		return nil, RPCError{err}
	}

	return txn, nil
}

// BroadcastTransaction broadcasts a transaction and returns its txid
func (c *WalletRPC) BroadcastTransaction(tx *coin.Transaction) (string, error) {
	txid, err := c.rpcClient.InjectTransaction(tx)
	if err != nil {
		return "", RPCError{err}
	}

	return txid, nil
}

// GetTransaction returns transaction by txid
func (c *WalletRPC) GetTransaction(txid string) (*webrpc.TxnResult, error) {
	txn, err := c.rpcClient.GetTransactionByID(txid)
	if err != nil {
		return nil, RPCError{err}
	}

	return txn, nil
}

// Balance returns the balance of a wallet
func (c *WalletRPC) Balance(walletFile string) (*cli.Balance, error) {
	bal, err := cli.CheckWalletBalance(c.rpcClient, walletFile)
	if err != nil {
		return nil, RPCError{err}
	}

	return &bal.Spendable, nil
}

func validateSendAmount(amtList []cli.SendAmount) error {
	// validate the recvAddr
	for _, amt := range amtList {
		if _, err := cipher.DecodeBase58Address(amt.Addr); err != nil {
			return err
		}

		if amt.Coins == 0 {
			return errors.New("Skycoin send amount is 0")
		}
	}

	return nil
}

// Send sends coins to recv address
func (c *WalletRPC) Send(walletFile, recvAddr string, amount int64) (string, error) {
	// validate the recvAddr
	if _, err := cipher.DecodeBase58Address(recvAddr); err != nil {
		return "", err
	}

	if amount == 0 {
		return "", errors.New("Can't send 0 coins")
	}

	sendAmount := cli.SendAmount{
		Addr:  recvAddr,
		Coins: uint64(amount),
	}

	wlt, err := c.WalletControl.GetWallet(walletFile)
	if err != nil {
		return "", err
	}
	changeAddr := wlt.Entries[0].Address.String()
	return cli.SendFromWallet(c.rpcClient, walletFile, changeAddr, []cli.SendAmount{sendAmount})
}

func (c *WalletRPC) GetBlocks(start, end uint64) (*visor.ReadableBlocks, error) {
	param := []uint64{start, end}
	blocks := visor.ReadableBlocks{}

	if err := c.rpcClient.Do(&blocks, "get_blocks", param); err != nil {
		return nil, err
	}

	return &blocks, nil
}
func (c *WalletRPC) GetBlockCount() (int, error) {
	st, err := c.rpcClient.GetStatus()
	if err != nil {
		return 0, err
	}
	return int(st.BlockNum), nil
}

func (c *WalletRPC) GetBlocksBySeq(seq uint64) (*visor.ReadableBlock, error) {
	ss := []uint64{seq}
	blocks := visor.ReadableBlocks{}

	if err := c.rpcClient.Do(&blocks, "get_blocks_by_seq", ss); err != nil {
		return nil, err
	}

	if len(blocks.Blocks) == 0 {
		return nil, nil
	}

	return &blocks.Blocks[0], nil
}

func (c *WalletRPC) GetLastBlocks() (*visor.ReadableBlock, error) {
	param := []uint64{1}
	blocks := visor.ReadableBlocks{}
	if err := c.rpcClient.Do(&blocks, "get_lastblocks", param); err != nil {
		return nil, err
	}

	if len(blocks.Blocks) == 0 {
		return nil, nil
	}
	return &blocks.Blocks[0], nil
}

// Send sends coins to batch recv address
func (c *WalletRPC) SendBatch(walletFile string, saList []cli.SendAmount) (string, error) {
	// validate the recvAddr
	wlt, err := c.WalletControl.GetWallet(walletFile)
	if err != nil {
		return "", err
	}
	changeAddr := wlt.Entries[0].Address.String()
	for _, sendAmount := range saList {
		if _, err := cipher.DecodeBase58Address(sendAmount.Addr); err != nil {
			return "", err
		}
		if sendAmount.Coins == 0 {
			return "", fmt.Errorf("Can't send 0 coins", sendAmount.Coins)
		}

	}

	return cli.SendFromWallet(c.rpcClient, walletFile, changeAddr, saList)
}
