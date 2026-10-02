package main

import "sync"

func NewCounter() *Counter {
	return &Counter{}
}

type Counter struct {
	mu     sync.Mutex
	dCount int
}

func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.dCount++
}

func (c *Counter) Value() int {
	return c.dCount
}
