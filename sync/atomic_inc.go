package main

import "sync/atomic"

func NewAtCounter() *AtCounter {
	return &AtCounter{}
}

type AtCounter struct {
	dCount atomic.Int64
}

func (c *AtCounter) Inc() {
	c.dCount.Add(1)
}

func (c *AtCounter) Value() int {
	return int(c.dCount.Load())
}
