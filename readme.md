# 🏦 Golang Order-Matching System

Minimal matching engine (limit & market orders), raw SQL persistence, and a simple REST API.

---

## 📁 1 Clone & bootstrap

```bash
# 1. grab the code
git clone https://github.com/yourname/golang-order-matching-system.git
cd golang-order-matching-system

# 2. tidy Go deps (creates go.sum)
go mod tidy
```

🐳 2 One-command stack with Docker Compose
Requires Docker 20+ (or Docker Desktop).

2-a Create docker-compose.yml (already in repo)

```bash
version: "3.8"
services:
  db:
    image: mysql:8
    container_name: oms-mysql
    restart: always
    environment:
      MYSQL_ROOT_PASSWORD: pass
      MYSQL_DATABASE: oms
    ports: [ "3306:3306" ]
    volumes:
      - db_data:/var/lib/mysql
      - ./db/migrate.sql:/docker-entrypoint-initdb.d/01-schema.sql

  api:
    build: .
    container_name: oms-api
    depends_on: [ db ]
    environment:
      MYSQL_DSN: root:pass@tcp(db:3306)/oms?parseTime=true
    ports: [ "8080:8080" ]

volumes: { db_data: {} }
```

Launch everything

```bash
docker compose up --build
```

MySQL 8 starts, seeds schema from db/migrate.sql.

Go binary builds once, then the API prints:
⇨ listening on :8080

Importable Postman collection

import oms_api.postman_collection.json ➜ Postman “Import”
