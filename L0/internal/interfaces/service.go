package interfaces

import (
	"context"

	"L0/internal/models"
)

//go:generate mockgen -source=service.go -destination=../mocks/mock_service.go -package=mocks

type OrderService interface {
	ProcessOrder(ctx context.Context, order *models.Order) error
	GetOrder(ctx context.Context, orderUID string) (*models.Order, error)
	GetAllOrders() []*models.Order
	RestoreCacheFromDB(ctx context.Context) error
}
