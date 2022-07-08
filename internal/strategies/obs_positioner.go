package strategies

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/bjornm82/trading-obs/pkg/broker"
)

// Cron every n minutes
// If a position closes half of the points to take profit, update order with SL to that price

func (e *Env) ObsPosition() (bool, error) {
	p, err := e.Broker.GetAccount()
	if err != nil {
		return false, err
	}

	trades, err := clean(p)
	if err != nil {
		return false, err
	}

	for k, v := range trades {
		cresp, err := e.Broker.GetCandles(v.Instrument, 2, "M5")
		if err != nil {
			return false, fmt.Errorf("unable to fetch results: %s", err)
		}

		var last = Candle{}
		if cresp.Candles[1].Complete {
			o, err := toFloat(cresp.Candles[1].Ask.O)
			if err != nil {
				return false, err
			}
			h, err := toFloat(cresp.Candles[1].Ask.H)
			if err != nil {
				return false, err
			}
			l, err := toFloat(cresp.Candles[1].Ask.L)
			if err != nil {
				return false, err
			}
			c, err := toFloat(cresp.Candles[1].Ask.C)
			if err != nil {
				return false, err
			}

			last.Ask.O = o
			last.Ask.H = h
			last.Ask.L = l
			last.Ask.C = c

			o, err = toFloat(cresp.Candles[1].Bid.O)
			if err != nil {
				return false, err
			}
			h, err = toFloat(cresp.Candles[1].Bid.H)
			if err != nil {
				return false, err
			}
			l, err = toFloat(cresp.Candles[1].Bid.L)
			if err != nil {
				return false, err
			}
			c, err = toFloat(cresp.Candles[1].Bid.C)
			if err != nil {
				return false, err
			}

			last.Bid.O = o
			last.Bid.H = h
			last.Bid.L = l
			last.Bid.C = c

			o, err = toFloat(cresp.Candles[1].Mid.O)
			if err != nil {
				return false, err
			}
			h, err = toFloat(cresp.Candles[1].Mid.H)
			if err != nil {
				return false, err
			}
			l, err = toFloat(cresp.Candles[1].Mid.L)
			if err != nil {
				return false, err
			}
			c, err = toFloat(cresp.Candles[1].Mid.C)
			if err != nil {
				return false, err
			}

			last.Mid.O = o
			last.Mid.H = h
			last.Mid.L = l
			last.Mid.C = c

			last.Volume = cresp.Candles[1].Volume
			last.Time = cresp.Candles[1].Time
			trades[k].LastCandle = last
		} else if cresp.Candles[0].Complete {
			o, err := toFloat(cresp.Candles[0].Ask.O)
			if err != nil {
				return false, err
			}
			h, err := toFloat(cresp.Candles[0].Ask.H)
			if err != nil {
				return false, err
			}
			l, err := toFloat(cresp.Candles[0].Ask.L)
			if err != nil {
				return false, err
			}
			c, err := toFloat(cresp.Candles[0].Ask.C)
			if err != nil {
				return false, err
			}

			last.Ask.O = o
			last.Ask.H = h
			last.Ask.L = l
			last.Ask.C = c

			o, err = toFloat(cresp.Candles[0].Bid.O)
			if err != nil {
				return false, err
			}
			h, err = toFloat(cresp.Candles[0].Bid.H)
			if err != nil {
				return false, err
			}
			l, err = toFloat(cresp.Candles[0].Bid.L)
			if err != nil {
				return false, err
			}
			c, err = toFloat(cresp.Candles[0].Bid.C)
			if err != nil {
				return false, err
			}

			last.Bid.O = o
			last.Bid.H = h
			last.Bid.L = l
			last.Bid.C = c

			o, err = toFloat(cresp.Candles[0].Mid.O)
			if err != nil {
				return false, err
			}
			h, err = toFloat(cresp.Candles[0].Mid.H)
			if err != nil {
				return false, err
			}
			l, err = toFloat(cresp.Candles[0].Mid.L)
			if err != nil {
				return false, err
			}
			c, err = toFloat(cresp.Candles[0].Mid.C)
			if err != nil {
				return false, err
			}

			last.Mid.O = o
			last.Mid.H = h
			last.Mid.L = l
			last.Mid.C = c
			last.Volume = cresp.Candles[0].Volume
			last.Time = cresp.Candles[0].Time
			trades[k].LastCandle = last
		} else {
			return false, fmt.Errorf("unable to fetch last completed candle")
		}
	}

	calcTrades := calculate(trades)

	// MOVING STOP LOSS PUT
	// https://api-fxpractice.oanda.com/v3/accounts/101-004-8979399-001/trades/512/orders
	// {"stopLoss":{"price":"7046.8"}}

	for k, v := range calcTrades {
		if v.BreakEven {
			log.Printf(
				"updating stoploss of trade ID: %s and order ID: %s from %.1f to %.1f",
				v.TradeID,
				v.SL.ID,
				v.SL.Price,
				v.Price,
			)

			_, err := e.Broker.PutSLToBreakEven(
				fmt.Sprintf("%.1f", v.Price),
				v.SL.ID,
				v.TradeID,
			)

			if err != nil {
				log.Println(err)
				log.Println("something wrong on the update order, need to retry again")
			}

			calcTrades = append(calcTrades[:k], calcTrades[k+1:]...)
		}
	}

	return true, nil
}

