package mem

import (
	"github.com/patrickmn/go-cache"
	"prisma/dto"
	"prisma/prisma-shema/db"
	"time"
)

func Cache(order dto.Order) *cache.Cache {
	c := cache.New(5*time.Minute, 10*time.Minute)
	c.Set("ключ_кеша", order, cache.DefaultExpiration)
	return c
}

func CacheItem(Item []db.ItemsModel) *cache.Cache {
	c := cache.New(5*time.Minute, 10*time.Minute)
	c.Set("ItemsGet", Item, cache.DefaultExpiration)
	return c
}
