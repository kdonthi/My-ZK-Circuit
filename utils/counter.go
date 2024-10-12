package utils

import "sync"

type Counter struct {
	mu  *sync.Mutex
	num int
}

func NewCounter() *Counter {
	return &Counter{
		num: 0,
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
