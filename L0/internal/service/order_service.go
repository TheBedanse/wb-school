package service

import (
	"context"
	"fmt"
	"log"

	"L0/internal/interfaces"
	"L0/internal/models"
)

var _ interfaces.OrderService = (*OrderService)(nil)

type OrderService struct {
	orderRepo interfaces.Repository
	cache     interfaces.Cache
	validator interfaces.Validator
}

func NewOrderService(orderRepo interfaces.Repository, cache interfaces.Cache) interfaces.OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		cache:     cache,
		validator: &models.Validator{},
	}
}

func (s *OrderService) RestoreCacheFromDB(ctx context.Context) error {
	log.Println("Restoring cache from database")

	orderUIDs, err := s.orderRepo.GetAllOrderUIDs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get orderUID for cache restoration: %w", err)
	}

	count := 0
	for _, orderUID := range orderUIDs {
		order, err := s.orderRepo.GetOrderByUID(ctx, orderUID)
		if err != nil {
			log.Printf("Error restoring order %s: %v", orderUID, err)
			continue
		}

		if err := s.cache.Set(order); err != nil {
			log.Printf("Failed to cache order %s: %v", orderUID, err)
			continue
		}
		count++
	}

	log.Printf("Cache restore. Loaded %d orders", count)
	return nil
}

func (s *OrderService) ProcessOrder(ctx context.Context, order *models.Order) error {
	if err := s.validator.ValidateOrder(order); err != nil {
		return fmt.Errorf("order validation failed: %w", err)
	}

	if err := s.cache.Set(order); err != nil {
		return fmt.Errorf("failed to cache order: %w", err)
	}
	log.Printf("Order cached: %s", order.OrderUID)

	if err := s.orderRepo.SaveOrder(ctx, order); err != nil {
		return fmt.Errorf("failed to save order to DB: %w", err)
	}

	log.Printf("Order processed successfully(Service): %s", order.OrderUID)
	return nil
}

func (s *OrderService) GetOrder(ctx context.Context, orderUID string) (*models.Order, error) {
	if orderUID == "" {
		return nil, fmt.Errorf("orderUID cannot be empty")
	}

	if order, exists := s.cache.Get(orderUID); exists {
		return order, nil
	}

	order, err := s.orderRepo.GetOrderByUID(ctx, orderUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order from DB: %w", err)
	}

	if err := s.cache.Set(order); err != nil {
		log.Printf("Warning: failed to cache order %s: %v", orderUID, err)
	}

	return order, nil
}

func (s *OrderService) GetAllOrders() []*models.Order {
	return s.cache.GetAll()
}
