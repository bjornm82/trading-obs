package broker

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	ENDPOINT_ACCOUNTS    = "/accounts/%s"
	ENDPOINT_INSTRUMENTS = "/accounts/%s/instruments"
)

type AccountResponse struct {
	Account struct {
		GuaranteedStopLossOrderMode string    `json:"guaranteedStopLossOrderMode"`
		HedgingEnabled              bool      `json:"hedgingEnabled"`
		ID                          string    `json:"id"`
		CreatedTime                 time.Time `json:"createdTime"`
		Currency                    string    `json:"currency"`
		CreatedByUserID             int       `json:"createdByUserID"`
		Alias                       string    `json:"alias"`
		MarginRate                  string    `json:"marginRate"`
		LastTransactionID           string    `json:"lastTransactionID"`
		Balance                     string    `json:"balance"`
		OpenTradeCount              int       `json:"openTradeCount"`
		OpenPositionCount           int       `json:"openPositionCount"`
		PendingOrderCount           int       `json:"pendingOrderCount"`
		Pl                          string    `json:"pl"`
		ResettablePL                string    `json:"resettablePL"`
		ResettablePLTime            string    `json:"resettablePLTime"`
		Financing                   string    `json:"financing"`
		Commission                  string    `json:"commission"`
		DividendAdjustment          string    `json:"dividendAdjustment"`
		GuaranteedExecutionFees     string    `json:"guaranteedExecutionFees"`
		Orders                      []struct {
			ID               string    `json:"id"`
			CreateTime       time.Time `json:"createTime"`
			Type             string    `json:"type"`
			TradeID          string    `json:"tradeID"`
			Price            string    `json:"price"`
			TimeInForce      string    `json:"timeInForce"`
			TriggerCondition string    `json:"triggerCondition"`
			State            string    `json:"state"`
			TriggerMode      string    `json:"triggerMode,omitempty"`
		} `json:"orders"`
		Positions []struct {
			Instrument string `json:"instrument"`
			Short      struct {
				Units                   string `json:"units"`
				Pl                      string `json:"pl"`
				ResettablePL            string `json:"resettablePL"`
				Financing               string `json:"financing"`
				DividendAdjustment      string `json:"dividendAdjustment"`
				GuaranteedExecutionFees string `json:"guaranteedExecutionFees"`
				UnrealizedPL            string `json:"unrealizedPL"`
			} `json:"short"`
			Long struct {
				Units                   string   `json:"units"`
				AveragePrice            string   `json:"averagePrice"`
				Pl                      string   `json:"pl"`
				ResettablePL            string   `json:"resettablePL"`
				Financing               string   `json:"financing"`
				DividendAdjustment      string   `json:"dividendAdjustment"`
				GuaranteedExecutionFees string   `json:"guaranteedExecutionFees"`
				TradeIDs                []string `json:"tradeIDs"`
				UnrealizedPL            string   `json:"unrealizedPL"`
			} `json:"long,omitempty"`
			Pl                      string `json:"pl"`
			ResettablePL            string `json:"resettablePL"`
			Financing               string `json:"financing"`
			Commission              string `json:"commission"`
			DividendAdjustment      string `json:"dividendAdjustment"`
			GuaranteedExecutionFees string `json:"guaranteedExecutionFees"`
			UnrealizedPL            string `json:"unrealizedPL"`
			MarginUsed              string `json:"marginUsed,omitempty"`
		} `json:"positions"`
		Trades []struct {
			ID                    string    `json:"id"`
			Instrument            string    `json:"instrument"`
			Price                 string    `json:"price"`
			OpenTime              time.Time `json:"openTime"`
			InitialUnits          string    `json:"initialUnits"`
			InitialMarginRequired string    `json:"initialMarginRequired"`
			State                 string    `json:"state"`
			CurrentUnits          string    `json:"currentUnits"`
			RealizedPL            string    `json:"realizedPL"`
			Financing             string    `json:"financing"`
			DividendAdjustment    string    `json:"dividendAdjustment"`
			TakeProfitOrderID     string    `json:"takeProfitOrderID"`
			StopLossOrderID       string    `json:"stopLossOrderID"`
			UnrealizedPL          string    `json:"unrealizedPL"`
			MarginUsed            string    `json:"marginUsed"`
		} `json:"trades"`
		UnrealizedPL                string `json:"unrealizedPL"`
		Nav                         string `json:"NAV"`
		MarginUsed                  string `json:"marginUsed"`
		MarginAvailable             string `json:"marginAvailable"`
		PositionValue               string `json:"positionValue"`
		MarginCloseoutUnrealizedPL  string `json:"marginCloseoutUnrealizedPL"`
		MarginCloseoutNAV           string `json:"marginCloseoutNAV"`
		MarginCloseoutMarginUsed    string `json:"marginCloseoutMarginUsed"`
		MarginCloseoutPositionValue string `json:"marginCloseoutPositionValue"`
		MarginCloseoutPercent       string `json:"marginCloseoutPercent"`
		WithdrawalLimit             string `json:"withdrawalLimit"`
		MarginCallMarginUsed        string `json:"marginCallMarginUsed"`
		MarginCallPercent           string `json:"marginCallPercent"`
	} `json:"account"`
	LastTransactionID string `json:"lastTransactionID"`
}

