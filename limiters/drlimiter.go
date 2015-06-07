package limiters

import (
	"errors"
	"log"
	"net"
	"net/http"
)

//type LimitTracker struct {
//	Lock *sync.RWMutex
//	// Per minute
//	limits map[string]int
//	counts map[string]int
//}
//
//func NewLimitTracker(limits map[string]int) *LimitTracker {
//	return &LimitTracker{
//		Lock:   &sync.RWMutex{},
//		limits: limits,
//		counts: make(map[string]int),
//	}
//}
//
//func (me *LimitTracker) Inc(host string) {
//	me.Lock.Lock()
//	me.Counts[host] += 1
//	me.Lock.Unlock()
//}
//
//func (me *LimitTracker) LimitExceeded(host string) bool {
//	me.Lock.RLock()
//	// Could defer the RUnlock(), but that's slower & causes alloc's
//	count, ok := me.Counts[host]
//	if !ok {
//		me.Lock.RUnlock()
//		return false
//	}
//	limit, ok := me.Limits[host]
//	if !ok {
//		me.Lock.RUnlock()
//		return false
//	}
//	me.Lock.RUnlock()
//	log.Printf("host=%s count=%d limit=%d", host, count, limit)
//	return count >= limit
//}

var (
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
)

const (
	// Not in Go stdlib??
	StatusTooManyRequests = 429
)

type DistributedLimiter struct {
	Tracker Tracker
}

func (me *DistributedLimiter) MakeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := me.checkLimits(w, r); err != nil {
			log.Printf("Request denied: %v", err)
			w.WriteHeader(StatusTooManyRequests)
			w.Write([]byte(err.Error()))
			return
		}
		log.Printf("Request allowed")
		next.ServeHTTP(w, r)
	})
}

func (me *DistributedLimiter) checkLimits(w http.ResponseWriter, r *http.Request) error {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// hrm...
		log.Printf("Problem getting remote addr: %v", err)
		return nil
	}
	log.Printf("Dealing with host %s", host)
	if me.Tracker.LimitExceeded(host) {
		return ErrRateLimitExceeded
	}
	me.Tracker.Inc(host)
	return nil
}
