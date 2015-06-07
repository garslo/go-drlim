package main

import (
	"net/http"
	"time"

	"github.com/garslo/drlim/limiters"
)

func main() {
	counts := limiters.NewMemoryCounter()
	limits := limiters.NewMemoryCounterBuilder().
		WithDefault(3).
		Build()
	limiter := limiters.DistributedLimiter{
		Tracker: limiters.NewExpiringTracker(5*time.Second, limits, counts),
	}

	mux := http.NewServeMux()
	handler := limiter.MakeMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("you got through!"))
		}),
	)
	mux.Handle("/", handler)
	http.ListenAndServe(":8989", mux)
}
