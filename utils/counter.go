package utils

import "sync"

// Counter keeps track of an incrementing integer that can be accessed from multiple threads
type Counter struct {
	mu  *sync.Mutex
	num int
}

func NewCounter() *Counter {
	return &Counter{
		num: 1,
		mu:  &sync.Mutex{},
	}
}

func (c *Counter) GetNext() int {
	c.mu.Lock()

	curr := c.num
	c.num++

	c.mu.Unlock()

	return curr
}
