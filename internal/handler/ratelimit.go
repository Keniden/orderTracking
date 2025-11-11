package handler

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/redis/go-redis/v9"
    "context"
)

// Simple token-bucket style limiter using Redis INCR + EXPIRE per key per window
func rateLimiter(rdb *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
    if rdb == nil {
        return func(c *gin.Context) { c.Next() }
    }
    return func(c *gin.Context) {
        // key can be per user or per token; fallback to IP
        token := c.GetHeader("Authorization")
        if token == "" { token = c.ClientIP() }
        key := "ratelimit:" + token + ":" + time.Now().UTC().Format("20060102T15:04")
        ctx := context.Background()
        n, err := rdb.Incr(ctx, key).Result()
        if err != nil {
            c.Next(); return
        }
        if n == 1 {
            rdb.Expire(ctx, key, window)
        }
        if n > int64(limit) {
            c.AbortWithStatusJSON(http.StatusTooManyRequests, ErrorResponse{Message: "rate limit exceeded"})
            return
        }
        c.Next()
    }
}
