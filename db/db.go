package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" // mysql driver
)

// Connect tries to connect to MySQL with retries.
func Connect(dsn string) (*sql.DB, error) {
	const (
		maxRetries = 10
		retryDelay = 2 * time.Second
	)

	var db *sql.DB
	var err error

	for i := 1; i <= maxRetries; i++ {
		db, err = sql.Open("mysql", dsn)
		if err == nil {
			db.SetMaxOpenConns(10)
			db.SetMaxIdleConns(5)
			db.SetConnMaxIdleTime(0)

			err = db.Ping()
			if err == nil {
				fmt.Println("✅ Connected to MySQL successfully.")
				return db, nil
			}
		}

		fmt.Printf("⚠️ Attempt %d: Failed to connect to MySQL: %v\n", i, err)
		time.Sleep(retryDelay)
	}

	return nil, fmt.Errorf("❌ could not connect to MySQL after %d attempts: %v", maxRetries, err)
}
