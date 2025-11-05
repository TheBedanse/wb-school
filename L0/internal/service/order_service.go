package service

import (
	"context"
	"log"

	"L0/internal/cache"
	"L0/internal/database"
	"L0/internal/models"
)

type OrderService struct {
	orderRepo *database.Database
	cache     *cache.Cache
}

func NewOrderService(orderRepo *database.Database, cache *cache.Cache) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		cache:     cache,
	}
}

func (s *OrderService) RestoreCacheFromDB(ctx context.Context) error {
	log.Println("Restoring cache from database...")

	orderUIDs, err := s.orderRepo.GetAllOrderUIDs(ctx)
	if err != nil {
		log.Printf("Warning: Could not restore cache from DB: %v", err)
		return nil
	}

	count := 0
	for _, orderUID := range orderUIDs {
		order, err := s.orderRepo.GetOrderByUID(ctx, orderUID)
		if err != nil {
			log.Printf("Error restoring order %s: %v", orderUID, err)
			continue
		}
		s.cache.Set(order)
		count++
	}

	log.Printf("Cache restored successfully. Loaded %d orders", count)
	return nil
}

func (s *OrderService) ProcessOrder(ctx context.Context, order *models.Order) error {
	s.cache.Set(order)
	log.Printf("Order cached: %s", order.OrderUID)

	if err := s.orderRepo.SaveOrder(ctx, order); err != nil {
		log.Printf("Warning: Failed to save order %s to DB: %v", order.OrderUID, err)
		return nil
	}

	log.Printf("Order saved to database(ProcessOrder): %s", order.OrderUID)
	return nil
}

func (s *OrderService) GetOrder(ctx context.Context, orderUID string) (*models.Order, error) {
	if order, exists := s.cache.Get(orderUID); exists {
		return order, nil
	}

	order, err := s.orderRepo.GetOrderByUID(ctx, orderUID)
	if err != nil {
		return nil, err
	}

	s.cache.Set(order)
	return order, nil
}

func (s *OrderService) GetAllOrders() []*models.Order {
	return s.cache.GetAll()
}
