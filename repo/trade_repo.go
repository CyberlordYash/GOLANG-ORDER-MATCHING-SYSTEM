package repo

import (
	"context"
	"database/sql"

	"golang-order-matching-system/models"
)

type TradeRepo struct{ db *sql.DB }

func NewTradeRepo(db *sql.DB) *TradeRepo { return &TradeRepo{db} }

func (r *TradeRepo) Insert(ctx context.Context, t *models.Trade) error {
	_, err := r.db.ExecContext(ctx, `
	  INSERT INTO trades(symbol, buy_order_id, sell_order_id, price, qty)
	  VALUES (?,?,?,?,?)`,
		t.Symbol, t.BuyOrderID, t.SellOrderID, t.Price, t.Qty)
	return err
}
func (r *TradeRepo) InsertTx(ctx context.Context, tx *sql.Tx, t *models.Trade) error {
	_, err := tx.ExecContext(ctx, `
	  INSERT INTO trades(symbol, buy_order_id, sell_order_id, price, qty)
	  VALUES (?,?,?,?,?)`,
		t.Symbol, t.BuyOrderID, t.SellOrderID, t.Price, t.Qty)
	return err
}

func (r *TradeRepo) DB() *sql.DB {
	return r.db
}

func (r *TradeRepo) ListRecent(ctx context.Context, symbol string, limit int) ([]models.Trade, error) {
	if limit <= 0 || limit > 1000 {
		limit = 100 
	}

	var rows *sql.Rows
	var err error

	if symbol == "" {
		rows, err = r.db.QueryContext(ctx, `
		  SELECT id, symbol, buy_order_id, sell_order_id, price, qty
		  FROM trades
		  ORDER BY ts DESC
		  LIMIT ?`, limit)
	} else {
		rows, err = r.db.QueryContext(ctx, `
		  SELECT id, symbol, buy_order_id, sell_order_id, price, qty
		  FROM trades
		  WHERE symbol = ?
		  ORDER BY ts DESC
		  LIMIT ?`, symbol, limit)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Trade
	for rows.Next() {
		var t models.Trade
		if err := rows.Scan(&t.ID, &t.Symbol, &t.BuyOrderID,
			&t.SellOrderID, &t.Price, &t.Qty); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}