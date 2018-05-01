package mobile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoinHour(t *testing.T) {
	should := requiredFee(10)
	assert.Equal(t, int64(5), should)
	should = requiredFee(11)
	assert.Equal(t, int64(6), should)

	balance := 10
	coins := 9
	chgAmt := balance - coins
	haveChange := chgAmt > 0

	chgHours, addrHours := distributeSpendHours(10, haveChange)
	assert.Equal(t, uint64(3), chgHours)
	assert.Equal(t, uint64(2), addrHours)

	chgHours, addrHours = distributeSpendHours(11, haveChange)
	assert.Equal(t, uint64(3), chgHours)
	assert.Equal(t, uint64(2), addrHours)

	balance = 10
	coins = 10
	chgAmt = balance - coins
	haveChange = chgAmt > 0
	chgHours, addrHours = distributeSpendHours(10, haveChange)
	assert.Equal(t, uint64(0), chgHours)
	assert.Equal(t, uint64(5), addrHours)
	chgHours, addrHours = distributeSpendHours(11, haveChange)
	assert.Equal(t, uint64(0), chgHours)
	assert.Equal(t, uint64(5), addrHours)
}
