# Rate Limiter Middleware

[![Test and Lint status](https://github.com/senseyman/go-rate-limiter/workflows/Go%20Test%20and%20Lint/badge.svg)](https://github.com/senseyman/go-rate-limiter/actions)

This is a rate limiting middleware written in Go. It allows you to limit the number of requests per minute from a single IP address.

## Usage

To use the rate limiter middleware in your Go application, follow these steps:

1. Install the package using `go get`:

```bash
   go get github.com/senseyman/go-rate-limiter
```

2. Import the package in your code:
```go
import "github.com/senseyman/go-rate-limiter"
```
3. Create a new rate limiter middleware instance with your desired configuration
```go
rateLimiter := ratelimiter.NewRateLimiter(100, time.Minute, yourStorage)
```

4. Attach the rate limiter middleware to your HTTP server's request handler
```go
http.Handle("/", rateLimiter.Middleware(yourHandler))
```
