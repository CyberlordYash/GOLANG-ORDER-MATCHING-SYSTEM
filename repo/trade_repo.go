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

