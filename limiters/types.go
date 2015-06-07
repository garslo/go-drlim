package limiters

import "net/http"

type Tracker interface {
	Inc(host string)
	LimitExceeded(host string) bool
}

type Counter interface {
	Zero()
	Update()
	Count(host string) int
	Inc(host string)
}

type Limiter interface {
	http.Handler
}
