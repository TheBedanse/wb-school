package cache

import (
	"L0/internal/models"
	"sync"
)

type Cache struct {
	mu     sync.RWMutex
	orders map[string]*models.Order
}

func NewCache() *Cache {
	return &Cache{
		orders: make(map[string]*models.Order),
	}
}

func (c *Cache) Set(order *models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.orders[order.OrderUID] = order
}

func (c *Cache) Get(orderUID string) (*models.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, exists := c.orders[orderUID]
	return order, exists
}

func (c *Cache) GetAll() []*models.Order {
	c.mu.RLock()
	defer c.mu.RUnlock()

	orders := make([]*models.Order, 0, len(c.orders))
	for _, order := range c.orders {
		orders = append(orders, order)
	}
	return orders
}

func (c *Cache) Delete(orderUID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.orders, orderUID)
}

func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.orders)
}
