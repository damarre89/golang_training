package main

import (
	"sync"
	"testing"
)

func assertAtCounter(t testing.TB, got *AtCounter, want int) {
	t.Helper()
	if got.Value() != want {
		t.Errorf("got %d, want %d", got.Value(), want)
	}
}

func TestAtCounter(t *testing.T) {
	t.Run("incrementing the counter 3 times leaves it at 3", func(t *testing.T) {
		counter := NewAtCounter()
		counter.Inc()
		counter.Inc()
		counter.Inc()

		assertAtCounter(t, counter, 3)
	})
	t.Run("it runs safely concurrently", func(t *testing.T) {
		wantedCount := 1000
		counter := NewAtCounter()

		var wg sync.WaitGroup
		wg.Add(wantedCount)

		for i := 0; i < wantedCount; i++ {
			go func() {
				counter.Inc()
				wg.Done()
			}()
		}
		wg.Wait()

		assertAtCounter(t, counter, wantedCount)
	})
}
