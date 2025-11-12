package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createValidOrder() *Order {
	return &Order{
		OrderUID:        "test-123",
		TrackNumber:     "TRACK-001",
		Entry:           "WBIL",
		Locale:          "en",
		CustomerID:      "customer-123",
		DeliveryService: "meest",
		SmID:            99,
		DateCreated:     time.Now(),
		Delivery: Delivery{
			Name:    "John Doe",
			Phone:   "+1234567890",
			Zip:     "12345",
			City:    "Moscow",
			Address: "Street 1",
			Region:  "Moscow",
			Email:   "test@example.com",
		},
		Payment: Payment{
			Transaction:  "trans-123",
			Currency:     "USD",
			Provider:     "provider",
			Amount:       1000,
			PaymentDt:    time.Now().Unix(),
			Bank:         "bank",
			DeliveryCost: 500,
			GoodsTotal:   500,
			CustomFee:    0,
		},
		Items: []Item{
			{
				ChrtID:      12345,
				TrackNumber: "TRACK-001",
				Price:       100,
				Rid:         "rid-123",
				Name:        "Test Item",
				Sale:        0,
				Size:        "M",
				TotalPrice:  100,
				NmID:        67890,
				Brand:       "Brand",
				Status:      200,
			},
		},
	}
}

func TestValidator_ValidateOrder(t *testing.T) {
	validator := &Validator{}

	t.Run("valid order", func(t *testing.T) {
		order := createValidOrder()
		err := validator.ValidateOrder(order)
		require.NoError(t, err)
	})

	t.Run("nil order", func(t *testing.T) {
		err := validator.ValidateOrder(nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "order is nil")
	})

	t.Run("empty order UID", func(t *testing.T) {
		order := createValidOrder()
		order.OrderUID = ""
		err := validator.ValidateOrder(order)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "order_uid is required")
	})

	t.Run("empty track number", func(t *testing.T) {
		order := createValidOrder()
		order.TrackNumber = ""
		err := validator.ValidateOrder(order)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "track_number is required")
	})

	t.Run("invalid delivery email", func(t *testing.T) {
		order := createValidOrder()
		order.Delivery.Email = "invalid-email"
		err := validator.ValidateOrder(order)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid email format")
	})

	t.Run("invalid phone", func(t *testing.T) {
		order := createValidOrder()
		order.Delivery.Phone = "invalid-phone"
		err := validator.ValidateOrder(order)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid phone format")
	})

	t.Run("negative amount", func(t *testing.T) {
		order := createValidOrder()
		order.Payment.Amount = -100
		err := validator.ValidateOrder(order)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "payment amount cannot be negative")
	})

	t.Run("empty items", func(t *testing.T) {
		order := createValidOrder()
		order.Items = []Item{}
		err := validator.ValidateOrder(order)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "at least one item is required")
	})

	t.Run("invalid item", func(t *testing.T) {
		order := createValidOrder()
		order.Items[0].ChrtID = 0
		err := validator.ValidateOrder(order)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "chrt_id must be positive")
	})
}

func TestValidator_IsValidEmail(t *testing.T) {
	validator := &Validator{}

	validEmails := []string{
		"test@example.com",
		"user.name@domain.co.uk",
		"user+tag@example.org",
	}

	invalidEmails := []string{
		"invalid-email",
		"@example.com",
		"user@",
		"user@.com",
	}

	for _, email := range validEmails {
		assert.True(t, validator.isValidEmail(email), "Email should be valid: %s", email)
	}

	for _, email := range invalidEmails {
		assert.False(t, validator.isValidEmail(email), "Email should be invalid: %s", email)
	}
}

func TestValidator_IsValidPhone(t *testing.T) {
	validator := &Validator{}

	validPhones := []string{
		"+1234567890",
		"+441234567890",
		"1234567890",
	}

	invalidPhones := []string{
		"invalid-phone",
		"+123",
		"abc123",
		"",
	}

	for _, phone := range validPhones {
		assert.True(t, validator.isValidPhone(phone), "Phone should be valid: %s", phone)
	}

	for _, phone := range invalidPhones {
		assert.False(t, validator.isValidPhone(phone), "Phone should be invalid: %s", phone)
	}
}
