package async

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Promise represents TypeScript's Promise<T>
type Promise[T any] struct {
	result chan T
	err    chan error
	done   chan bool
	mu     sync.RWMutex
	state  PromiseState
	value  T
	error  error
}

// PromiseState represents the state of a Promise
type PromiseState int

const (
	Pending PromiseState = iota
	Fulfilled
	Rejected
)

func (s PromiseState) String() string {
	switch s {
	case Pending:
		return "Pending"
	case Fulfilled:
		return "Fulfilled"
	case Rejected:
		return "Rejected"
	default:
		return "Unknown"
	}
}

// Executor function type for Promise constructor
type Executor[T any] func() (T, error)

// NewPromise creates a new Promise (like new Promise() in TypeScript)
func NewPromise[T any](executor Executor[T]) *Promise[T] {
	p := &Promise[T]{
		result: make(chan T, 1),
		err:    make(chan error, 1),
		done:   make(chan bool, 1),
		state:  Pending,
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				p.mu.Lock()
				p.state = Rejected
				p.error = fmt.Errorf("panic: %v", r)
				p.mu.Unlock()
				p.err <- p.error
				p.done <- true
			}
		}()

		result, err := executor()
		p.mu.Lock()
		if err != nil {
			p.state = Rejected
			p.error = err
			p.mu.Unlock()
			p.err <- err
		} else {
			p.state = Fulfilled
			p.value = result
			p.mu.Unlock()
			p.result <- result
		}
		p.done <- true
	}()

	return p
}

// Resolve creates a resolved Promise (like Promise.resolve() in TypeScript)
func Resolve[T any](value T) *Promise[T] {
	p := &Promise[T]{
		result: make(chan T, 1),
		err:    make(chan error, 1),
		done:   make(chan bool, 1),
		state:  Fulfilled,
		value:  value,
	}
	p.result <- value
	p.done <- true
	return p
}

// Reject creates a rejected Promise (like Promise.reject() in TypeScript)
func Reject[T any](err error) *Promise[T] {
	p := &Promise[T]{
		result: make(chan T, 1),
		err:    make(chan error, 1),
		done:   make(chan bool, 1),
		state:  Rejected,
		error:  err,
	}
	p.err <- err
	p.done <- true
	return p
}

// Then chains promises (like .then() in TypeScript)
func Then[T, U any](p *Promise[T], onFulfilled func(T) U, onRejected func(error) U) *Promise[U] {
	return NewPromise[U](func() (U, error) {
		select {
		case result := <-p.result:
			if onFulfilled != nil {
				return onFulfilled(result), nil
			}
			var zero U
			return zero, nil
		case err := <-p.err:
			if onRejected != nil {
				return onRejected(err), nil
			}
			var zero U
			return zero, err
		}
	})
}

// ThenPromise chains promises that return promises (like .then() returning Promise)
func ThenPromise[T, U any](p *Promise[T], onFulfilled func(T) *Promise[U]) *Promise[U] {
	return NewPromise[U](func() (U, error) {
		result, err := p.Await()
		if err != nil {
			var zero U
			return zero, err
		}
		if onFulfilled != nil {
			nextPromise := onFulfilled(result)
			return nextPromise.Await()
		}
		var zero U
		return zero, nil
	})
}

// Catch handles promise rejection (like .catch() in TypeScript)
func Catch[T any](p *Promise[T], onRejected func(error) T) *Promise[T] {
	return NewPromise[T](func() (T, error) {
		select {
		case result := <-p.result:
			return result, nil
		case err := <-p.err:
			if onRejected != nil {
				return onRejected(err), nil
			}
			var zero T
			return zero, err
		}
	})
}

// Finally executes code regardless of promise outcome (like .finally() in TypeScript)
func Finally[T any](p *Promise[T], onFinally func()) *Promise[T] {
	return NewPromise[T](func() (T, error) {
		defer func() {
			if onFinally != nil {
				onFinally()
			}
		}()
		return p.Await()
	})
}

// Await waits for the promise to resolve (like await in TypeScript)
func (p *Promise[T]) Await() (T, error) {
	<-p.done
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	if p.state == Fulfilled {
		return p.value, nil
	}
	var zero T
	return zero, p.error
}

// AwaitWithTimeout waits for promise with timeout
func (p *Promise[T]) AwaitWithTimeout(timeout time.Duration) (T, error) {
	select {
	case <-p.done:
		return p.Await()
	case <-time.After(timeout):
		var zero T
		return zero, fmt.Errorf("promise timeout after %v", timeout)
	}
}

// AwaitWithContext waits for promise with context cancellation
func (p *Promise[T]) AwaitWithContext(ctx context.Context) (T, error) {
	select {
	case <-p.done:
		return p.Await()
	case <-ctx.Done():
		var zero T
		return zero, ctx.Err()
	}
}

