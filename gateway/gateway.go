package gateway

import (
	"bnsportal/storage"
	"time"

	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type Gateway struct {
	Port       string
	Domain     string
	DefaultRPC string
	Storage    *storage.Storage
	RedisCfg   *redis.Options
}

func (gw *Gateway) Start() error {
	var cacheStore persist.CacheStore

	cacheStore = persist.NewMemoryStore(1 * time.Minute)

	if gw.RedisCfg != nil {
		cacheStore = persist.NewRedisStore(redis.NewClient(gw.RedisCfg))
	}

	r := gin.Default()

	r.Use(cors.Default())

	r.Use(gzip.Gzip(gzip.DefaultCompression))

	apiv1 := r.Group("/api/v1")

	apiv1.GET("/name/:name", cache.CacheByRequestURI(cacheStore, 5*time.Second), gw.GetName)
	apiv1.GET("/names/:address", cache.CacheByRequestURI(cacheStore, 5*time.Second), gw.GetAddressNames)
	apiv1.GET("/names", cache.CacheByRequestURI(cacheStore, 1*time.Second), gw.GetNames)
	apiv1.GET("/available/:name", cache.CacheByRequestURI(cacheStore, 1*time.Second), gw.CheckNameAvailable)

	err := r.Run("0.0.0.0:" + gw.Port)
	if err != nil {
		return err
	}
	return nil
}
