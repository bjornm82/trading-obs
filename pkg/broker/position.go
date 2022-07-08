package broker

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	ENDPOINT_POSITION  = "/accounts/%s/positions/%s"
	ENDPOINT_POSITIONS = "/accounts/%s/openPositions"
)

func (r *Repo) GetPosition(instrument string) (PositionResponse, error) {
	resp := PositionResponse{}
	if instrument == "" {
		return resp, errors.New("position instrument not given")
	}

	i, err := r.GetInstruments()
	if err != nil {
		return resp, fmt.Errorf("unable to retreive list of instruments: %s", err)
	}

	var exists = false
	for _, v := range i.Instruments {
		if v.Name == instrument {
			exists = true
		}
	}
	if !exists {
		return resp, fmt.Errorf("unable to find instrument with name: %s", instrument)
	}

	a, err := r.cl.GetAccountID()
	if err != nil {
		return resp, err
	}
	b, err := r.cl.Read(
		fmt.Sprintf(ENDPOINT_POSITION, a, instrument))
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(b, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (r *Repo) GetPositions() (PositionsResponse, error) {
	resp := PositionsResponse{}
	a, err := r.cl.GetAccountID()
	if err != nil {
		return resp, err
	}
	b, err := r.cl.Read(
		fmt.Sprintf(ENDPOINT_POSITIONS, a))
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(b, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

type PositionsResponse struct {
	LastTransactionID string `json:"lastTransactionID"`
	Positions         []struct {
		Instrument string `json:"instrument"`
		Long       struct {
			AveragePrice string   `json:"averagePrice"`
			Pl           string   `json:"pl"`
			ResettablePL string   `json:"resettablePL"`
			TradeIDs     []string `json:"tradeIDs"`
			Units        string   `json:"units"`
			UnrealizedPL string   `json:"unrealizedPL"`
		} `json:"long"`
		Pl           string `json:"pl"`
		ResettablePL string `json:"resettablePL"`
		Short        struct {
			AveragePrice string   `json:"averagePrice"`
			Pl           string   `json:"pl"`
			ResettablePL string   `json:"resettablePL"`
			TradeIDs     []string `json:"tradeIDs"`
			Units        string   `json:"units"`
			UnrealizedPL string   `json:"unrealizedPL"`
		} `json:"short"`
		UnrealizedPL string `json:"unrealizedPL"`
	} `json:"positions"`
}

type PositionResponse struct {
	Position struct {
		Instrument string `json:"instrument"`
		Long       struct {
			Units                   string `json:"units"`
			Pl                      string `json:"pl"`
			ResettablePL            string `json:"resettablePL"`
			Financing               string `json:"financing"`
			DividendAdjustment      string `json:"dividendAdjustment"`
			GuaranteedExecutionFees string `json:"guaranteedExecutionFees"`
			UnrealizedPL            string `json:"unrealizedPL"`
		} `json:"long"`
		Short struct {
			Units                   string `json:"units"`
			Pl                      string `json:"pl"`
			ResettablePL            string `json:"resettablePL"`
			Financing               string `json:"financing"`
			DividendAdjustment      string `json:"dividendAdjustment"`
			GuaranteedExecutionFees string `json:"guaranteedExecutionFees"`
			UnrealizedPL            string `json:"unrealizedPL"`
		} `json:"short"`
		Pl                      string `json:"pl"`
		ResettablePL            string `json:"resettablePL"`
		Financing               string `json:"financing"`
		Commission              string `json:"commission"`
		DividendAdjustment      string `json:"dividendAdjustment"`
		GuaranteedExecutionFees string `json:"guaranteedExecutionFees"`
		UnrealizedPL            string `json:"unrealizedPL"`
	} `json:"position"`
	LastTransactionID string `json:"lastTransactionID"`
}
