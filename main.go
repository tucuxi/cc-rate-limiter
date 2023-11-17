package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/tucuxi/cc-rate-limiter/internal/ratelimit"
)

var limiter = ratelimit.NewFixedWindowLimiter(60, 60)

func main() {
	http.HandleFunc("/limited", limitedHandler)
	http.HandleFunc("/unlimited", unlimitedHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func limitedHandler(w http.ResponseWriter, r *http.Request) {
	if !limiter.Take() {
		w.WriteHeader(http.StatusTooManyRequests)
		return
	}
	fmt.Fprintln(w, "Limited")
}

func unlimitedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "No limits!")
}
