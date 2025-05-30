package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	ips         map[string]*rate.Limiter
	mu          *sync.RWMutex
	rate        rate.Limit
	burst       int
	ttl         time.Duration
	lastCleanup time.Time
}

func NewRateLimiter(r rate.Limit, b int, ttl time.Duration) *RateLimiter {
	return &RateLimiter{
		ips:         make(map[string]*rate.Limiter),
		mu:          &sync.RWMutex{},
		rate:        r,
		burst:       b,
		ttl:         ttl,
		lastCleanup: time.Now(),
	}
}

func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.ips[ip] = limiter
	}

	// Cleanup old entries if needed
	if time.Since(rl.lastCleanup) > rl.ttl {
		rl.cleanup()
	}

	return limiter
}

func (rl *RateLimiter) cleanup() {
	// Simple cleanup - just remove all entries
	// In a production environment, you might want to implement a more sophisticated cleanup strategy
	rl.ips = make(map[string]*rate.Limiter)
	rl.lastCleanup = time.Now()
}

func RateLimit(rps float64, burst int) gin.HandlerFunc {
	limiter := NewRateLimiter(rate.Limit(rps), burst, 1*time.Hour)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.getLimiter(ip).Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
