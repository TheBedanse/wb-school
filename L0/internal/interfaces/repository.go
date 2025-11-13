package interfaces

import (
	"context"

	"L0/internal/models"
)

//go:generate mockgen -source=repository.go -destination=../mocks/mock_repository.go -package=mocks

type Repository interface {
	SaveOrder(ctx context.Context, order *models.Order) error
	GetOrderByUID(ctx context.Context, orderUID string) (*models.Order, error)
	GetAllOrderUIDs(ctx context.Context) ([]string, error)
	Close()
}
