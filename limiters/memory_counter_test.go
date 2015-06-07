package limiters_test

import (
	"sync"
	"testing"

	"github.com/garslo/go-drlim/limiters"
)

func TestMemoryCounterBuilderGivesGoodDefaults(t *testing.T) {
	counter := limiters.NewMemoryCounterBuilder().Build()
	count := counter.Count("foo")
	if count != 0 {
		t.Errorf("Expected default count to be zero, got %d", count)
	}
}

func TestWeCanZero(t *testing.T) {
	counter := limiters.NewMemoryCounter()
	host := "foo"
	counter.Inc(host)
	counter.Zero()
	count := counter.Count(host)
	if count != 0 {
		t.Errorf("Expected count to be zero, got %d", count)
	}
}

func TestThatIncActuallyIncrements(t *testing.T) {
	counter := limiters.NewMemoryCounter()
	host := "foo"
	counter.Inc(host)
	count := counter.Count(host)
	if count != 1 {
		t.Errorf("Expected count to increment to 1, got %d", count)
	}
}

func TestThatUpdateZerosTheCounts(t *testing.T) {
	counter := limiters.NewMemoryCounter()
	host := "foo"
	counter.Inc("foo")
	counter.Update()
	count := counter.Count(host)
	if count != 0 {
		t.Errorf("Expected count to be zero, got %d", count)
	}
}

func BenchmarkIncrementing(b *testing.B) {
	counter := limiters.NewMemoryCounter()
	for i := 0; i < b.N; i++ {
		counter.Inc("foo")
	}
}

func BenchmarkCounting(b *testing.B) {
	counter := limiters.NewMemoryCounter()
	for i := 0; i < b.N; i++ {
		counter.Count("foo")
	}
}

func BenchmarkZeroing(b *testing.B) {
	counter := limiters.NewMemoryCounter()
	for i := 0; i < b.N; i++ {
		counter.Zero()
	}
}

func BenchmarkIncrementingConcurrently(b *testing.B) {
	numActors := 1000
	counter := limiters.NewMemoryCounter()
	wg := &sync.WaitGroup{}
	wg.Add(numActors)
	for i := 0; i < numActors; i++ {
		go func() {
			for j := 0; j < b.N/numActors; j++ {
				counter.Inc("foo")
				counter.Count("foo")
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
