package encrypt

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWallet(t *testing.T) {
	content := "a xx yy 1243 !!@ bn dx"
	password := "123456781qaz2wsx"
	cipher, err := Encrypt([]byte(password), content)
	assert.NoError(t, err)
	assert.NotEmpty(t, cipher)
	fmt.Printf("cipher: %s\n", cipher)

	origin, err := Decrypt([]byte(password), cipher)
	assert.NoError(t, err)
	assert.Equal(t, origin, content)
}
