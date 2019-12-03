package common

import "sync/atomic"

type CounterInt64 struct {
	value int64
}

func NewCounterInt64() CounterInt64 {
	return CounterInt64{}
}

func (c *CounterInt64) NextValue() int64 {
	return atomic.AddInt64(&c.value, 1)
}
