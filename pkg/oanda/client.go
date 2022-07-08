package oanda

import (
	"github.com/xtordoir/goanda/api"
	"github.com/xtordoir/goanda/models"
)

type OandaManager interface {
	GetPricing(symbol string, timeframe string) (*models.Candles, error)
}

func New(api_url, stream_url, token, account, app_name string) Client {
	ctx := api.Context{
		ApiURL:       api_url,
		StreamApiURL: stream_url,
		Token:        token,
		Account:      account,
		Application:  app_name,
	}

	return Client{ctx: ctx}
}

type Client struct {
	ctx api.Context
}
