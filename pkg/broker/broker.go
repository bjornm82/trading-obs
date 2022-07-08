package broker

import "time"

type BrokerManager interface {
	GetOrders() (OrdersResponse, error)
	GetPendingOrders() (PendingOrdersResponse, error)
	CreateOrder(units int, instrument, order_type string) (OrderResponse, error)
	CreateSlTpMarketOrder(units int, instrument string, price, sl, tp float64) (OrderResponse, error)
	GetOpenTrades() (OpenTradesResponse, error)
	GetPosition(instrument string) (PositionResponse, error)
	GetPositions() (PositionsResponse, error)
	GetInstruments() (InstrumentResponse, error)
	GetAccount() (AccountResponse, error)
	PutSLToBreakEven(price, order_id, trade_id string) (UpdateOrderResponse, error)
	GetCandles(instrument string, count int, granularity string) (CandlesResponse, error)
	GetPreciseCandles(instrument string, date_from, date_to time.Time, granularity string) (CandlesResponse, error)
	TickStream(instruments []string, tchan chan Tick, end chan bool)
}

type Repo struct {
	cl ClientManager
}

func New(cl ClientManager) BrokerManager {
	return &Repo{
		cl: cl,
	}
}
