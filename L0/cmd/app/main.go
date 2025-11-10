package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"L0/internal/cache"
	"L0/internal/config"
	"L0/internal/database"
	"L0/internal/handler"
	"L0/internal/kafka"
	"L0/internal/service"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	cfg := config.LoadConfig()

	// if err := autoMigrate(cfg.DBPassword); err != nil {
	// 	log.Fatal("Failed to apply migrations:", err)
	// }

	db, err := database.NewDB(cfg.DBPassword)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	orderCache := cache.NewCache()
	orderService := service.NewOrderService(db, orderCache)

	ctx := context.Background()
	if err := orderService.RestoreCacheFromDB(ctx); err != nil {
		log.Printf("Warning: failed to restore cache from DB: %v", err)
	} else {
		log.Printf("Cache restored successfully. Loaded %d orders", orderCache.Size())
	}

	consumer := kafka.NewMockConsumer(orderService)
	defer consumer.Close()

	go consumer.Start(ctx)

	orderHandler, err := handler.NewOrderHandler(orderService)
	if err != nil {
		log.Fatal("Error creating order handler:", err)
	}

	http.HandleFunc("/", orderHandler.ShowHomePage)
	http.HandleFunc("/order/", orderHandler.ShowOrder)
	http.HandleFunc("/api/order/", orderHandler.GetOrderJSON)

	server := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("Server starting on :%s", cfg.HTTPPort)
		log.Printf("Access the application at: http://localhost:%s", cfg.HTTPPort)
		log.Printf("Mock Kafka Consumer is generating test orders every 15 seconds")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server error:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited properly")
}

func autoMigrate(dbPassword string) error {
	connStr := fmt.Sprintf("postgres://L0User:%s@localhost:5432/L0?sslmode=disable", dbPassword)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer db.Close()

	return goose.Up(db, "schema")
}
