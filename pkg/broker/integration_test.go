package broker

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	ACCOUNT_ID = "101-004-8979399-001"
	TOKEN      = "a7f4f34b8f41a1d1d467047517c9b8a0-73ef8ea2c18a28fbb89c0fc21786b9b9"
)

var cl = NewClient(
	ACCOUNT_ID,
	TOKEN,
	false,
)

var broker = New(cl)

func TestGetAccountID(t *testing.T) {
	v, err := cl.GetAccountID()
	assert.NoError(t, err)
	assert.Equal(t, ACCOUNT_ID, v)
}

func TestGetPreciseCandles(t *testing.T) {
	const layout = "2006-01-02T15:04:05"
	// TODO: 4 o'clock is 6 o'clock on the API and 9 at home? Strange, needs some attention
	then, err := time.Parse(layout, "2022-07-06T04:00:00")
	if err != nil {
		t.Error(err)
		return
	}
	from := then.Add(-time.Hour * 2)
	to := then.Add(-time.Hour * 1)

	v, err := broker.GetPreciseCandles("DE30_EUR", from, to, "M5")

	log.Println(v.Candles[0].Time)
	assert.NoError(t, err)
	assert.Equal(t, ACCOUNT_ID, v)
}

func TestGetInstruments(t *testing.T) {
	i, err := broker.GetInstruments()
	assert.NoError(t, err)
	assert.Greater(t, len(i.Instruments), 1)

	var exists = false
	var displayName = ""
	for _, v := range i.Instruments {
		if v.Name == "UK100_GBP" {
			exists = true
			displayName = v.DisplayName
		}
	}

	assert.True(t, exists)
	assert.Equal(t, "UK 100", displayName)
}

func TestGetOrders(t *testing.T) {
	if t.Name() != "integration" {
		t.Skip()
	}

	b, err := broker.GetOrders()

	assert.NoError(t, err)
	assert.NotEqual(t, "0", b.LastTransactionID)
	assert.NotEqual(t, "", b.LastTransactionID)
	assert.Equal(t, "438", b.LastTransactionID)
}

func TestCreateOrder(t *testing.T) {
	if t.Name() != "integration" {
		t.Skip()
	}

	b, err := broker.CreateOrder(1, "UK100_GBP", "MARKET")

	assert.NoError(t, err)
	assert.NotEqual(t, "0", b.LastTransactionID)
	assert.NotEqual(t, "", b.LastTransactionID)
	assert.Equal(t, "438", b.LastTransactionID)
}

func TestGetPositionsByInstrument(t *testing.T) {
	if t.Name() != "integration" {
		t.Skip()
	}

	b, err := broker.GetPosition("UK100_GBP")

	assert.NoError(t, err)
	assert.Equal(t, "0.0", b.Position.Long.Units)
	assert.Equal(t, "0", b.Position.Short.Units)
}

func TestGetOpenPositions(t *testing.T) {
	if t.Name() != "integration" {
		t.Skip()
	}

	b, err := broker.GetPositions()
	assert.NoError(t, err)
	assert.Len(t, b.Positions, 0)
}