func calculate(t Trades) Trades {
	for k, v := range t {
		if v.Price != v.SL.Price {
			if v.Direction == "BUY" {
				var diff = v.Price - v.SL.Price
				var break_even_price = v.Price + diff
				t[k].BreakEven = false
				if v.LastCandle.Bid.C >= break_even_price {
					t[k].BreakEven = true
				}
			}
			if v.Direction == "SELL" {
				var diff = v.SL.Price - v.Price
				var break_even_price = v.Price - diff
				t[k].BreakEven = false
				if v.LastCandle.Ask.C <= break_even_price {
					t[k].BreakEven = true
				}
			}
		}
	}

	return t
}

func clean(ar broker.AccountResponse) (Trades, error) {
	trades := Trades{}
	for _, acTrades := range ar.Account.Trades {
		trade := Trade{}
		trade.Instrument = acTrades.Instrument
		trade.TradeID = acTrades.ID
		u, err := toFloat(acTrades.CurrentUnits)
		if err != nil {
			return trades, err
		}
		trade.Units = u
		if u > 0 {
			trade.Direction = "BUY"
		} else if u < 0 {
			trade.Direction = "SELL"
		} else {
			return trades, fmt.Errorf("unable to determine direction derived from value %f", u)
		}

		trade.State = acTrades.State
		p, err := toFloat(acTrades.Price)
		if err != nil {
			return trades, err
		}
		trade.Price = p

		for _, acOrders := range ar.Account.Orders {
			if acOrders.TradeID == acTrades.ID {
				if acOrders.ID == acTrades.TakeProfitOrderID {
					p, err := toFloat(acOrders.Price)
					if err != nil {
						return trades, err
					}
					o := Order{
						ID:    acOrders.ID,
						Price: p,
						Type:  acOrders.Type,
						State: acOrders.State,
					}
					trade.TP = o
				}

				if acOrders.ID == acTrades.StopLossOrderID {
					p, err := toFloat(acOrders.Price)
					if err != nil {
						return trades, err
					}
					o := Order{
						ID:    acOrders.ID,
						Price: p,
						Type:  acOrders.Type,
						State: acOrders.State,
					}
					trade.SL = o
				}
			}
		}

		trades = append(trades, trade)
	}

	return trades, nil
}

func toFloat(value string) (float64, error) {
	const bitSize = 64
	return strconv.ParseFloat(value, bitSize)
}

type Trades []Trade

type Trade struct {
	Direction  string  `json:"direction"`
	TradeID    string  `json:"trade_id"`
	Instrument string  `json:"instrument"`
	Price      float64 `json:"price"`
	Units      float64 `json:"units"`
	SL         Order   `json:"sl"`
	TP         Order   `json:"tp"`
	State      string  `json:"state"`
	LastCandle Candle  `json:"last_candle"`
	BreakEven  bool    `json:"break_even"`
}

type Order struct {
	ID    string  `json:"id"`
	Type  string  `json:"type"`
	Price float64 `json:"price"`
	State string  `json:"state"`
}

type Candle struct {
	Volume int       `json:"volume"`
	Time   time.Time `json:"time"`
	Bid    struct {
		O float64 `json:"o"`
		H float64 `json:"h"`
		L float64 `json:"l"`
		C float64 `json:"c"`
	} `json:"bid"`
	Mid struct {
		O float64 `json:"o"`
		H float64 `json:"h"`
		L float64 `json:"l"`
		C float64 `json:"c"`
	} `json:"mid"`
	Ask struct {
		O float64 `json:"o"`
		H float64 `json:"h"`
		L float64 `json:"l"`
		C float64 `json:"c"`
	} `json:"ask"`
}
