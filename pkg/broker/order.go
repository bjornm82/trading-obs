package broker

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

const (
	ENDPOINT_ORDERS         = "/accounts/%s/orders"
	ENDPOINT_PENDING_ORDERS = "/accounts/%s/pendingOrders"
	ENDPOINT_PUT_ORDERS     = "/accounts/%s/orders/%s"
)

const (
	MARKET = iota
	LIMIT
	STOP
	TAKE_PROFIT
	STOP_LOSS
)

// CREATE SELL ORDER WITH SL AND TP WITH PRICE
// https://api-fxpractice.oanda.com/v3/accounts/101-004-8979399-001/orders
// {"order":{"instrument":"UK100_GBP","type":"MARKET","units":"-0.1","stopLossOnFill":{"price":"7093.1"},"takeProfitOnFill":{"price":"7053.1"}}}

// CREATE SELL ORDER WITH SL AND TP WITH DISTANCE AND PRICE
// https://api-fxpractice.oanda.com/v3/accounts/101-004-8979399-001/orders
// {"order":{"instrument":"UK100_GBP","type":"MARKET","units":"0.1","stopLossOnFill":{"distance":"4"},"takeProfitOnFill":{"price":"7100.2"}}}

// https://api-fxpractice.oanda.com/v3/accounts/101-004-8979399-001/positions/UK100_GBP/close
// {"longUnits":"ALL"}

// MOVING STOP LOSS PUT
// https://api-fxpractice.oanda.com/v3/accounts/101-004-8979399-001/trades/512/orders
// {"stopLoss":{"price":"7046.8"}}

// MOVING TAKE PROFIT PUT
// https://api-fxpractice.oanda.com/v3/accounts/101-004-8979399-001/trades/512/orders
// {"takeProfit":{"price":"7095.9"}}

func (r *Repo) CreateOrder(units int, instrument, order_type string) (OrderResponse, error) {
	resp := OrderResponse{}

	if instrument == "" {
		return resp, errors.New("instrument value can not be empty")
	}

	o := OrderPayload{
		Order: OrderBody{
			Units:        units,
			Instrument:   instrument,
			TimeInForce:  "FOK",
			Type:         order_type,
			PositionFill: "DEFAULT",
		},
	}

	b, err := json.Marshal(o)
	if err != nil {
		return resp, err
	}

	a, err := r.cl.GetAccountID()
	if err != nil {
		return resp, err
	}

	b, err = r.cl.Create(fmt.Sprintf(ENDPOINT_ORDERS, a), b)
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(b, &resp)
	return resp, err
}

func (r *Repo) CreateSlTpMarketOrder(units int, instrument string, price, sl, tp float64) (OrderResponse, error) {
	resp := OrderResponse{}

	if instrument == "" {
		return resp, errors.New("instrument value can not be empty")
	}

	o := OrderPayload{
		Order: OrderBody{
			Instrument:       instrument,
			Type:             "MARKET",
			Units:            units,
			StopLossOnFill:   &OnFill{Price: fmt.Sprintf("%.2f", sl)},
			TakeProfitOnFill: &OnFill{Price: fmt.Sprintf("%.2f", tp)},
		},
	}

	b, err := json.Marshal(o)
	if err != nil {
		return resp, err
	}

	a, err := r.cl.GetAccountID()
	if err != nil {
		return resp, err
	}

	b, err = r.cl.Create(fmt.Sprintf(ENDPOINT_ORDERS, a), b)
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(b, &resp)
	return resp, err
}

