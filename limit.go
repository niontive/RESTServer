package main

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

// Limit requests to one per millisecond
var limiter = rate.NewLimiter(rate.Every(time.Millisecond), 1)

func limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if limiter.Allow() == false {
			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