// GetState returns the current state of the promise
func (p *Promise[T]) GetState() PromiseState {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.state
}

// IsPending returns true if promise is still pending
func (p *Promise[T]) IsPending() bool {
	return p.GetState() == Pending
}

// IsFulfilled returns true if promise is fulfilled
func (p *Promise[T]) IsFulfilled() bool {
	return p.GetState() == Fulfilled
}

// IsRejected returns true if promise is rejected
func (p *Promise[T]) IsRejected() bool {
	return p.GetState() == Rejected
}

// All waits for all promises to resolve (like Promise.all() in TypeScript)
func All[T any](promises ...*Promise[T]) *Promise[[]T] {
	return NewPromise[[]T](func() ([]T, error) {
		results := make([]T, len(promises))
		errors := make([]error, len(promises))
		var wg sync.WaitGroup
		
		for i, promise := range promises {
			wg.Add(1)
			go func(index int, p *Promise[T]) {
				defer wg.Done()
				result, err := p.Await()
				if err != nil {
					errors[index] = err
				} else {
					results[index] = result
				}
			}(i, promise)
		}
		
		wg.Wait()
		
		// Check for errors
		for _, err := range errors {
			if err != nil {
				return nil, err
			}
		}
		
		return results, nil
	})
}

// AllSettled waits for all promises to settle (like Promise.allSettled() in TypeScript)
func AllSettled[T any](promises ...*Promise[T]) *Promise[[]PromiseResult[T]] {
	return NewPromise[[]PromiseResult[T]](func() ([]PromiseResult[T], error) {
		results := make([]PromiseResult[T], len(promises))
		var wg sync.WaitGroup
		
		for i, promise := range promises {
			wg.Add(1)
			go func(index int, p *Promise[T]) {
				defer wg.Done()
				result, err := p.Await()
				if err != nil {
					results[index] = PromiseResult[T]{
						Status: "rejected",
						Reason: err,
					}
				} else {
					results[index] = PromiseResult[T]{
						Status: "fulfilled",
						Value:  result,
					}
				}
			}(i, promise)
		}
		
		wg.Wait()
		return results, nil
	})
}

// PromiseResult represents the result of a settled promise
type PromiseResult[T any] struct {
	Status string `json:"status"` // "fulfilled" or "rejected"
	Value  T      `json:"value,omitempty"`
	Reason error  `json:"reason,omitempty"`
}

// Race returns the first promise to settle (like Promise.race() in TypeScript)
func Race[T any](promises ...*Promise[T]) *Promise[T] {
	return NewPromise[T](func() (T, error) {
		result := make(chan T, 1)
		err := make(chan error, 1)
		
		for _, promise := range promises {
			go func(p *Promise[T]) {
				res, e := p.Await()
				if e != nil {
					select {
					case err <- e:
					default:
					}
				} else {
					select {
					case result <- res:
					default:
					}
				}
			}(promise)
		}
		
		select {
		case res := <-result:
			return res, nil
		case e := <-err:
			var zero T
			return zero, e
		}
	})
}

// Any returns the first fulfilled promise (like Promise.any() in TypeScript)
func Any[T any](promises ...*Promise[T]) *Promise[T] {
	return NewPromise[T](func() (T, error) {
		result := make(chan T, 1)
		errors := make([]error, 0, len(promises))
		var mu sync.Mutex
		var wg sync.WaitGroup
		
		for _, promise := range promises {
			wg.Add(1)
			go func(p *Promise[T]) {
				defer wg.Done()
				res, err := p.Await()
				if err != nil {
					mu.Lock()
					errors = append(errors, err)
					mu.Unlock()
				} else {
					select {
					case result <- res:
					default:
					}
				}
			}(promise)
		}
		
		go func() {
			wg.Wait()
			close(result)
		}()
		
		if res, ok := <-result; ok {
			return res, nil
		}
		
		var zero T
		return zero, fmt.Errorf("all promises rejected: %v", errors)
	})
}

// Sleep creates a promise that resolves after a duration (like setTimeout in TypeScript)
func Sleep[T any](duration time.Duration, value T) *Promise[T] {
	return NewPromise[T](func() (T, error) {
		time.Sleep(duration)
		return value, nil
	})
}

// Delay creates a promise that resolves with void after a duration
func Delay(duration time.Duration) *Promise[interface{}] {
	return Sleep[interface{}](duration, nil)
}

// Timeout wraps a promise with a timeout
func Timeout[T any](promise *Promise[T], duration time.Duration) *Promise[T] {
	return Race(promise, NewPromise[T](func() (T, error) {
		time.Sleep(duration)
		var zero T
		return zero, fmt.Errorf("operation timed out after %v", duration)
	}))
}