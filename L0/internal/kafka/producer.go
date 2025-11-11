package kafka

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"L0/internal/models"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers string) *Producer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(strings.Split(brokers, ",")...),
		Topic:    "orders",
		Balancer: &kafka.LeastBytes{},
	}

	return &Producer{writer: writer}
}

func (p *Producer) SendOrder(ctx context.Context, order *models.Order) error {
	orderJSON, err := json.Marshal(order)
	if err != nil {
		return err
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(order.OrderUID),
		Value: orderJSON,
	})
	if err != nil {
		return err
	}

	log.Printf("Order sent to Kafka: %s", order.OrderUID)
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

// Функция для генерации и отправки заказа
func GenerateAndSendOrder(producer *Producer) error {
	order := GenerateTestOrder()
	ctx := context.Background()
	return producer.SendOrder(ctx, order)
}
