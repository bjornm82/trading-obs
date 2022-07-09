package broker

import "time"

type BrokerManager interface {
	GetOrders() (OrdersResponse, error)
	GetPendingOrders() (PendingOrdersResponse, error)
	CreateOrder(units float64, instrument, order_type string) (OrderResponse, error)
	CreateSlTpMarketOrder(units float64, instrument string, price, sl, tp float64) (OrderResponse, error)
	GetOpenTrades() (OpenTradesResponse, error)
	GetPosition(instrument string) (PositionResponse, error)
	GetPositions() (PositionsResponse, error)
	GetInstruments() (InstrumentResponse, error)
	GetAccount() (AccountResponse, error)
	PutSLToBreakEven(price, order_id, trade_id string) (UpdateOrderResponse, error)
	GetCandles(instrument string, count int, granularity string) (CandlesResponse, error)
	GetPreciseCandles(instrument string, date_from, date_to time.Time, granularity string) (CandlesResponse, error)
	TickStream(instruments []string, tchan chan Tick, done chan bool)
	GetSemiPriceStream(instrument []string) (SemiTick, error)
}

type Repo struct {
	cl ClientManager
}

func New(cl ClientManager) BrokerManager {
	return &Repo{
		cl: cl,
	}
}
