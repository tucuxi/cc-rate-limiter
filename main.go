package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/tucuxi/cc-rate-limiter/internal/pkg/ratelimit"
)

func main() {
	var limiter = ratelimit.NewSlidingWindowLimiter(ratelimit.PerSecond(5))

	http.HandleFunc("/limited", limit(unlimitedHandler, limiter))
	http.HandleFunc("/unlimited", unlimitedHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func limit(handler http.HandlerFunc, limiter ratelimit.Limiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if limiter.Take() {
			handler(w, r)
		} else {
			w.WriteHeader(http.StatusTooManyRequests)
		}
	}
}

func unlimitedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Done")
}
