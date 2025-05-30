package api

import (
	"time"

	"golang-order-matching-system/engine"
	"golang-order-matching-system/models"
)

type PlaceOrderReq struct {
	Symbol   string   `json:"symbol" binding:"required"`
	Side     string   `json:"side"   binding:"oneof=buy sell"`
	Type     string   `json:"type"   binding:"oneof=limit market"`
	Price    *float64 `json:"price,omitempty"` 
	Quantity int64    `json:"quantity" binding:"gt=0"`
}

type PlaceOrderResp struct {
	OrderID    int64         `json:"order_id"`
	Executions []engine.Fill `json:"executions"`
}

func (r *PlaceOrderReq) ToEngine(id int64) *engine.Order {
	return &engine.Order{
		ID:      id,
		Symbol:  r.Symbol,
		Side:    engine.OrderSide(r.Side),
		IsLimit: r.Type == "limit",
		Price:   getPrice(r),
		Qty:     r.Quantity,
		Ts:      time.Now().UnixNano(),
	}
}

func getPrice(r *PlaceOrderReq) float64 {
	if r.Price != nil {
		return *r.Price
	}
	return 0
}

func (r *PlaceOrderReq) ToModel(id int64) *models.Order {
	return &models.Order{
		ID:           id,
		Symbol:       r.Symbol,
		Side:         r.Side,
		Type:         r.Type,
		Price:        r.Price,
		QtyInitial:   r.Quantity,
		QtyRemaining: r.Quantity,
		Status:       "open",
	}
}
