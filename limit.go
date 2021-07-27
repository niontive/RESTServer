package main

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

const rateLimit = 25 // Limit requests to one per 25 milliseconds

var limiter = rate.NewLimiter(rate.Every(rateLimit*time.Millisecond), 1)

//
// Limit request for HTTP server
//
func limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if limiter.Allow() == false {
			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
