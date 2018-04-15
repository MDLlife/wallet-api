package coin

// TxIn records the tx vin info, txid is the prevous txid, Index is the out index in previous tx.
type TxIn struct {
	Txid    string
	Address string
	Vout    uint32
}

// GetPrivKey is a callback func used for SignTx func to get relevant private key of specific address.
type GetPrivKey func(addr string) (string, error)

// AddressEntry represents the wallet address
type AddressEntry struct {
	Address string `json:"address"`
	Public  string `json:"pubkey"`
	Secret  string `json:"seckey"`
}
