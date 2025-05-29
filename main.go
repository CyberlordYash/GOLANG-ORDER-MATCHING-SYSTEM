package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"golang-order-matching-system/api"
	"golang-order-matching-system/config"
	"golang-order-matching-system/db"
	"golang-order-matching-system/engine"
	"golang-order-matching-system/repo"
)

func main() {
	// load env → config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	// DB connection pool
	sqlDB, err := db.Connect(cfg.DSN)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}

	// in-memory matching engine
	eng := engine.New()

	// repositories
	orderRepo := repo.NewOrderRepo(sqlDB)
	tradeRepo := repo.NewTradeRepo(sqlDB)

	// HTTP handlers (dependency injection)
	h := &api.Handler{
		Eng: eng,
		OR:  orderRepo,
		TR:  tradeRepo,
	}

	// router
	r := gin.Default()
	r.POST("/orders", h.PlaceOrder)
	r.DELETE("/orders/:id", h.CancelOrder)
	r.GET("/orderbook", h.GetOrderBook)
	r.GET("/trades", h.ListTrades)

	log.Println("⇨ listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
