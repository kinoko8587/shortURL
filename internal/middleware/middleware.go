package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type TokenBucket struct {
	capacity     int
	tokens       int
	refillRate   int
	refillPeriod time.Duration
	lastRefill   time.Time
	mu           sync.Mutex
}

func NewTokenBucket(capacity, refillRate int, refillPeriod time.Duration) *TokenBucket {
	return &TokenBucket{
		capacity:     capacity,
		tokens:       capacity,
		refillRate:   refillRate,
		refillPeriod: refillPeriod,
		lastRefill:   time.Now(),
	}
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)

	if elapsed >= tb.refillPeriod {
		refills := int(elapsed / tb.refillPeriod)
		tb.tokens += refills * tb.refillRate
		if tb.tokens > tb.capacity {
			tb.tokens = tb.capacity
		}
		tb.lastRefill = now
	}

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}

// RateLimitMiddlewareByIPAndUser limits by IP and optionally by user ID (from header "X-User-ID")
func RateLimitMiddlewareByIPAndUser(capacity, refillRate int, refillPeriod time.Duration) gin.HandlerFunc {
	buckets := sync.Map{} // map[string]*TokenBucket

	return func(c *gin.Context) {
		ip := getIP(c)
		userID := c.GetHeader("X-User-ID") // fallback to IP if no user ID
		key := ip
		if userID != "" {
			key = ip + ":" + userID
		}

		val, _ := buckets.LoadOrStore(key, NewTokenBucket(capacity, refillRate, refillPeriod))
		tb := val.(*TokenBucket)

		if !tb.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}

		c.Next()
	}
}

func getIP(c *gin.Context) string {
	// handle reverse proxy or load balancer
	ip := c.ClientIP()
	if idx := strings.Index(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	return ip
}
