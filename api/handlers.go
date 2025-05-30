package api

import (
	"context"
	"net/http"
    "database/sql"
	"time"
    "strconv"
	"github.com/gin-gonic/gin"
	"golang-order-matching-system/engine"
	"golang-order-matching-system/models"
	"golang-order-matching-system/repo"
)

type Handler struct {
	Eng *engine.Engine
	OR  *repo.OrderRepo
	TR  *repo.TradeRepo
}

type BookLevel struct {
	Price float64 `json:"price"`
	Qty   int64   `json:"qty"`
}

func (h *Handler) PlaceOrder(c *gin.Context) {
    var req PlaceOrderReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx := context.Background()
    tx, err := h.OR.DB().BeginTx(ctx, nil)
    if err != nil { /* handle */ }
    defer tx.Rollback()

    // 1. create *empty* order row to get real ID
    model := req.ToModel(0)                // id will be filled later
    id, err := h.OR.InsertTx(ctx, tx, model)
    if err != nil { /* handle */ }

    // 2. run engine with the real ID
    engOrder := req.ToEngine(id)
    fills := h.Eng.Process(engOrder)

    // 3. update remaining qty / status in the same tx
    if err := h.OR.UpdateRemainingTx(ctx, tx,
        id, engOrder.Qty, statusFromQty(engOrder.Qty)); err != nil {
        /* handle */
    }

    // 4. persist trades with correct IDs
    for _, f := range fills {
        t := &models.Trade{
            Symbol:      engOrder.Symbol,
            BuyOrderID:  chooseBuyID(f, req.Side, id),
            SellOrderID: chooseSellID(f, req.Side, id),
            Price:       f.Price,
            Qty:         f.Qty,
        }
        if err := h.TR.InsertTx(ctx, tx, t); err != nil { /* handle */ }
    }

    if err := tx.Commit(); err != nil { /* handle */ }

    c.JSON(http.StatusOK, PlaceOrderResp{
        OrderID:    id,
        Executions: fills,
    })
}

func statusFromQty(remaining int64) string {
    if remaining == 0 {
        return "filled"
    }
    return "open"
}

func chooseBuyID(f engine.Fill, side string, takerID int64) int64 {
    if side == "buy" {
        return takerID        // taker is the buy side
    }
    return f.MakerID          // maker was the buy side
}

func chooseSellID(f engine.Fill, side string, takerID int64) int64 {
    if side == "sell" {
        return takerID
    }
    return f.MakerID
}



func (h *Handler) ListTrades(c *gin.Context) {
	symbol := c.Query("symbol")            // optional filter
	limitStr := c.DefaultQuery("limit", "100")

	limit, err := strconv.Atoi(limitStr)   // allow ?limit=200 etc.
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be an integer"})
		return
	}

	ctx := context.Background()
	trades, err := h.TR.ListRecent(ctx, symbol, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, trades)
}







func (h *Handler) CancelOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	ctx := context.Background()
	if err := h.OR.CancelOrder(ctx, id); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "order not open or not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "cancelled", "order_id": id})
}



func (h *Handler) GetOrderBook(c *gin.Context) {
	symbol := c.Query("symbol")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "symbol query param required"})
		return
	}

	ctx := context.Background()
	bids, asks, err := h.OR.GetOrderBook(ctx, symbol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"symbol": symbol,
		"bids":   bids,
		"asks":   asks,
	})
}



func timeNowNano() int64 { return time.Now().UnixNano() }
