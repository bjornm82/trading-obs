package strategies

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/bjornm82/trading-obs/pkg/broker"
)

var instruments = []string{
	"US30_USD",
	"NAS100_USD",
	"SPX500_USD",
	"DE30_EUR",
	"UK100_GBP",
	"FR40_EUR",
	"US2000_USD",
}

type Instruments struct {
	Instruments []Params
}

type Params struct {
	Name string  `json:"name"`
	Max  float64 `json:"max"`
	Min  float64 `json:"min"`
}

func (e *Env) tickProcessor(c chan broker.Tick, done chan bool, i Instruments) {
	for {
		tick := <-c
		for k, v := range i.Instruments {
			if tick.Instrument == v.Name {
				var mid = (tick.Bid + tick.Ask) / 2
				if mid > (v.Max + 1) {
					a, err := e.Broker.GetAccount()
					if err != nil {
						log.Println(err)
						return
					}

					log.Println("KICKOFF LONG FOR: ")
					log.Println("Mid price tick: ")
					log.Println(mid)
					log.Println("Instrument")
					log.Println(v)
					log.Print("POINTS AT RISK: ")
					var riskPoints = (v.Max - v.Min) / 2
					log.Println(riskPoints)
					log.Print("TOTAL BALANCE NOW: ")
					balance, err := strconv.ParseFloat(a.Account.Balance, 64)
					if err != nil {
						log.Println(err)
						return
					}
					log.Println(balance)
					var one_percent = balance / 100
					var units = one_percent / riskPoints
					var diff = v.Max - v.Min
					var sl = mid - (diff / 2)
					var tp = mid + diff
					log.Print("UNITS")
					log.Println(units)
					log.Print("SL: ")
					log.Println(sl)
					log.Print("TP: ")
					log.Println(tp)
					log.Println("")
					log.Println("")

					// e.Broker.CreateSlTpMarketOrder()
					i.Instruments = append(i.Instruments[:k], i.Instruments[k+1:]...)
				}
				if mid < (v.Min - 1) {
					log.Println("KICKOFF SHORT FOR: ")
					log.Println("Mid price tick: ")
					log.Println(mid)
					log.Println("Instrument")
					log.Println(v)
					log.Println("")
					log.Println("")
					i.Instruments = append(i.Instruments[:k], i.Instruments[k+1:]...)
				}
			}
		}
		if len(i.Instruments)-1 == 0 {
			log.Println("closing ticker, no instruments to be done")
			<-c
			<-done
		}
	}
}

// Run the order function which would
func (e *Env) Run(from, to time.Time) error {
	i := Instruments{}

	ar, err := e.Broker.GetAccount()
	if err != nil {
		return err
	}

	for k, in := range instruments {
		for _, t := range ar.Account.Trades {
			if t.Instrument == in {
				log.Println(fmt.Sprintf(
					"removing instrument %s due to open trade with ID: %s",
					t.Instrument,
					t.ID,
				))
				instruments = append(instruments[:k], instruments[k+1:]...)
			}
		}
	}

	for _, instrument := range instruments {
		resp, err := e.Broker.GetPreciseCandles(instrument, from, to, "M5")
		if err != nil {
			return err
		}
		par, err := GetHighLow(resp)
		if err != nil {
			return err
		}
		i.Instruments = append(i.Instruments, par)
	}

	// channels for data
	tchan := make(chan broker.Tick)
	done := make(chan bool)

	// start processors for data
	go e.tickProcessor(tchan, done, i)

	log.Println("start ticking")
	e.Broker.TickStream(instruments, tchan, done)
	log.Println("closing tchan")

	log.Println(i)

	return nil
}

func GetHighLow(resp broker.CandlesResponse) (Params, error) {
	p := Params{}
	var highest float64
	var lowest float64
	lowest = 1000000000000
	for _, v := range resp.Candles {
		log.Println(v.Time)
		high, err := strconv.ParseFloat(v.Mid.H, 64)
		if err != nil {
			return p, err
		}
		if high > highest {
			highest = high
		}
		low, err := strconv.ParseFloat(v.Mid.L, 64)
		if err != nil {
			return p, err
		}
		if low < lowest {
			lowest = low
		}
	}

	p.Max = highest
	p.Min = lowest
	p.Name = resp.Instrument

	return p, nil
}
