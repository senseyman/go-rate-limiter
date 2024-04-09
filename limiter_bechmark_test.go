package limiter

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func BenchmarkRateLimiterMiddleware(b *testing.B) {
	storage := NewInMemoryStorage()
	rateLimiter := NewRateLimiter(2, time.Minute, storage)

	req := httptest.NewRequest("GET", "/", http.NoBody)
	req.RemoteAddr = testLocalhostAddr

	recorder := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rateLimiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		})).ServeHTTP(recorder, req)
	}
}
