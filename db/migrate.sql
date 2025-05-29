-- schema and seed data for order-matching system ----------

CREATE DATABASE IF NOT EXISTS oms;
USE oms;

CREATE TABLE IF NOT EXISTS orders (
  id             BIGINT AUTO_INCREMENT PRIMARY KEY,
  symbol         VARCHAR(16)  NOT NULL,
  side           ENUM('buy','sell') NOT NULL,
  type           ENUM('limit','market') NOT NULL,
  price          DECIMAL(18,4),          -- NULL for market
  qty_initial    BIGINT       NOT NULL,
  qty_remaining  BIGINT       NOT NULL,
  status         ENUM('open','filled','cancelled') NOT NULL,
  ts             TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_symbol_side_price (symbol, side, price),
  INDEX idx_status (status)
);

CREATE TABLE IF NOT EXISTS trades (
  id            BIGINT AUTO_INCREMENT PRIMARY KEY,
  symbol        VARCHAR(16) NOT NULL,
  buy_order_id  BIGINT NOT NULL,
  sell_order_id BIGINT NOT NULL,
  price         DECIMAL(18,4) NOT NULL,
  qty           BIGINT NOT NULL,
  ts            TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (buy_order_id)  REFERENCES orders(id),
  FOREIGN KEY (sell_order_id) REFERENCES orders(id)
);

-- seed a symbol so GET /orderbook isn’t empty
INSERT IGNORE INTO orders(symbol, side, type, price, qty_initial,
                          qty_remaining, status)
VALUES ('ACME', 'buy', 'limit', 10.50, 100, 100, 'open');
