package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/snavarro/microtracker/config"
	"golang.org/x/time/rate"
)

// RateLimiter represents a rate limiter
type RateLimiter struct {
	config   *config.RateLimitConfig
	limiters map[string]map[string]*rateLimiterInfo
	mu       *sync.RWMutex
}

type rateLimiterInfo struct {
	limiter    *rate.Limiter
	lastAccess time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(cfg *config.RateLimitConfig) *RateLimiter {
	limiter := &RateLimiter{
		config:   cfg,
		limiters: make(map[string]map[string]*rateLimiterInfo),
		mu:       &sync.RWMutex{},
	}

	// Start cleanup routine
	go limiter.cleanupLoop()

	return limiter
}

// cleanupLoop periodically removes old rate limiters
func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		for endpoint, ipLimiters := range rl.limiters {
			for ip, info := range ipLimiters {
				if time.Since(info.lastAccess) > time.Duration(rl.getTTL(endpoint))*time.Minute {
					delete(ipLimiters, ip)
				}
			}
			if len(ipLimiters) == 0 {
				delete(rl.limiters, endpoint)
			}
		}
		rl.mu.Unlock()
	}
}

// getLimiter returns the rate limiter for the given endpoint and IP
func (rl *RateLimiter) getLimiter(endpoint, ip string) *rateLimiterInfo {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if _, exists := rl.limiters[endpoint]; !exists {
		rl.limiters[endpoint] = make(map[string]*rateLimiterInfo)
	}

	info, exists := rl.limiters[endpoint][ip]
	if !exists {
		info = &rateLimiterInfo{
			limiter: rate.NewLimiter(
				rate.Limit(float64(rl.getRequestsPerMinute(endpoint))/60.0),
				rl.getBurstSize(endpoint),
			),
			lastAccess: time.Now(),
		}
		rl.limiters[endpoint][ip] = info
	}

	info.lastAccess = time.Now()
	return info
}

// getRequestsPerMinute returns the requests per minute for the given endpoint
func (rl *RateLimiter) getRequestsPerMinute(endpoint string) int {
	if limit, exists := rl.config.Endpoints[endpoint]; exists {
		return limit.RequestsPerMinute
	}
	return rl.config.Default.RequestsPerMinute
}

// getBurstSize returns the burst size for the given endpoint
func (rl *RateLimiter) getBurstSize(endpoint string) int {
	if limit, exists := rl.config.Endpoints[endpoint]; exists {
		return limit.BurstSize
	}
	return rl.config.Default.BurstSize
}

// getTTL returns the TTL in minutes for the given endpoint
func (rl *RateLimiter) getTTL(endpoint string) int {
	if limit, exists := rl.config.Endpoints[endpoint]; exists {
		return limit.TTLMinutes
	}
	return rl.config.Default.TTLMinutes
}

// RateLimit returns a gin middleware for rate limiting
func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		endpoint := fmt.Sprintf("%s:%s", c.Request.Method, c.Request.URL.Path)
		ip := c.ClientIP()
		info := rl.getLimiter(endpoint, ip)

		if !info.limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":    "Rate limit exceeded",
				"success":  false,
				"endpoint": endpoint,
				"limit":    rl.getRequestsPerMinute(endpoint),
				"burst":    rl.getBurstSize(endpoint),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
