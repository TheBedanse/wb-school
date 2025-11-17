package cache

import (
	"L0/internal/interfaces"
	"L0/internal/models"
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

var _ interfaces.Cache = (*Cache)(nil)

type Cache struct {
	mu         sync.RWMutex
	orders     map[string]*cacheEntry
	maxSize    int
	dafaultTTL time.Duration
}
type cacheEntry struct {
	order      *models.Order
	expiresAt  time.Time
	lastAccess time.Time
}

func NewCache() interfaces.Cache {
	return &Cache{
		orders:     make(map[string]*cacheEntry),
		maxSize:    1000,
		dafaultTTL: 10 * time.Minute,
	}
}

func (c *Cache) Set(order *models.Order) error {
	if order == nil || order.OrderUID == "" {
		return fmt.Errorf("invalid order")
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	c.orders[order.OrderUID] = &cacheEntry{
		order:      order,
		expiresAt:  time.Now().Add(c.dafaultTTL),
		lastAccess: time.Now(),
	}
	return nil
}

func (c *Cache) Get(orderUID string) (*models.Order, bool) {
	if orderUID == "" {
		return nil, false
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exists := c.orders[orderUID]
	if !exists {
		return nil, false
	}

	if time.Now().After(entry.expiresAt) {
		delete(c.orders, orderUID)
		return nil, false
	}

	entry.lastAccess = time.Now()
	return entry.order, true
}

func (c *Cache) GetAll() []*models.Order {
	c.mu.RLock()
	defer c.mu.RUnlock()

	orders := make([]*models.Order, 0, len(c.orders))
	now := time.Now()

	for _, entry := range c.orders {
		if now.After(entry.expiresAt) {
			continue
		}
		orders = append(orders, entry.order)
	}

	return orders
}

func (c *Cache) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()

	for orderUID, entry := range c.orders {
		if now.After(entry.expiresAt) {
			delete(c.orders, orderUID)
		}
	}

}

func (c *Cache) StartCleanupWorker(ctx context.Context) {
	ticker := time.NewTicker(3 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.Cleanup()
		case <-ctx.Done():
			log.Println("Cleanup worker stopped")
			return
		}
	}
}

func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.orders)
}
