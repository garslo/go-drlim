package limiters

import "sync"

type MemoryCounterBuilder struct {
	counter *memoryCounter
}

func NewMemoryCounterBuilder() MemoryCounterBuilder {
	return MemoryCounterBuilder{
		counter: newMemoryCounter(make(map[string]int)),
	}
}

func (me MemoryCounterBuilder) WithCounts(counts map[string]int) MemoryCounterBuilder {
	me.counter.counts = counts
	return me
}

func (me MemoryCounterBuilder) WithDefault(defaultCount int) MemoryCounterBuilder {
	me.counter.defaultCount = defaultCount
	return me
}

func (me MemoryCounterBuilder) Build() Counter {
	return me.counter
}

type memoryCounter struct {
	lock         *sync.RWMutex
	counts       map[string]int
	defaultCount int
}

func newMemoryCounter(initial map[string]int) *memoryCounter {
	return &memoryCounter{
		counts: initial,
		lock:   &sync.RWMutex{},
	}
}

func NewMemoryCounter() Counter {
	return NewMemoryCounter2(make(map[string]int))
}

func NewMemoryCounter2(initial map[string]int) Counter {
	return newMemoryCounter(initial)
}

func (me *memoryCounter) Zero() {
	me.lock.Lock()
	for key, _ := range me.counts {
		me.counts[key] = 0
	}
	me.lock.Unlock()
}

func (me *memoryCounter) Update() {
	me.Zero()
}

func (me *memoryCounter) Count(host string) int {
	me.lock.RLock()
	counts, ok := me.counts[host]
	me.lock.RUnlock()
	if !ok {
		return me.defaultCount
	}
	return counts
}

func (me *memoryCounter) Inc(host string) {
	me.lock.Lock()
	me.counts[host] += 1
	me.lock.Unlock()
}
