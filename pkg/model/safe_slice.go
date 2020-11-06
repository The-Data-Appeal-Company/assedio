package model

import (
	"sync"
)

type Slice interface {
	Append(Record)
	Get(i int) Record
	Len() int
}

type ThreadSafeSlice struct {
	mu    *sync.RWMutex
	slice []Record
}

func NewThreadSafeSlice() *ThreadSafeSlice {
	return &ThreadSafeSlice{
		mu:    &sync.RWMutex{},
		slice: make([]Record, 0),
	}
}

func (t *ThreadSafeSlice) Append(value Record) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.slice = append(t.slice, value)
}

func (t *ThreadSafeSlice) Get(i int) Record {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.slice[i]
}

func (t *ThreadSafeSlice) Len() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.slice)
}
