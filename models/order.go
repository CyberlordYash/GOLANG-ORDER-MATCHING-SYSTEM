package models

type Order struct {
	ID           int64   `json:"id,omitempty"`
	Symbol       string  `json:"symbol"`
	Side         string  `json:"side"`           // buy | sell
	Type         string  `json:"type"`           // limit | market
	Price        *float64 `json:"price,omitempty"` // nil for market
	QtyInitial   int64   `json:"qty_initial"`
	QtyRemaining int64   `json:"qty_remaining"`
	Status       string  `json:"status"`         // open | filled | cancelled
}
