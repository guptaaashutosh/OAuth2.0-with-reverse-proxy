package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetRedisClientToContext(redisClient *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("redis-client", redisClient)
		ctx.Next()
	}

}
