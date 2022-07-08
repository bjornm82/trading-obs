package broker

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	ENDPOINT_STREAM_PRICING = "/accounts/%s/pricing/stream"
)

// TickStream starts a stream of ticks, hiding the Prices structs which are autoRestarted
func (r *Repo) TickStream(instruments []string, tchan chan Tick, done chan bool) {
	pchan := make(chan ClientPrice)
	// AutoRestart for PricingStream
	go autoRestart("PricingStream", 0, func() { r.PricingStream(instruments, pchan) })

	fmt.Println("Starting loop on Prices")

	for {
		price := <-pchan
		tchan <- ClientPrice2Tick(&price)
	}
}

// AutoRestart for the PricingStream function as connection reset can result in panic
func autoRestart(name string, nPanics int64, f func()) {
	defer func() {
		if v := recover(); v != nil {
			// A panic is detected.
			log.Printf("%s is crashed. Panic #%d. Restarting in 5 seconds.", name, nPanics+1)
			time.Sleep(5 * time.Second)
			go autoRestart(name, nPanics+1, f) // restart
		}
	}()
	f()
}

// PricingStream starts a stream of prices
func (r *Repo) PricingStream(instruments []string, pchan chan ClientPrice) {
	a, err := r.cl.GetAccountID()
	if err != nil {
		log.Fatalln(err)
	}

	url := fmt.Sprintf(ENDPOINT_STREAM_PRICING, a)
	qurl := url + "?instruments=" + strings.Join(instruments, ",")

	client := &http.Client{}
	req, _ := http.NewRequest("GET", PRACTICE_STREAMING_HOST+qurl, nil)
	req.Header.Add("Authorization", "Bearer "+r.cl.GetToken())
	response, err := client.Do(req)

	// b, err := r.cl.ReadStream(qurl)
	// var resp = bytes.NewReader(b)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		reader := bufio.NewReader(response.Body)
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				log.Println(err)
				panic("Connection on clientStream is lost")
			}
			var p ClientPrice
			json.Unmarshal([]byte(line), &p)
			if p.Type == "HEARTBEAT" {
			} else {
				pchan <- p
			}
		}
	}
	//return nil
}

// Tick is a bid/ask for an instrument at a given time
type Tick struct {
	Instrument string
	Time       time.Time
	Bid        float64
	Ask        float64
}

// Price of a Tick (avergae Bid-Ask)
func (tick *Tick) Price() float64 {
	return (tick.Ask + tick.Bid) / 2
}

// ClientPrice2Tick converts a ClientPrice to a Tick, by taking the first Bid and Ask
func ClientPrice2Tick(price *ClientPrice) Tick {
	return Tick{
		Instrument: price.Instrument,
		Time:       price.Time,
		Bid:        (*price).Bids[0].Price,
		Ask:        (*price).Asks[0].Price,
	}
}

// ClientPrice is the price of an instrument
type ClientPrice struct {
	Instrument string        `json:"instrument"`
	Type       string        `json:"type"`
	Time       time.Time     `json:"time"`
	Bids       []PriceBucket `json:"bids"`
	Asks       []PriceBucket `json:"asks"`
}

// PricingHeartbeat is a heartbeat to keep connection alive
type PricingHeartbeat struct {
	Type string    `json:"type"`
	Time time.Time `json:"time"`
}

// Prices is the object response from GetPricing call
type Prices struct {
	Prices []ClientPrice `json:"prices"`
}

// PriceBucket is a type for Bids or Asks
type PriceBucket struct {
	Price     float64 `json:"price,string"`
	Liquidity int     `json:"liquidity"`
}
