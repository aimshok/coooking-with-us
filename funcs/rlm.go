package funcs

import (
	"fmt"
	"net/http"
	"time"
)

// Handles risky operations that may fail
func RiskyOperationHandler(w http.ResponseWriter, r *http.Request) {
	Logger.Info("Executing risky operation")

	err := riskyOperation()
	if err != nil {
		Logger.WithError(err).Error("Risky operation failed")
		http.Error(w, "Risky operation failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Risky operation succeeded"))
}

// Simulates a risky operation with a random failure
func riskyOperation() error {
	if time.Now().Unix()%2 == 0 {
		return fmt.Errorf("simulated error")
	}
	return nil
}

// Middleware to enforce rate limiting
func RateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !rateLimiter.Allow() {
			Logger.Warn("Rate limit exceeded")
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		next(w, r)
	}
}
