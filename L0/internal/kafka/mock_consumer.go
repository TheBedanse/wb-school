package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"L0/internal/models"
	"L0/internal/service"

	"github.com/brianvoe/gofakeit/v7"
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
	orderUID := gofakeit.UUID()
	trackNumber := generateTrackNumber()

	return &models.Order{
		OrderUID:    orderUID,
		TrackNumber: trackNumber,
		Entry:       "WBIL",
		Delivery: models.Delivery{
			Name:    gofakeit.Name(),
			Phone:   gofakeit.Phone(),
			Zip:     gofakeit.Zip(),
			City:    gofakeit.City(),
			Address: gofakeit.Street(),
			Region:  gofakeit.State(),
			Email:   gofakeit.Email(),
		},
		Payment: models.Payment{
			Transaction:  orderUID,
			RequestID:    "",
			Currency:     gofakeit.CurrencyShort(),
			Provider:     gofakeit.Company(),
			Amount:       gofakeit.Number(1000, 10000),
			PaymentDt:    gofakeit.Date().Unix(),
			Bank:         gofakeit.BankName(),
			DeliveryCost: gofakeit.Number(100, 1000),
			GoodsTotal:   gofakeit.Number(10, 1000),
			CustomFee:    gofakeit.Number(0, 100),
		},
		Items:             generateFakeItems(gofakeit.Number(1, 5), trackNumber),
		Locale:            gofakeit.LanguageAbbreviation(),
		InternalSignature: "",
		CustomerID:        gofakeit.UUID(),
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmID:              gofakeit.Number(1, 100),
		DateCreated:       gofakeit.Date(),
		OofShard:          "1",
	}
}

func generateTrackNumber() string {
	return "WB" + gofakeit.DigitN(12)
}

func generateFakeItems(count int, trackNum string) []models.Item {
	items := make([]models.Item, count)

	for i := 0; i < count; i++ {
		items[i] = models.Item{
			ChrtID:      gofakeit.Int64(),
			TrackNumber: trackNum,
			Price:       gofakeit.Number(100, 5000),
			Rid:         gofakeit.UUID(),
			Name:        gofakeit.ProductName(),
			Sale:        gofakeit.Number(0, 50),
			Size:        gofakeit.RandomString([]string{"0", "S", "M", "L", "XL"}),
			TotalPrice:  gofakeit.Number(100, 5000),
			NmID:        gofakeit.Int64(),
			Brand:       gofakeit.Company(),
			Status:      gofakeit.Number(100, 400),
		}
	}

	return items
}
