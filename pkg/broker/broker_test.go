package broker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testStr = []byte("{\"value\": \"hello\"}")

func TestNew(t *testing.T) {

	m := Mock{}
	m.ReturnByte = testStr
	m.ReturnError = nil
	broker := New(&m)
	o, err := broker.GetOrders()

	assert.Equal(t, "", o.LastTransactionID)
	assert.NoError(t, err)
}
