package services

import (
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type OrderService struct {
	db    *gorm.DB
	cache *redis.Client
	queue string
}

func (o *OrderService) Publish(msg string) int64 {
	return o.cache.Publish(o.queue, msg).Val()
}

func NewOrderService(d *gorm.DB, c *redis.Client, q string) *OrderService {
	return &OrderService{db: d, cache: c, queue: q}
}
