package cache

import (
	"sync"

	"github.com/darksuei/suei-intelligence/internal/config"
	domain "github.com/darksuei/suei-intelligence/internal/domain/cache"
	"github.com/darksuei/suei-intelligence/internal/infrastructure/cache/memory"
	"github.com/darksuei/suei-intelligence/internal/infrastructure/cache/redis"
)

var (
	instance domain.Cache
	once     sync.Once
)

// GetCache returns a singleton cache instance
func GetCache() domain.Cache {

	once.Do(func() {
		switch config.Cache().CacheType {
			case domain.CacheTypeRedis:
				instance = redis.NewCache(config.Cache())
			case domain.CacheTypeMemory:
				instance = memory.NewCache()
			default:
				instance = memory.NewCache()
		}
	})
	return instance
}