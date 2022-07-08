package account

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var bank = New()

func TestBankGetStartCaptital(t *testing.T) {
	assert.Equal(t, float64(START_CAPITAL), bank.GetCapital())
}

func TestBankUpdateCapital(t *testing.T) {
	var value = 99.90

	bank.UpdateCapital(float64(value))
	assert.Equal(t, float64(START_CAPITAL+value), bank.GetCapital())
}

func TestBankUpdateToStartCapital(t *testing.T) {
	var value = -99.90

	bank.UpdateCapital(float64(value))
	assert.Equal(t, float64(START_CAPITAL), bank.GetCapital())
}
