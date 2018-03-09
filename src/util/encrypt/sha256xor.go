package encrypt

import (
	"github.com/skycoin/skycoin/src/cipher/encrypt"
)

var glosha encrypt.Sha256Xor

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
