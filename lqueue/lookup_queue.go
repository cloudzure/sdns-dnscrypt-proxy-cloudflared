package lqueue

import (
	"sync"
	"time"
)

// LQueue type
type LQueue struct {
	mu sync.RWMutex

	l map[uint64]chan struct{}

	duration time.Duration
}

// New func
func New(duration time.Duration) *LQueue {
	return &LQueue{
		l: make(map[uint64]chan struct{}),

		duration: duration,
	}
}

// Get func
func (q *LQueue) Get(key uint64) <-chan struct{} {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if c, ok := q.l[key]; ok {
		return c
	}

	return nil
}

// Wait func
func (q *LQueue) Wait(key uint64) {
	q.mu.RLock()

	if c, ok := q.l[key]; ok {
		q.mu.RUnlock()
		select {
		case <-c:
		case <-time.After(q.duration):
		}
		return
	}

	q.mu.RUnlock()
}

// Add func
func (q *LQueue) Add(key uint64) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.l[key] = make(chan struct{})
}

// Done func
func (q *LQueue) Done(key uint64) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if c, ok := q.l[key]; ok {
		close(c)
	}

	delete(q.l, key)
}
