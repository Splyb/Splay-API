package middleware

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

var limiter = NewRateLimiter(5, time.Minute)

func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Demasiadas solicitudes, inténtalo más tarde", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type RateLimiter struct {
	limiter  *rate.Limiter
	lastUsed time.Time
	mu       sync.Mutex
}

func NewRateLimiter(r int, d time.Duration) *RateLimiter {
	return &RateLimiter{
		limiter: rate.NewLimiter(rate.Every(d/time.Duration(r)), r),
	}
}

func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	if rl.limiter.Allow() {
		rl.lastUsed = time.Now()
		return true
	}
	return false
}
