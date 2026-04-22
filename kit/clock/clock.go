package clock

import (
	"sync"
	"time"
)

var (
	mu   sync.RWMutex
	now  currentTime
	once sync.Once
)

type currentTime func() time.Time

func initialize() {
	once.Do(func() {
		mu.Lock()
		defer mu.Unlock()
		if now == nil {
			now = time.Now
		}
	})
}

func Now() time.Time {
	initialize()
	mu.RLock()
	defer mu.RUnlock()
	return now()
}

// With replaces the time source and returns a restore function.
// Use in tests: defer clock.With(func() time.Time { return fixed })()
func With(current currentTime) func() {
	mu.Lock()
	defer mu.Unlock()
	now = current
	once = sync.Once{}
	return initialize
}
