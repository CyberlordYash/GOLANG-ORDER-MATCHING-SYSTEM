package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DSN string // MySQL DSN: user:pass@tcp(host:port)/dbname?parseTime=true
}

// Load reads .env (if present) then env vars.
// Falls back to sensible defaults so it’s hassle-free in dev.
func Load() (*Config, error) {
	_ = godotenv.Load() // ignore missing file

	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		user := os.Getenv("MYSQL_USER")
		if user == "" {
			user = "root"
		}
		pass := os.Getenv("MYSQL_PASSWORD")
		if pass == "" {
			pass = "pass"
		}
		host := os.Getenv("MYSQL_HOST")
		if host == "" {
			host = "127.0.0.1"
		}
		port := os.Getenv("MYSQL_PORT")
		if port == "" {
			port = "3306"
		}
		name := os.Getenv("MYSQL_DB")
		if name == "" {
			name = "oms"
		}
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, name)
	}
	return &Config{DSN: dsn}, nil
}
