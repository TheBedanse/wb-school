package interfaces

import "L0/internal/models"

//go:generate mockgen -source=cache.go -destination=../mocks/mock_cache.go -package=mocks

type Cache interface {
	Set(order *models.Order) error
	Get(orderUID string) (*models.Order, bool)
	GetAll() []*models.Order
	Size() int
	Cleanup()
}