func (r *Repo) PutSLToBreakEven(price, order_id, trade_id string) (UpdateOrderResponse, error) {
	resp := UpdateOrderResponse{}

	a, err := r.cl.GetAccountID()
	if err != nil {
		return resp, err
	}

	sl := OrderSLBody{Order: Order{
		Price:   price,
		Type:    "STOP_LOSS",
		TradeID: trade_id,
	}}

	b, err := json.Marshal(sl)
	if err != nil {
		return resp, err
	}

	b, err = r.cl.Update(fmt.Sprintf(ENDPOINT_PUT_ORDERS, a, order_id), b)
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(b, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

type UpdateOrderResponse struct {
	OrderCancelTransaction struct {
		ID                string    `json:"id"`
		AccountID         string    `json:"accountID"`
		UserID            int       `json:"userID"`
		BatchID           string    `json:"batchID"`
		RequestID         string    `json:"requestID"`
		Time              time.Time `json:"time"`
		Type              string    `json:"type"`
		OrderID           string    `json:"orderID"`
		ReplacedByOrderID string    `json:"replacedByOrderID"`
		Reason            string    `json:"reason"`
	} `json:"orderCancelTransaction"`
	OrderCreateTransaction struct {
		ID                      string    `json:"id"`
		AccountID               string    `json:"accountID"`
		UserID                  int       `json:"userID"`
		BatchID                 string    `json:"batchID"`
		RequestID               string    `json:"requestID"`
		Time                    time.Time `json:"time"`
		Type                    string    `json:"type"`
		TradeID                 string    `json:"tradeID"`
		TimeInForce             string    `json:"timeInForce"`
		TriggerCondition        string    `json:"triggerCondition"`
		TriggerMode             string    `json:"triggerMode"`
		Price                   string    `json:"price"`
		Reason                  string    `json:"reason"`
		ReplacesOrderID         string    `json:"replacesOrderID"`
		CancellingTransactionID string    `json:"cancellingTransactionID"`
	} `json:"orderCreateTransaction"`
	RelatedTransactionIDs []string `json:"relatedTransactionIDs"`
	LastTransactionID     string   `json:"lastTransactionID"`
}

type OrderSLBody struct {
	Order Order `json:"order"`
}

type Order struct {
	Price   string `json:"price"`
	Type    string `json:"type"`
	TradeID string `json:"tradeID"`
}

type OrdersResponse struct {
	Orders []struct {
		ID               string    `json:"id"`
		CreateTime       time.Time `json:"createTime"`
		Type             string    `json:"type"`
		TradeID          string    `json:"tradeID"`
		Price            string    `json:"price"`
		TimeInForce      string    `json:"timeInForce"`
		TriggerCondition string    `json:"triggerCondition"`
		TriggerMode      string    `json:"triggerMode,omitempty"`
		State            string    `json:"state"`
	} `json:"orders"`
	LastTransactionID string `json:"lastTransactionID"`
}

type PendingOrdersResponse struct {
	LastTransactionID string `json:"lastTransactionID"`
	Orders            []struct {
		ClientExtensions struct {
			Comment string `json:"comment"`
			ID      string `json:"id"`
			Tag     string `json:"tag"`
		} `json:"clientExtensions,omitempty"`
		CreateTime       time.Time `json:"createTime"`
		ID               string    `json:"id"`
		Instrument       string    `json:"instrument,omitempty"`
		PartialFill      string    `json:"partialFill,omitempty"`
		PositionFill     string    `json:"positionFill,omitempty"`
		Price            string    `json:"price"`
		ReplacesOrderID  string    `json:"replacesOrderID,omitempty"`
		State            string    `json:"state"`
		TimeInForce      string    `json:"timeInForce"`
		TriggerCondition string    `json:"triggerCondition"`
		Type             string    `json:"type"`
		Units            string    `json:"units,omitempty"`
		StopLossOnFill   struct {
			Price       string `json:"price"`
			TimeInForce string `json:"timeInForce"`
		} `json:"stopLossOnFill,omitempty"`
		TradeID string `json:"tradeID,omitempty"`
	} `json:"orders"`
}

func (r *Repo) GetPendingOrders() (PendingOrdersResponse, error) {
	o := PendingOrdersResponse{}
	a, err := r.cl.GetAccountID()
	if err != nil {
		return o, err
	}
	b, err := r.cl.Read(fmt.Sprintf(ENDPOINT_PENDING_ORDERS, a))

	if err != nil {
		return o, err
	}

	err = json.Unmarshal(b, &o)
	if err != nil {
		return o, err
	}

	return o, nil
}

func (r *Repo) GetOrders() (OrdersResponse, error) {
	o := OrdersResponse{}
	a, err := r.cl.GetAccountID()
	if err != nil {
		return o, err
	}
	b, err := r.cl.Read(fmt.Sprintf(ENDPOINT_ORDERS, a))

	if err != nil {
		return o, err
	}

	err = json.Unmarshal(b, &o)
	if err != nil {
		return o, err
	}

	return o, nil
}

type OrderExtensions struct {
	Comment string `json:"comment,omitempty"`
	ID      string `json:"id,omitempty"`
	Tag     string `json:"tag,omitempty"`
}

type OnFill struct {
	TimeInForce string `json:"timeInForce,omitempty"`
	Price       string `json:"price,omitempty"` // must be a string for float precision
}

type OrderBody struct {
	Units            int              `json:"units"`
	Instrument       string           `json:"instrument"`
	TimeInForce      string           `json:"timeInForce"`
	Type             string           `json:"type"`
	PositionFill     string           `json:"positionFill,omitempty"`
	Price            string           `json:"price,omitempty"`
	TakeProfitOnFill *OnFill          `json:"takeProfitOnFill,omitempty"`
	StopLossOnFill   *OnFill          `json:"stopLossOnFill,omitempty"`
	ClientExtensions *OrderExtensions `json:"clientExtensions,omitempty"`
	TradeID          string           `json:"tradeId,omitempty"`
}

type OrderPayload struct {
	Order OrderBody `json:"order"`
}
type OrderResponse struct {
	LastTransactionID      string `json:"lastTransactionID"`
	OrderCreateTransaction struct {
		AccountID    string    `json:"accountID"`
		BatchID      string    `json:"batchID"`
		ID           string    `json:"id"`
		Instrument   string    `json:"instrument"`
		PositionFill string    `json:"positionFill"`
		Reason       string    `json:"reason"`
		Time         time.Time `json:"time"`
		TimeInForce  string    `json:"timeInForce"`
		Type         string    `json:"type"`
		Units        string    `json:"units"`
		UserID       int       `json:"userID"`
	} `json:"orderCreateTransaction"`
	OrderFillTransaction struct {
		AccountBalance string    `json:"accountBalance"`
		AccountID      string    `json:"accountID"`
		BatchID        string    `json:"batchID"`
		Financing      string    `json:"financing"`
		ID             string    `json:"id"`
		Instrument     string    `json:"instrument"`
		OrderID        string    `json:"orderID"`
		Pl             string    `json:"pl"`
		Price          string    `json:"price"`
		Reason         string    `json:"reason"`
		Time           time.Time `json:"time"`
		TradeOpened    struct {
			TradeID string `json:"tradeID"`
			Units   string `json:"units"`
		} `json:"tradeOpened"`
		Type   string `json:"type"`
		Units  string `json:"units"`
		UserID int    `json:"userID"`
	} `json:"orderFillTransaction"`
	RelatedTransactionIDs []string `json:"relatedTransactionIDs"`
}
