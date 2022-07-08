package broker

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	LIVE_HOST               = "https://api-fxtrade.oanda.com/v3"
	PRACTICE_HOST           = "https://api-fxpractice.oanda.com/v3"
	LIVE_STREAMING_HOST     = "https://stream-fxtrade.oanda.com/v3"
	PRACTICE_STREAMING_HOST = "https://stream-fxpractice.oanda.com/v3"
)

type Headers struct {
	contentType    string
	agent          string
	DatetimeFormat string
	auth           string
}

type ClientManager interface {
	GetAccountID() (string, error)
	GetToken() string
	Request(endpoint string) ([]byte, int, error)
	Read(endpoint string) ([]byte, error)
	ReadStream(endpoint string) ([]byte, error)
	Update(endpoint string, data []byte) ([]byte, error)
	Create(endpoint string, data []byte) ([]byte, error)
	// Delete(endpoint string) ([]byte, error)
}

type Client struct {
	hostname       string
	hostnameStream string
	port           int
	ssl            bool
	token          string
	accountID      string
	DatetimeFormat string
	headers        *Headers
}

const OANDA_AGENT string = "v20-golang/0.0.1"

func NewClient(accountID string, token string, live bool) ClientManager {
	var hostname = ""
	var hostnameStream = ""
	if live {
		hostname = LIVE_HOST
		hostnameStream = LIVE_STREAMING_HOST
	} else {
		hostname = PRACTICE_HOST
		hostnameStream = PRACTICE_STREAMING_HOST
	}

	var buf bytes.Buffer

	buf.WriteString("Bearer ")
	buf.WriteString(token)
	authHeader := buf.String()

	// Create headers for oanda to be used in requests
	headers := &Headers{
		contentType:    "application/json",
		agent:          OANDA_AGENT,
		DatetimeFormat: "RFC3339",
		auth:           authHeader,
	}

	return &Client{
		hostname:       hostname,
		hostnameStream: hostnameStream,
		port:           443,
		ssl:            true,
		token:          token,
		headers:        headers,
		accountID:      accountID,
	}
}

// TODO: Remove this and move the streaming properly
func (c *Client) GetToken() string {
	return c.token
}

func (c *Client) GetAccountID() (string, error) {
	if c.accountID == "" {
		return "", errors.New("no account ID set nothing to return")
	}
	return c.accountID, nil
}

func (c *Client) Request(endpoint string) ([]byte, int, error) {
	client := http.Client{
		Timeout: time.Second * 5,
	}

	url, err := createUrl(c.hostname, endpoint)
	if err != nil {
		return nil, 0, err
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}

	return makeRequest(c, endpoint, client, req)
}

func (c *Client) Read(endpoint string) ([]byte, error) {
	client := http.Client{
		Timeout: time.Second * 5,
	}

	url, err := createUrl(c.hostname, endpoint)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	b, code, err := makeRequest(c, endpoint, client, req)
	if err != nil {
		return nil, fmt.Errorf("unable to perform request %s", err)
	}

	if code != http.StatusOK {
		return nil, fmt.Errorf("request failed, expected %d but got value %d", http.StatusOK, code)
	}

	return b, nil
}

func (c *Client) ReadStream(endpoint string) ([]byte, error) {
	client := http.Client{
		Timeout: time.Second * 5,
	}

	url, err := createUrl(c.hostnameStream, endpoint)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	b, code, err := makeRequest(c, endpoint, client, req)
	if err != nil {
		return nil, fmt.Errorf("unable to perform request %s", err)
	}

	if code != http.StatusOK {
		return nil, fmt.Errorf("request failed, expected %d but got value %d", http.StatusOK, code)
	}

	return b, nil
}

func (c *Client) Create(endpoint string, data []byte) ([]byte, error) {
	client := http.Client{
		Timeout: time.Second * 5,
	}

	url, err := createUrl(c.hostname, endpoint)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	b, code, err := makeRequest(c, endpoint, client, req)
	if err != nil {
		return nil, fmt.Errorf("unable to perform request %s", err)
	}

	if code != http.StatusCreated {
		return nil, fmt.Errorf("request failed, expected %d but got value %d", http.StatusCreated, code)
	}

	return b, nil
}

func (c *Client) Update(endpoint string, data []byte) ([]byte, error) {
	client := http.Client{
		Timeout: time.Second * 5,
	}

	url, err := createUrl(c.hostname, endpoint)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("unable to perform request %s", err)
	}
	b, code, err := makeRequest(c, endpoint, client, req)
	if err != nil {
		return nil, fmt.Errorf("unable to perform request %s", err)
	}

	if code != http.StatusCreated {
		return nil, fmt.Errorf("request failed, expected %d but got value %d", http.StatusOK, code)
	}

	return b, nil
}

func createUrl(host string, endpoint string) (string, error) {
	if endpoint[0:1] != "/" {
		return "", errors.New("invalid URL: endpoint value should start with /")
	}

	return fmt.Sprintf("%s%s", host, endpoint), nil
}

func makeRequest(c *Client, endpoint string, client http.Client, req *http.Request) ([]byte, int, error) {
	req.Header.Set("User-Agent", c.headers.agent)
	req.Header.Set("Authorization", c.headers.auth)
	req.Header.Set("Content-Type", c.headers.contentType)

	res, err := client.Do(req)
	if err != nil {
		return nil, res.StatusCode, fmt.Errorf("unable to perform request to client: %s", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, res.StatusCode, fmt.Errorf("unable to perform request to client: %s", err)
	}

	err = checkApiErr(body, endpoint)
	if err != nil {
		return nil, res.StatusCode, fmt.Errorf("unable to perform request to client: %s", err)
	}
	return body, res.StatusCode, nil
}

func checkApiErr(body []byte, route string) error {
	bodyString := string(body[:])
	if strings.Contains(bodyString, "errorMessage") {
		return errors.New("\nOANDA API Error: " + bodyString + "\nOn route: " + route)
	}

	return nil
}
