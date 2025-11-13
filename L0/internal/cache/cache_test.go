package cache

import (
	"testing"

	"L0/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCache_SetAndGet(t *testing.T) {
	cache := NewCache()
	order := &models.Order{
		OrderUID:    "testUID",
		TrackNumber: "TrackNumber",
		Entry:       "WB",
	}

	err := cache.Set(order)
	require.NoError(t, err)

	retrieved, exists := cache.Get("testUID")
	require.True(t, exists)
	assert.Equal(t, order.OrderUID, retrieved.OrderUID)
	assert.Equal(t, order.TrackNumber, retrieved.TrackNumber)
}

func TestCache_GetNonExistent(t *testing.T) {
	cache := NewCache()

	order, exists := cache.Get("non-existent")
	assert.False(t, exists)
	assert.Nil(t, order)
}

func TestCache_GetAll(t *testing.T) {
	cache := NewCache()

	order1 := &models.Order{OrderUID: "order-1"}
	order2 := &models.Order{OrderUID: "order-2"}

	_ = cache.Set(order1)
	_ = cache.Set(order2)

	all := cache.GetAll()
	assert.Len(t, all, 2)
}

func TestCache_Size(t *testing.T) {
	cache := NewCache()

	assert.Equal(t, 0, cache.Size())

	_ = cache.Set(&models.Order{OrderUID: "order-1"})
	assert.Equal(t, 1, cache.Size())

	_ = cache.Set(&models.Order{OrderUID: "order-2"})
	assert.Equal(t, 2, cache.Size())
}

func TestCache_SetInvalidOrder(t *testing.T) {
	cache := NewCache()

	err := cache.Set(nil)
	assert.Error(t, err)

	err = cache.Set(&models.Order{OrderUID: ""})
	assert.Error(t, err)
}

func TestCache_Cleanup(t *testing.T) {
	cache := NewCache()

	order := &models.Order{OrderUID: "test-order"}
	_ = cache.Set(order)

	_, exists := cache.Get("test-order")
	assert.True(t, exists)

	cache.Cleanup()

	_, exists = cache.Get("test-order")
	assert.True(t, exists)
}
