package account

import (
	"time"
)

const START_CAPITAL = 10000
const SPREAD = 1.7

func New() BankManager {
	return &Bank{
		Capital: START_CAPITAL,
		Spread:  SPREAD,
	}
}

type BankManager interface {
	UpdateCapital(value float64)
	GetCapital() float64
	GetOrders() Orders
	AddOrder(direction, symbol string, amount, price, stopLoss, takeProfit float64) bool
}

type Bank struct {
	Capital float64 `json:"capital"`
	Spread  float64 `json:"spread"`
	Orders  Orders  `json:"orders"`
}

func (b *Bank) UpdateCapital(value float64) {
	b.Capital = b.Capital + value
}

func (b *Bank) GetCapital() float64 {
	return b.Capital
}

func (b *Bank) GetOrders() Orders {
	return b.Orders
}

func (b *Bank) AddOrder(direction, symbol string, amount, price, stopLoss, takeProfit float64) bool {
	return b.Orders.Add(direction, symbol, amount, price, stopLoss, takeProfit)
}

type Orders struct {
	Order []Order
}

func (o Orders) Len() int {
	return len(o.Order)
}

func (o Orders) Less(i, j int) bool {
	return o.Order[i].ID < o.Order[j].ID
}

func (o Orders) Swap(i, j int) {
	o.Order[i], o.Order[j] = o.Order[j], o.Order[i]
}

func (o *Orders) Add(direction, symbol string, amount, price, stopLoss, takeProfit float64) bool {
	if o.LastOrder().Position {
		return false
	}

	order := Order{
		ID:        o.LastOrder().ID + 1,
		Direction: direction,
		Symbol:    symbol,
		Amount:    amount,
		StopLoss:  stopLoss,
		Open: Contract{
			Timestamp: time.Now(),
			Price:     price,
		},
	}

	o.Order = append(o.Order, order)

	return true
}

func (o *Orders) LastOrder() Order {
	return Order{}
}

type Order struct {
	ID         int      `json:"id"`
	Direction  string   `json:"direction"`
	Symbol     string   `json:"symbol"`
	Amount     float64  `json:"amount"`
	Open       Contract `json:"open"`
	Close      Contract `json:"close"`
	Position   bool     `json:"position"`
	StopLoss   float64  `json:"stop_loss"`
	TakeProfit float64  `json:"take_profit"`
}

type Contract struct {
	Timestamp time.Time `json:"timestamp"`
	Price     float64   `json:"price"`
}

func (b *Order) InPosition() bool {
	return b.Position
}

func (b *Order) CloseOrder(price float64) bool {
	return b.InPosition()
}

func (b *Order) Sell(amount, price float64) bool {
	return !b.InPosition()
}

func (b *Order) Buy(amount, price float64) bool {
	return !b.InPosition()
}
