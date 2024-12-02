// internal/middleware/rate_limiter.go
package middleware

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

type RateLimiter struct {
	redis     *redis.Client
	maxReqs   int
	window    time.Duration
	keyPrefix string
}

func NewRateLimiter(redis *redis.Client, maxReqs int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		redis:     redis,
		maxReqs:   maxReqs,
		window:    window,
		keyPrefix: "ratelimit:",
	}
}

func (rl *RateLimiter) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get client IP
		key := fmt.Sprintf("%s%s", rl.keyPrefix, c.IP())

		// Use Redis pipeline for atomic operations
		pipe := rl.redis.Pipeline()
		incr := pipe.Incr(c.Context(), key)
		pipe.Expire(c.Context(), key, rl.window)

		// Execute pipeline
		_, err := pipe.Exec(c.Context())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Rate limiting error",
			})
		}

		// Check count
		count, err := incr.Result()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Rate limiting error",
			})
		}

		// Set headers
		c.Set("X-RateLimit-Limit", fmt.Sprintf("%d", rl.maxReqs))
		c.Set("X-RateLimit-Remaining", fmt.Sprintf("%d", rl.maxReqs-int(count)))
		c.Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(rl.window).Unix()))

		if count > int64(rl.maxReqs) {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests",
			})
		}

		return c.Next()
	}
}
