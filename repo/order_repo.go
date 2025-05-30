package repo

import (
	"context"
	"database/sql"

	"golang-order-matching-system/models"
)

type OrderRepo struct{ db *sql.DB }

func NewOrderRepo(db *sql.DB) *OrderRepo { return &OrderRepo{db} }

func (r *OrderRepo) Insert(ctx context.Context, o *models.Order) error {
	_, err := r.db.ExecContext(ctx, `
	  INSERT INTO orders(symbol, side, type, price, qty_initial,
	                     qty_remaining, status)
	  VALUES (?,?,?,?,?,?,?)`,
		o.Symbol, o.Side, o.Type, o.Price,
		o.QtyInitial, o.QtyRemaining, o.Status)
	return err
}

// Old version — no tx
func (r *OrderRepo) UpdateRemaining(ctx context.Context, id, remaining int64, status string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE orders SET qty_remaining=?, status=? WHERE id=?`,
		remaining, status, id)
	return err
}

// ✅ NEW version — accepts tx
func (r *OrderRepo) UpdateRemainingTx(ctx context.Context, tx *sql.Tx, id, remaining int64, status string) error {
	_, err := tx.ExecContext(ctx,
		`UPDATE orders SET qty_remaining=?, status=? WHERE id=?`,
		remaining, status, id)
	return err
}


func (r *OrderRepo) InsertTx(ctx context.Context, tx *sql.Tx,
    o *models.Order) (int64, error) {

    res, err := tx.ExecContext(ctx, `
        INSERT INTO orders(symbol, side, type, price, qty_initial,
                           qty_remaining, status)
        VALUES (?,?,?,?,?,?,?)`,
        o.Symbol, o.Side, o.Type, o.Price,
        o.QtyInitial, o.QtyRemaining, o.Status)
    if err != nil {
        return 0, err
    }
    id, err := res.LastInsertId()
    return id, err
}


func (r *OrderRepo) DB() *sql.DB {
	return r.db
}

/* ---------- cancel ONE open order -------------------------------------- */

func (r *OrderRepo) CancelOrder(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE orders SET status='cancelled'
		 WHERE id=? AND status='open'`,
		id)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

/* ---------- order-book snapshot ---------------------------------------- */

type BookLevel struct {
	Price float64 `json:"price"`
	Qty   int64   `json:"qty"`
}

func (r *OrderRepo) GetOrderBook(ctx context.Context, symbol string) (bids, asks []BookLevel, err error) {
	// bids (buy side) – highest price first
	rows, err := r.db.QueryContext(ctx, `
		SELECT price, SUM(qty_remaining)
		FROM orders
		WHERE symbol=? AND side='buy' AND status='open'
		GROUP BY price ORDER BY price DESC`, symbol)
	if err != nil {
		return
	}
	for rows.Next() {
		var lvl BookLevel
		if err = rows.Scan(&lvl.Price, &lvl.Qty); err != nil {
			return
		}
		bids = append(bids, lvl)
	}
	rows.Close()

	// asks (sell side) – lowest price first
	rows, err = r.db.QueryContext(ctx, `
		SELECT price, SUM(qty_remaining)
		FROM orders
		WHERE symbol=? AND side='sell' AND status='open'
		GROUP BY price ORDER BY price ASC`, symbol)
	if err != nil {
		return
	}
	for rows.Next() {
		var lvl BookLevel
		if err = rows.Scan(&lvl.Price, &lvl.Qty); err != nil {
			return
		}
		asks = append(asks, lvl)
	}
	return
}
