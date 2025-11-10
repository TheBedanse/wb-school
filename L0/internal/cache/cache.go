package cache

import (
	"L0/internal/models"
	"fmt"
	"sync"
	"time"
)

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

func NewCache() *Cache {
	return &Cache{
		orders:     make(map[string]*cacheEntry),
		maxSize:    1000,
		dafaultTTL: 10 * time.Minute,
	}
}

func (c *Cache) Set(order *models.Order) error {
	if order == nil || order.OrderUID == "" {
		return ErrInvalidOrder
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

func (c *Cache) Delete(orderUID string) error {
	if orderUID == "" {
		return ErrInvalidOrderUID
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.orders, orderUID)
	return nil
}

func (c *Cache) Cleanup() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	count := 0

	for orderUID, entry := range c.orders {
		if now.After(entry.expiresAt) {
			delete(c.orders, orderUID)
			count++
		}
	}

	return count
}

func (c *Cache) SetWithTTL(order *models.Order, ttl time.Duration) error {
	if order == nil || order.OrderUID == "" {
		return ErrInvalidOrder
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.orders) >= c.maxSize {
		c.evictOldest()
	}

	c.orders[order.OrderUID] = &cacheEntry{
		order:      order,
		expiresAt:  time.Now().Add(ttl),
		lastAccess: time.Now(),
	}
	return nil
}

func (c *Cache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time

	for key, entry := range c.orders {
		if oldestKey == "" || entry.lastAccess.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.lastAccess
		}
	}

	if oldestKey != "" {
		delete(c.orders, oldestKey)
	}
}

func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.orders)
}

func (c *Cache) GetStats() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	now := time.Now()
	expiredCount := 0

	for _, entry := range c.orders {
		if now.After(entry.expiresAt) {
			expiredCount++
		}
	}

	return map[string]interface{}{
		"total_entries":   len(c.orders),
		"expired_entries": expiredCount,
		"max_size":        c.maxSize,
	}
}

var (
	ErrInvalidOrder    = fmt.Errorf("invalid order")
	ErrInvalidOrderUID = fmt.Errorf("invalid orderUID")
	ErrCacheFull       = fmt.Errorf("cache is full")
)
