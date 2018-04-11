package encrypt

import (
	"github.com/skycoin/skycoin/src/cipher/encrypt"
)

//var glosha encrypt.Sha256Xor
var glosha encrypt.ScryptChacha20poly1305

func init() {
	glosha = encrypt.DefaultScryptChacha20poly1305
	glosha.N = 1 << 16
}

//Encrypt encrypt text
func Encrypt(key []byte, text string) (string, error) {
	encry, err := glosha.Encrypt([]byte(text), key)
	if err != nil {
		return "", err
	}
	return string(encry), nil
}

// Decrypt decrypt text
func Decrypt(key []byte, text string) (string, error) {
	decry, err := glosha.Decrypt([]byte(text), key)
	if err != nil {
		return "", err
	}
	return string(decry), nil
}
