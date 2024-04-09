package limiter

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	maxRequests int           // Maximum number of requests allowed within the duration.
	duration    time.Duration // duration window for the rate limit.
	storage     Storage
	mutex       sync.Mutex
}

// NewRateLimiter creates a new RateLimiter middleware instance.
func NewRateLimiter(maxRequests int, duration time.Duration, storage Storage) *RateLimiter {
	if storage == nil {
		storage = &InMemoryStorage{}
	}
	return &RateLimiter{
		maxRequests: maxRequests,
		duration:    duration,
		storage:     storage,
	}
}

// Middleware function to enforce rate limiting.
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr
		key := fmt.Sprintf("rate_limit:%s", clientIP)

		rl.mutex.Lock()
		defer rl.mutex.Unlock()

		count, err := rl.storage.Get(key)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if count >= rl.maxRequests {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		// Increment the request count for the client.
		err = rl.storage.Set(key, count+1, rl.duration)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	})
}
