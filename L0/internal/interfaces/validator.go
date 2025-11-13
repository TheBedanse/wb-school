package interfaces

import "L0/internal/models"

//go:generate mockgen -source=validator.go -destination=../mocks/mock_validator.go -package=mocks
type Validator interface {
	ValidateOrder(order *models.Order) error
}
