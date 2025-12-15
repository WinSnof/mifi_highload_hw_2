package utils

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

// limiter: 1000 req/s, burst 5000
var limiter = rate.NewLimiter(rate.Limit(1000), 5000)

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
		defer cancel()

		if err := limiter.WaitN(ctx, 1); err != nil {
			http.Error(w, "Too Many Requests (Queue Timeout)", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
