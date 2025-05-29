package engine


// --- public -----------------------------------------------------------------

type OrderSide string

const (
	Buy  OrderSide = "buy"
	Sell OrderSide = "sell"
)

// Order sent to the engine (no pointers to keep it simple).
type Order struct {
	ID      int64
	Symbol  string
	Side    OrderSide
	IsLimit bool
	Price   float64 // 0 for market order
	Qty     int64
	Ts      int64 // unix nano for FIFO
}

type Fill struct {
	TakerID int64
	MakerID int64
	Price   float64
	Qty     int64
}

// --- internal ---------------------------------------------------------------

type priceLevel struct {
	price  float64
	orders []*Order // FIFO queue
}

// generic heap (min); custom Less for buys to invert it
type orderPQ struct {
	levels []*priceLevel
	isBuy  bool
}

func (pq orderPQ) Len() int { return len(pq.levels) }
func (pq orderPQ) Less(i, j int) bool {
	if pq.isBuy {
		return pq.levels[i].price > pq.levels[j].price // max-heap
	}
	return pq.levels[i].price < pq.levels[j].price // min-heap
}
func (pq orderPQ) Swap(i, j int) { pq.levels[i], pq.levels[j] = pq.levels[j], pq.levels[i] }
func (pq *orderPQ) Push(x any)   { pq.levels = append(pq.levels, x.(*priceLevel)) }
func (pq *orderPQ) Pop() (v any) {
	n := len(pq.levels)
	v = pq.levels[n-1]
	pq.levels = pq.levels[:n-1]
	return
}

// --- orderBook --------------------------------------------------------------

type orderBook struct {
	buys  orderPQ // best price first
	sells orderPQ // best price first
}

func newOrderBook() *orderBook {
	return &orderBook{
		buys:  orderPQ{isBuy: true},
		sells: orderPQ{isBuy: false},
	}
}
