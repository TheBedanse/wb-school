package service

import (
	"context"
	"errors"
	"testing"

	"L0/internal/mocks"
	"L0/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestOrderService_ProcessOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockCache := mocks.NewMockCache(ctrl)
	mockValidator := mocks.NewMockValidator(ctrl)

	service := &OrderService{
		orderRepo: mockRepo,
		cache:     mockCache,
		validator: mockValidator,
	}

	ctx := context.Background()
	order := &models.Order{OrderUID: "test-123"}

	t.Run("successful processing", func(t *testing.T) {
		mockValidator.EXPECT().ValidateOrder(order).Return(nil)
		mockCache.EXPECT().Set(order).Return(nil)
		mockRepo.EXPECT().SaveOrder(ctx, order).Return(nil)

		err := service.ProcessOrder(ctx, order)

		require.NoError(t, err)
	})

	t.Run("validation failed", func(t *testing.T) {
		mockValidator.EXPECT().ValidateOrder(order).Return(errors.New("validation error"))

		err := service.ProcessOrder(ctx, order)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
	})

	t.Run("cache set failed", func(t *testing.T) {
		mockValidator.EXPECT().ValidateOrder(order).Return(nil)
		mockCache.EXPECT().Set(order).Return(errors.New("cache error"))

		err := service.ProcessOrder(ctx, order)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to cache order")
	})

	t.Run("repository save failed", func(t *testing.T) {
		mockValidator.EXPECT().ValidateOrder(order).Return(nil)
		mockCache.EXPECT().Set(order).Return(nil)
		mockRepo.EXPECT().SaveOrder(ctx, order).Return(errors.New("db error"))

		err := service.ProcessOrder(ctx, order)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to save order to DB")
	})
}

func TestOrderService_GetOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockCache := mocks.NewMockCache(ctrl)

	service := &OrderService{
		orderRepo: mockRepo,
		cache:     mockCache,
		validator: &models.Validator{},
	}

	ctx := context.Background()
	order := &models.Order{OrderUID: "test-123"}

	t.Run("from cache", func(t *testing.T) {
		mockCache.EXPECT().Get("test-123").Return(order, true)

		result, err := service.GetOrder(ctx, "test-123")

		require.NoError(t, err)
		assert.Equal(t, order, result)
	})

	t.Run("from repository", func(t *testing.T) {
		mockCache.EXPECT().Get("test-123").Return((*models.Order)(nil), false)
		mockRepo.EXPECT().GetOrderByUID(ctx, "test-123").Return(order, nil)
		mockCache.EXPECT().Set(order).Return(nil)

		result, err := service.GetOrder(ctx, "test-123")

		require.NoError(t, err)
		assert.Equal(t, order, result)
	})

	t.Run("empty orderUID", func(t *testing.T) {
		result, err := service.GetOrder(ctx, "")

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "orderUID cannot be empty")
	})

	t.Run("repository error", func(t *testing.T) {
		mockCache.EXPECT().Get("test-123").Return((*models.Order)(nil), false)
		mockRepo.EXPECT().GetOrderByUID(ctx, "test-123").Return((*models.Order)(nil), errors.New("db error"))

		result, err := service.GetOrder(ctx, "test-123")

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to get order from DB")
	})
}

func TestOrderService_GetAllOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockCache := mocks.NewMockCache(ctrl)

	service := &OrderService{
		orderRepo: mockRepo,
		cache:     mockCache,
		validator: &models.Validator{},
	}

	orders := []*models.Order{
		{OrderUID: "order-1"},
		{OrderUID: "order-2"},
	}

	mockCache.EXPECT().GetAll().Return(orders)

	result := service.GetAllOrders()

	assert.Equal(t, orders, result)
}

func TestOrderService_RestoreCacheFromDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockCache := mocks.NewMockCache(ctrl)

	service := &OrderService{
		orderRepo: mockRepo,
		cache:     mockCache,
		validator: &models.Validator{},
	}

	ctx := context.Background()

	t.Run("successful restore", func(t *testing.T) {
		orderUIDs := []string{"order-1", "order-2"}
		order1 := &models.Order{OrderUID: "order-1"}
		order2 := &models.Order{OrderUID: "order-2"}

		mockRepo.EXPECT().GetAllOrderUIDs(ctx).Return(orderUIDs, nil)
		mockRepo.EXPECT().GetOrderByUID(ctx, "order-1").Return(order1, nil)
		mockRepo.EXPECT().GetOrderByUID(ctx, "order-2").Return(order2, nil)
		mockCache.EXPECT().Set(order1).Return(nil)
		mockCache.EXPECT().Set(order2).Return(nil)

		err := service.RestoreCacheFromDB(ctx)

		require.NoError(t, err)
	})

	t.Run("error getting order uids", func(t *testing.T) {
		mockRepo.EXPECT().GetAllOrderUIDs(ctx).Return(nil, errors.New("db error"))

		err := service.RestoreCacheFromDB(ctx)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get orderUID for cache restoration")
	})

	t.Run("error getting order", func(t *testing.T) {
		orderUIDs := []string{"order-1", "order-2"}
		order2 := &models.Order{OrderUID: "order-2"}

		mockRepo.EXPECT().GetAllOrderUIDs(ctx).Return(orderUIDs, nil)
		mockRepo.EXPECT().GetOrderByUID(ctx, "order-1").Return((*models.Order)(nil), errors.New("db error"))
		mockRepo.EXPECT().GetOrderByUID(ctx, "order-2").Return(order2, nil)
		mockCache.EXPECT().Set(order2).Return(nil)

		err := service.RestoreCacheFromDB(ctx)

		require.NoError(t, err)
	})

	t.Run("error caching order", func(t *testing.T) {
		orderUIDs := []string{"order-1", "order-2"}
		order1 := &models.Order{OrderUID: "order-1"}
		order2 := &models.Order{OrderUID: "order-2"}

		mockRepo.EXPECT().GetAllOrderUIDs(ctx).Return(orderUIDs, nil)
		mockRepo.EXPECT().GetOrderByUID(ctx, "order-1").Return(order1, nil)
		mockRepo.EXPECT().GetOrderByUID(ctx, "order-2").Return(order2, nil)
		mockCache.EXPECT().Set(order1).Return(errors.New("cache error"))
		mockCache.EXPECT().Set(order2).Return(nil)

		err := service.RestoreCacheFromDB(ctx)

		require.NoError(t, err)
	})
}
