# Go Order-Matching System 🏦

A minimal, raw-SQL limit-/market-order matching engine with REST API.

## Quick Start

```bash
# clone & enter
git clone https://github.com/yourname/golang-order-matching-system.git
cd golang-order-matching-system

# install deps
go mod tidy

# spin up MySQL 8
docker run --name oms-mysql -e MYSQL_ROOT_PASSWORD=pass \
  -p 3306:3306 -d mysql:8
mysql -h127.0.0.1 -uroot -ppass < db/migrate.sql

# run API
go run ./...
```
