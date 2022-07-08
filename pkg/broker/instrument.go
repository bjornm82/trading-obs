package broker

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	ENDPOINT_INSTRUMENTS_CANDLES         = "/instruments/%s/candles?count=%d&price=BAM&granularity=%s"
	ENDPOINT_INSTRUMENTS_PRECISE_CANDLES = "/instruments/%s/candles?price=BAM&from=%s&to=%s&granularity=%s"
)

// S5	5 second candlesticks, minute alignment
// S10	10 second candlesticks, minute alignment
// S15	15 second candlesticks, minute alignment
// S30	30 second candlesticks, minute alignment
// M1	1 minute candlesticks, minute alignment
// M2	2 minute candlesticks, hour alignment
// M4	4 minute candlesticks, hour alignment
// M5	5 minute candlesticks, hour alignment
// M10	10 minute candlesticks, hour alignment
// M15	15 minute candlesticks, hour alignment
// M30	30 minute candlesticks, hour alignment
// H1	1 hour candlesticks, hour alignment
// H2	2 hour candlesticks, day alignment
// H3	3 hour candlesticks, day alignment
// H4	4 hour candlesticks, day alignment
// H6	6 hour candlesticks, day alignment
// H8	8 hour candlesticks, day alignment
// H12	12 hour candlesticks, day alignment
// D	1 day candlesticks, day alignment
// W	1 week candlesticks, aligned to start of week
// M	1 month candlesticks, aligned to first day of the month

func (r *Repo) GetPreciseCandles(instrument string, date_from, date_to time.Time, granularity string) (CandlesResponse, error) {
	o := CandlesResponse{}
	ok, err := r.InstrumentExists(instrument)
	if err != nil {
		return o, err
	}
	if !ok {
		return o, fmt.Errorf("instrument given: %s does not exist for the account", instrument)
	}

	from := date_from.Format("2006-01-02T15:04:05")
	to := date_to.Format("2006-01-02T15:04:05")

	// ?price=BAM&from=2022-07-06T06:00:00.000000000Z&to=2022-07-06T07:00:00.000000000Z&granularity=M5
	var querystring = fmt.Sprintf(
		ENDPOINT_INSTRUMENTS_PRECISE_CANDLES,
		instrument,
		from,
		to,
		granularity,
	)
	b, err := r.cl.Read(querystring)
	if err != nil {
		return o, err
	}

	err = json.Unmarshal(b, &o)
	if err != nil {
		return o, err
	}

	return o, nil
}

func (r *Repo) GetCandles(instrument string, count int, granularity string) (CandlesResponse, error) {
	o := CandlesResponse{}
	ok, err := r.InstrumentExists(instrument)
	if err != nil {
		return o, err
	}
	if !ok {
		return o, fmt.Errorf("instrument given: %s does not exist for the account", instrument)
	}

	b, err := r.cl.Read(fmt.Sprintf(
		ENDPOINT_INSTRUMENTS_CANDLES,
		instrument,
		count,
		granularity,
	))

	if err != nil {
		return o, err
	}

	err = json.Unmarshal(b, &o)
	if err != nil {
		return o, err
	}

	return o, nil
}

type CandlesResponse struct {
	Instrument  string `json:"instrument"`
	Granularity string `json:"granularity"`
	Candles     []struct {
		Complete bool      `json:"complete"`
		Volume   int       `json:"volume"`
		Time     time.Time `json:"time"`
		Bid      struct {
			O string `json:"o"`
			H string `json:"h"`
			L string `json:"l"`
			C string `json:"c"`
		} `json:"bid"`
		Mid struct {
			O string `json:"o"`
			H string `json:"h"`
			L string `json:"l"`
			C string `json:"c"`
		} `json:"mid"`
		Ask struct {
			O string `json:"o"`
			H string `json:"h"`
			L string `json:"l"`
			C string `json:"c"`
		} `json:"ask"`
	} `json:"candles"`
}
