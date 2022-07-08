package oanda

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPricing(t *testing.T) {
	ctx := Context{
		"https://api-fxpractice.oanda.com",
		"https://stream-fxpractice.oanda.com/",
		"a7f4f34b8f41a1d1d467047517c9b8a0-73ef8ea2c18a28fbb89c0fc21786b9b9",
		"101-004-8979399-001",
		"OBS DAX",
	}

	api := ctx.CreateAPI()

	pr, err := api.GetCandles("DE30_EUR", 5000, "M5")
	assert.NoError(t, err)
	assert.Equal(t, "M5", pr.Granularity)
	assert.Equal(t, "DE30_EUR", pr.Instrument)
	assert.Equal(t, 5000, len(pr.Candles))
}

type Instruments struct {
	Instruments []Instrument `json:"instruments"`
}

type Instrument struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}
