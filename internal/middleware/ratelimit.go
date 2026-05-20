package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/rizky/go-scaffold/pkg/response"
)

const cleanupInterval = 10 * time.Minute
const bucketTTL = 30 * time.Minute

type tokenBucket struct {
	tokens     float64
	maxTokens  float64
	refillRate float64
	lastRefill time.Time
	lastAccess time.Time
	mu         sync.Mutex
}

func newTokenBucket(maxTokens, refillRate float64) *tokenBucket {
	return &tokenBucket{
		tokens:     maxTokens,
		maxTokens:  maxTokens,
		refillRate: refillRate,
		lastRefill: time.Now(),
		lastAccess: time.Now(),
	}
}

func (tb *tokenBucket) allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	tb.tokens = min(tb.maxTokens, tb.tokens+elapsed*tb.refillRate)
	tb.lastRefill = now
	tb.lastAccess = now

	if tb.tokens >= 1 {
		tb.tokens--
		return true
	}
	return false
}

type RateLimiter struct {
	buckets map[string]*tokenBucket
	mu      sync.RWMutex
}

func NewRateLimiter() *RateLimiter {
	rl := &RateLimiter{buckets: make(map[string]*tokenBucket)}
	go rl.periodicCleanup()
	return rl
}

func (rl *RateLimiter) periodicCleanup() {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		for key, bucket := range rl.buckets {
			bucket.mu.Lock()
			stale := time.Since(bucket.lastAccess) > bucketTTL
			bucket.mu.Unlock()
			if stale {
				delete(rl.buckets, key)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Limit(maxTokens, refillRate float64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.RemoteAddr

			rl.mu.RLock()
			bucket, exists := rl.buckets[key]
			rl.mu.RUnlock()

			if !exists {
				bucket = newTokenBucket(maxTokens, refillRate)
				rl.mu.Lock()
				rl.buckets[key] = bucket
				rl.mu.Unlock()
			}

			if !bucket.allow() {
				response.Error(w, http.StatusTooManyRequests, "rate_limited", "too many requests")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
