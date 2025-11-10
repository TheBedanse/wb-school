package kafka

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"L0/internal/models"
	"L0/internal/service"
)

type MockConsumer struct {
	orderService *service.OrderService
}

func NewMockConsumer(orderService *service.OrderService) *MockConsumer {
	return &MockConsumer{
		orderService: orderService,
	}
}

func (m *MockConsumer) Start(ctx context.Context) {
	log.Println("Starting Mock Kafka Consumer - generating test orders every 15 seconds")

	go m.startCacheCleanup(ctx)

	if err := m.generateAndProcessOrder(ctx); err != nil {
		log.Printf("Error processing initial mock order: %v", err)
	}

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping Mock Kafka Consumer")
			return
		case <-ticker.C:
			if err := m.generateAndProcessOrder(ctx); err != nil {
				log.Printf("Error processing mock order: %v", err)
			}
		}
	}
}

func (m *MockConsumer) generateAndProcessOrder(ctx context.Context) error {
	order := GenerateTestOrder()

	if err := m.orderService.ProcessOrder(ctx, order); err != nil {
		return fmt.Errorf("failed to process order %s: %w", order.OrderUID, err)
	}

	log.Printf("Mock order processed successfully: %s", order.OrderUID)
	return nil
}

func (m *MockConsumer) startCacheCleanup(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			cleaned := m.orderService.CleanupCache()
			if cleaned > 0 {
				log.Printf("Cache cleanup: removed %d expired entries", cleaned)
			}

			stats := m.orderService.GetCacheStats()
			log.Printf("Cache stats: total=%d, expired=%d, max_size=%d",
				stats["total_entries"], stats["expired_entries"], stats["max_size"])
		}
	}
}

func (m *MockConsumer) Close() error {
	log.Println("Mock Kafka Consumer closed")
	return nil
}

func GenerateTestOrder() *models.Order {
	orderUID := generateOrderUID()
	trackNumber := generateTrackNumber()
	rid := generateRid()

	return &models.Order{
		OrderUID:    orderUID,
		TrackNumber: trackNumber,
		Entry:       "WBIL",
		Delivery: models.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: models.Payment{
			Transaction:  orderUID,
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       rand.Intn(5000) + 1000,
			PaymentDt:    time.Now().Unix(),
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []models.Item{
			{
				ChrtID:      int64(rand.Intn(100000)),
				TrackNumber: trackNumber,
				Price:       453,
				Rid:         rid,
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmID:        int64(rand.Intn(100000)),
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmID:              99,
		DateCreated:       time.Now(),
		OofShard:          "1",
	}
}

func generateOrderUID() string {
	return fmt.Sprintf("order-%d-%d", time.Now().Unix(), rand.Intn(1000))
}

func generateTrackNumber() string {
	return fmt.Sprintf("WBILMTESTTRACK-%d", rand.Intn(10000))
}

func generateRid() string {
	return fmt.Sprintf("ab4219087a764ae0btest-%d", rand.Intn(1000))
}
