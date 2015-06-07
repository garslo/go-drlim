package limiters

import (
	"sync"
	"time"
)

type expiringTracker struct {
	lock   *sync.RWMutex
	counts Counter
	limits Counter
}

func NewExpiringTracker(expireInterval time.Duration, limits, counts Counter) Tracker {
	t := &expiringTracker{
		lock:   &sync.RWMutex{},
		limits: limits,
		counts: counts,
	}
	go t.manageExpiration(expireInterval)
	return t
}

func (me *expiringTracker) manageExpiration(expireInterval time.Duration) {
	for {
		<-time.Tick(expireInterval)
		me.counts.Update()
	}
}

func (me *expiringTracker) Inc(host string) {
	me.counts.Inc(host)
}

func (me *expiringTracker) LimitExceeded(host string) bool {
	count := me.counts.Count(host)
	limit := me.limits.Count(host)
	return count >= limit
}