func (r *Repo) InstrumentExists(value string) (bool, error) {
	return true, nil
}

func (r *Repo) GetAccount() (AccountResponse, error) {
	resp := AccountResponse{}

	a, err := r.cl.GetAccountID()
	if err != nil {
		return resp, err
	}
	b, err := r.cl.Read(
		fmt.Sprintf(ENDPOINT_ACCOUNTS, a))
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(b, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (r *Repo) GetInstruments() (InstrumentResponse, error) {
	resp := InstrumentResponse{}

	a, err := r.cl.GetAccountID()
	if err != nil {
		return resp, err
	}
	b, err := r.cl.Read(
		fmt.Sprintf(ENDPOINT_INSTRUMENTS, a))
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(b, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

type InstrumentResponse struct {
	Instruments []struct {
		Name                        string `json:"name"`
		Type                        string `json:"type"`
		DisplayName                 string `json:"displayName"`
		PipLocation                 int    `json:"pipLocation"`
		DisplayPrecision            int    `json:"displayPrecision"`
		TradeUnitsPrecision         int    `json:"tradeUnitsPrecision"`
		MinimumTradeSize            string `json:"minimumTradeSize"`
		MaximumTrailingStopDistance string `json:"maximumTrailingStopDistance"`
		MinimumTrailingStopDistance string `json:"minimumTrailingStopDistance"`
		MaximumPositionSize         string `json:"maximumPositionSize"`
		MaximumOrderUnits           string `json:"maximumOrderUnits"`
		MarginRate                  string `json:"marginRate"`
		GuaranteedStopLossOrderMode string `json:"guaranteedStopLossOrderMode"`
		Tags                        []struct {
			Type string `json:"type"`
			Name string `json:"name"`
		} `json:"tags"`
		Financing struct {
			LongRate            string `json:"longRate"`
			ShortRate           string `json:"shortRate"`
			FinancingDaysOfWeek []struct {
				DayOfWeek   string `json:"dayOfWeek"`
				DaysCharged int    `json:"daysCharged"`
			} `json:"financingDaysOfWeek"`
		} `json:"financing"`
		MinimumGuaranteedStopLossDistance       string `json:"minimumGuaranteedStopLossDistance,omitempty"`
		GuaranteedStopLossOrderExecutionPremium string `json:"guaranteedStopLossOrderExecutionPremium,omitempty"`
		GuaranteedStopLossOrderLevelRestriction struct {
			Volume     string `json:"volume"`
			PriceRange string `json:"priceRange"`
		} `json:"guaranteedStopLossOrderLevelRestriction,omitempty"`
	} `json:"instruments"`
	LastTransactionID string `json:"lastTransactionID"`
}
