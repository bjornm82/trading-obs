package strategies

import (
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/xtordoir/goanda/models"
)

type ObsDax struct {
	Candles *models.Candles
}

type Rules struct {
	StrategySet bool
	OrderBuy    Orders
	OrderSell   Orders
}

type Orders struct {
	Entry      float64
	SL         float64
	TP         float64
	InPosition bool
}

func CreateRules(low, high float64) Rules {
	var profitPoints = high - low
	var riskPoints = profitPoints / 2

	return Rules{
		StrategySet: true,
		OrderBuy: Orders{
			Entry: high,
			SL:    high - riskPoints,
			TP:    high + profitPoints,
		},
		OrderSell: Orders{
			Entry: low,
			SL:    low + riskPoints,
			TP:    low - profitPoints,
		},
	}
}

func (o *ObsDax) Run(e Env) {
	var highest = float64(0)
	var lowest = float64(10000000)
	var rules = Rules{
		OrderBuy:  Orders{InPosition: false},
		OrderSell: Orders{InPosition: false},
	}
	for k, v := range o.Candles.Candles {
		if v.Time.Hour() == 6 && v.Time.Minute() == 55 {
			var hourPrior = o.Candles.Candles[k-11 : k+1]

			for _, h := range hourPrior {
				if h.Mid.H > highest {
					highest = h.Mid.H
				}
				if h.Mid.L < lowest {
					lowest = h.Mid.L
				}
			}
			rules = CreateRules(lowest, highest)
			log.Println(rules)
		}
		if v.Mid.H > rules.OrderBuy.Entry && !rules.OrderSell.InPosition && !rules.OrderBuy.InPosition && rules.StrategySet {
			rules.OrderBuy.InPosition = true
		}
		if v.Mid.L < rules.OrderBuy.SL && rules.OrderBuy.InPosition {
			spew.Dump(k)
			spew.Dump(v)
			log.Fatalln("LOSS on BUY order")
		}
		if v.Mid.L > rules.OrderBuy.TP && rules.OrderBuy.InPosition {
			spew.Dump(k)
			spew.Dump(v)
			log.Println(rules)
			log.Fatalln("WON on BUY order")
		}
		if v.Mid.L < rules.OrderSell.Entry && !rules.OrderSell.InPosition && !rules.OrderBuy.InPosition && rules.StrategySet {
			rules.OrderSell.InPosition = true
		}
		if v.Mid.L > rules.OrderSell.SL && rules.OrderSell.InPosition {
			spew.Dump(k)
			spew.Dump(v)
			log.Fatalln("LOSS on SELL order")
		}
		if v.Mid.L < rules.OrderSell.TP && rules.OrderSell.InPosition {
			spew.Dump(k)
			spew.Dump(v)
			log.Fatalln("WON on SELL order")
		}
	}
	spew.Dump(e)
}
