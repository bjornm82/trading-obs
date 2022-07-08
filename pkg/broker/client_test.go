package broker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateURLShouldFail(t *testing.T) {
	v, err := createUrl(PRACTICE_HOST, "instruments")
	assert.Error(t, err)
	assert.Equal(t, "", v)
}

func TestCreatePracticeURL(t *testing.T) {
	v, err := createUrl(PRACTICE_HOST, "/instruments")
	assert.NoError(t, err)
	assert.Equal(t, "https://api-fxpractice.oanda.com/v3/instruments", v)
}

func TestCreateLiveURL(t *testing.T) {
	v, err := createUrl(LIVE_HOST, "/instruments")
	assert.NoError(t, err)
	assert.Equal(t, "https://api-fxtrade.oanda.com/v3/instruments", v)
}
