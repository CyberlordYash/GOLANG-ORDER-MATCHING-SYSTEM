
# Golang Order‑Matching System 🏦

A minimal limit/market order‑matching engine written in Go.  
* Raw **SQL** for persistence (no ORM)  
* In‑memory matching core (heap‑based order book)  
* REST + JSON transport via **Gin**

---

## ✨ High‑Level Logic

1. **POST /orders**  
   *Insert a new order, get a real DB id, match it in memory, persist any fills—all in one DB transaction.*

2. **Engine**  
   * Two price‐level heaps (bids max‑heap, asks min‑heap) + FIFO queues → O(log N) insert/remove  
   * Price‑time priority: highest bid / lowest ask first, then earliest timestamp.

3. **Persistence**  
   * `orders` table stores live & historical orders (`status=open|filled|cancelled`).  
   * `trades` table stores executions with FK back to both orders.

---

## 🐳 Quick Start

///bash
git clone https://github.com/yourname/golang-order-matching-system.git
cd golang-order-matching-system
go mod tidy                       # creates go.sum

docker compose up --build         # MySQL 8 + API on :8080
///bash

The API logs `⇨ listening on :8080` when ready.

---

## 🔗 Endpoint Reference (current state)

| Method | Path | Payload / Query | Response | Status |
|--------|------|-----------------|----------|--------|
| `POST` | `/orders` | `{ "symbol":"ACME", "side":"buy", "type":"limit", "price":10.50, "quantity":100 }` | `200 OK` → `{ order_id, executions[] }` | **Implemented** |
| `DELETE` | `/orders/:id` | – | `501 Not Implemented` | Stub |
| `GET` | `/orderbook?symbol=ACME` | query `symbol` | `501 Not Implemented` | Stub |
| `GET` | `/trades?symbol=ACME&limit=100` | optional `symbol`, `limit` | `501 Not Implemented` | Stub |

### Example success (POST /orders)

```jsonc
{
  "order_id": 1717000000123456000,
  "executions": [
    { "taker_id": 1717000000123456000, "maker_id": 1, "price": 10.5, "qty": 70 }
  ]
}
```

### Example error

```jsonc
{ "error": "price is required for limit orders" }
```

---

## 📂 Folder Layout

```
.
├── main.go
├── api/             # HTTP handlers & DTOs
├── engine/          # Matching core (pure Go)
├── repo/            # Raw SQL data‑access
├── models/          # Plain structs
├── db/              # migrate.sql + db helpers
├── config/          # env → DSN
└── Dockerfile
```

---

## 🚀 How `docker compose up --build` works

1. **db** service → MySQL 8, seeds schema via mounted `db/migrate.sql`.
2. **api** service → multistage Dockerfile builds static Go binary, starts API.
3. Containers come up; visit `http://localhost:8080`.

Re‑run without code changes: `docker compose up` (skips rebuild).  
Rebuild API only: `docker compose build api && docker compose up`.

---

## 🛠️ Next Steps

* Implement **cancel, orderbook, trades** handlers.  
* Add unit tests (`go test ./engine/...`).  
* Add CI (GitHub Actions) & healthchecks.

Happy matching 🚀
