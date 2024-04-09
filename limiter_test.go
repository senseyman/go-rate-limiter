package limiter

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	testLocalhostAddr = "127.0.0.1"
)

func TestRateLimiter_Middleware(t *testing.T) {
	testCases := []struct {
		name           string
		getStorage     func() Storage
		getRateLimiter func(maxRequests int, duration time.Duration, storage Storage) *RateLimiter
		hasErr         bool
		exec           func(limiter *RateLimiter) error
	}{
		{
			name: "success/OK",
			getStorage: func() Storage {
				return NewInMemoryStorage()
			},
			getRateLimiter: NewRateLimiter,
			hasErr:         false,
			exec: func(limiter *RateLimiter) error {
				for i := 0; i < 2; i++ {
					req := httptest.NewRequest("GET", "/", http.NoBody)
					req.RemoteAddr = testLocalhostAddr

					recorder := httptest.NewRecorder()

					limiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
						w.WriteHeader(http.StatusOK)
					})).ServeHTTP(recorder, req)

					if recorder.Code != http.StatusOK {
						return fmt.Errorf("expected status code %d; got %d", http.StatusOK, recorder.Code)
					}
				}
				return nil
			},
		},
		{
			name: "success/TooManyRequests",
			getStorage: func() Storage {
				return NewInMemoryStorage()
			},
			getRateLimiter: NewRateLimiter,
			hasErr:         false,
			exec: func(limiter *RateLimiter) error {
				for i := 0; i < 2; i++ {
					req := httptest.NewRequest("GET", "/", http.NoBody)
					req.RemoteAddr = fmt.Sprintf("192.168.0.%d", i) // Different IP address for each request

					recorder := httptest.NewRecorder()

					limiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
						w.WriteHeader(http.StatusOK)
					})).ServeHTTP(recorder, req)

					if recorder.Code != http.StatusOK {
						return fmt.Errorf("expected status code %d; got %d", http.StatusOK, recorder.Code)
					}
				}
				return nil
			},
		},
		{
			name: "err/TooManyRequests",
			getStorage: func() Storage {
				return NewInMemoryStorage()
			},
			getRateLimiter: NewRateLimiter,
			hasErr:         true,
			exec: func(limiter *RateLimiter) error {
				for i := 0; i < 10; i++ {
					req := httptest.NewRequest("GET", "/", http.NoBody)
					req.RemoteAddr = testLocalhostAddr

					recorder := httptest.NewRecorder()

					limiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
						w.WriteHeader(http.StatusOK)
					})).ServeHTTP(recorder, req)

					if recorder.Code != http.StatusOK {
						return fmt.Errorf("expected status code %d; got %d", http.StatusOK, recorder.Code)
					}
				}
				return nil
			},
		},
		{
			name: "success/InternalServerError",
			getStorage: func() Storage {
				return &MockStorageWithError{}
			},
			getRateLimiter: NewRateLimiter,
			hasErr:         false,
			exec: func(limiter *RateLimiter) error {
				req := httptest.NewRequest("GET", "/", http.NoBody)
				req.RemoteAddr = testLocalhostAddr

				recorder := httptest.NewRecorder()

				limiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusOK)
				})).ServeHTTP(recorder, req)

				if recorder.Code != http.StatusInternalServerError {
					return fmt.Errorf("expected status code %d; got %d", http.StatusInternalServerError, recorder.Code)
				}

				return nil
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storage := tc.getStorage()
			rl := tc.getRateLimiter(2, time.Minute, storage)
			err := tc.exec(rl)
			if tc.hasErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// MockStorageWithError mocks RateLimitStorage interface with an error.
type MockStorageWithError struct{}

func (MockStorageWithError) Get(_ string) (int, error) {
	return 0, nil
}

func (MockStorageWithError) Set(_ string, _ int, _ time.Duration) error {
	return fmt.Errorf("mock error")
}
