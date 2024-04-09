package limiter

import (
	"time"
)

type Storage interface {
	Get(key string) (int, error)
	Set(key string, value int, expiration time.Duration) error
}
