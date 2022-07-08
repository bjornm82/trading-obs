package broker

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	ENDPOINT_OPEN_TRADES = "/accounts/%s/openTrades"
)

type OpenTradesResponse struct {
	LastTransactionID string `json:"lastTransactionID"`
	Trades            []struct {
		CurrentUnits     string    `json:"currentUnits"`
		Financing        string    `json:"financing"`
		ID               string    `json:"id"`
		InitialUnits     string    `json:"initialUnits"`
		Instrument       string    `json:"instrument"`
		OpenTime         time.Time `json:"openTime"`
		Price            string    `json:"price"`
		RealizedPL       string    `json:"realizedPL"`
		State            string    `json:"state"`
		UnrealizedPL     string    `json:"unrealizedPL"`
		ClientExtensions struct {
			ID string `json:"id"`
		} `json:"clientExtensions,omitempty"`
	} `json:"trades"`
}

func (r *Repo) GetOpenTrades() (OpenTradesResponse, error) {
	o := OpenTradesResponse{}
	a, err := r.cl.GetAccountID()
	if err != nil {
		return o, err
	}
	b, err := r.cl.Read(fmt.Sprintf(ENDPOINT_OPEN_TRADES, a))

	if err != nil {
		return o, err
	}

	err = json.Unmarshal(b, &o)
	if err != nil {
		return o, err
	}

	return o, nil
}
