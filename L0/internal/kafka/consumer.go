package kafka

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"L0/internal/models"
	"L0/internal/service"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader       *kafka.Reader
	orderService *service.OrderService
}

func NewConsumer(orderService *service.OrderService, brokers string) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  strings.Split(brokers, ","),
		Topic:    "orders",
		GroupID:  "order-service",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	return &Consumer{
		reader:       reader,
		orderService: orderService,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	log.Printf("Starting Kafka Consumer for topic: orders")

	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Kafka read error: %v", err)
			continue
		}

		log.Printf("Received message: partition=%d, offset=%d", msg.Partition, msg.Offset)

		var order models.Order
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Printf("Failed to parse order: %v", err)
			continue
		}

		if err := c.orderService.ProcessOrder(ctx, &order); err != nil {
			log.Printf("Failed to process order %s: %v", order.OrderUID, err)
		} else {
			log.Printf("Order processed successfully: %s", order.OrderUID)
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
