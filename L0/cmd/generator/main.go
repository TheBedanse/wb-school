package main

import (
	"context"
	"log"
	"time"

	"L0/internal/config"
	"L0/internal/kafka"
)

func main() {
	cfg := config.LoadConfig()
	ctx := context.Background()

	producer := kafka.NewProducer(cfg.KafkaBroker)
	defer producer.Close()

	log.Println("Generating test orders every 15 seconds")

	if err := generateAndSendOrder(producer); err != nil {
		log.Printf("Error sending initial order: %v", err)
	}

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping Mock Kafka Consumer")
			return
		case <-ticker.C:
			if err := generateAndSendOrder(producer); err != nil {
				log.Printf("Error sending order: %v", err)
			}
		}
	}
}

func generateAndSendOrder(producer *kafka.Producer) error {
	order := kafka.GenerateTestOrder()
	ctx := context.Background()
	return producer.SendOrder(ctx, order)
}
