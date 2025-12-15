package service

import "sync"

// Resettable generic parameter restriction.
type Resettable interface {
	Reset()
}

// Pool generic container for objects of the same type.
type Pool[T Resettable] struct {
	pool sync.Pool
}

// New create new Pool.
func New[T Resettable](newFn func() T) *Pool[T] {
	return &Pool[T]{
		pool: sync.Pool{
			New: func() any {
				return newFn()
			},
		},
	}
}

// Get returns an object from the pool.
func (p *Pool[T]) Get() T {
	return p.pool.Get().(T)
}

// Put resets the state of an object and returns it to the pool.
func (p *Pool[T]) Put(obj T) {
	obj.Reset()
	p.pool.Put(obj)
}
