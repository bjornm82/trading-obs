package oanda

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/xtordoir/goanda/models"
)

// Context is the api Context
type Context struct {
	ApiURL       string
	StreamApiURL string
	Token        string
	Account      string
	Application  string
}

// CreateAPI Creates an api instance from the Context
func (context *Context) CreateAPI() API {
	return API{
		context: *context,
	}
}

// API is an api instance with a context to call endpoints
type API struct {
	context Context
}

func (cl Client) GetPricing(symbol string, timeframe string) (*models.Candles, error) {
	api := cl.ctx.CreateAPI()
	return api.GetCandles(symbol, 100, timeframe)
}

func (api *API) GetInstruments() ([]byte, error) {
	client := &http.Client{}
	apiURL := api.context.ApiURL
	token := api.context.Token
	account := api.context.Account
	url := fmt.Sprintf(
		"%s/v3/accounts/%s/instruments",
		apiURL,
		account,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d expected but %d given", http.StatusOK, response.StatusCode)
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetCandles fetches a number of candles for a given instrument and granularity
func (api *API) GetCandles(instrument string, num int, granularity string) (*models.Candles, error) {
	// TODO DEDUPLICATE THIS
	client := &http.Client{}
	apiURL := api.context.ApiURL
	token := api.context.Token
	qStr := fmt.Sprintf("?granularity=%s&count=%d", granularity, num)
	url := fmt.Sprintf(
		"%s/v3/instruments/%s/candles%s",
		apiURL,
		instrument,
		qStr,
	)

	req, errr := http.NewRequest("GET", url, nil)
	if errr != nil {
		return nil, errr
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d expected but %d given", http.StatusOK, response.StatusCode)
	}

	data, errb := ioutil.ReadAll(response.Body)
	if errb != nil {
		return nil, errb
	}

	candles, errp := parseCandles(&data)

	return &candles, errp
}

func parseCandles(msg *[]byte) (models.Candles, error) {
	var p models.Candles
	err := json.Unmarshal(*msg, &p)
	return p, err
}
