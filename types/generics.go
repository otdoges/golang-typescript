package types

import "fmt"

// Optional represents TypeScript's T | undefined
type Optional[T any] struct {
	value   T
	hasValue bool
}

// Some creates an Optional with a value
func Some[T any](value T) Optional[T] {
	return Optional[T]{value: value, hasValue: true}
}

// None creates an empty Optional
func None[T any]() Optional[T] {
	return Optional[T]{hasValue: false}
}

// IsSome returns true if Optional has a value
func (o Optional[T]) IsSome() bool {
	return o.hasValue
}

// IsNone returns true if Optional is empty
func (o Optional[T]) IsNone() bool {
	return !o.hasValue
}

// Get returns the value (panics if None)
func (o Optional[T]) Get() T {
	if !o.hasValue {
		panic("Optional.Get() called on None")
	}
	return o.value
}

// GetOrDefault returns value or default if None
func (o Optional[T]) GetOrDefault(defaultValue T) T {
	if o.hasValue {
		return o.value
	}
	return defaultValue
}

// Map transforms the Optional value
func (o Optional[T]) Map(fn func(T) T) Optional[T] {
	if o.hasValue {
		return Some(fn(o.value))
	}
	return None[T]()
}

// FlatMap chains Optional operations
func FlatMap[T, U any](o Optional[T], fn func(T) Optional[U]) Optional[U] {
	if o.hasValue {
		return fn(o.value)
	}
	return None[U]()
}

// Result represents TypeScript's Promise resolve/reject or Either type
type Result[T, E any] struct {
	value T
	error E
	isOk  bool
}

// Ok creates a successful Result
func Ok[T, E any](value T) Result[T, E] {
	return Result[T, E]{value: value, isOk: true}
}

// Err creates an error Result
func Err[T, E any](error E) Result[T, E] {
	return Result[T, E]{error: error, isOk: false}
}

// IsOk returns true if Result is successful
func (r Result[T, E]) IsOk() bool {
	return r.isOk
}

// IsErr returns true if Result is an error
func (r Result[T, E]) IsErr() bool {
	return !r.isOk
}

// Unwrap returns the value (panics on error)
func (r Result[T, E]) Unwrap() T {
	if !r.isOk {
		panic(fmt.Sprintf("Result.Unwrap() called on error: %v", r.error))
	}
	return r.value
}

// UnwrapOr returns value or default on error
func (r Result[T, E]) UnwrapOr(defaultValue T) T {
	if r.isOk {
		return r.value
	}
	return defaultValue
}

// MapResult transforms successful Result
func MapResult[T, U, E any](r Result[T, E], fn func(T) U) Result[U, E] {
	if r.isOk {
		return Ok[U, E](fn(r.value))
	}
	return Err[U, E](r.error)
}

// Tuple represents TypeScript's tuple types [T, U]
type Tuple2[T, U any] struct {
	First  T
	Second U
}

// NewTuple2 creates a 2-element tuple
func NewTuple2[T, U any](first T, second U) Tuple2[T, U] {
	return Tuple2[T, U]{First: first, Second: second}
}

// Tuple3 represents TypeScript's tuple types [T, U, V]
type Tuple3[T, U, V any] struct {
	First  T
	Second U
	Third  V
}

// NewTuple3 creates a 3-element tuple
func NewTuple3[T, U, V any](first T, second U, third V) Tuple3[T, U, V] {
	return Tuple3[T, U, V]{First: first, Second: second, Third: third}
}

// Predicate represents TypeScript's predicate functions
type Predicate[T any] func(T) bool

// Function represents TypeScript's function types
type Function[T, U any] func(T) U

// Action represents TypeScript's void functions
type Action[T any] func(T)

// Supplier represents TypeScript's no-arg functions
type Supplier[T any] func() T

// Consumer represents TypeScript's consumer functions
type Consumer[T any] func(T)

// BiFunction represents TypeScript's two-argument functions
type BiFunction[T, U, V any] func(T, U) V

// Comparable constraint for generic comparison
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 | ~string
}

// Min returns the minimum of two comparable values
func Min[T Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Max returns the maximum of two comparable values
func Max[T Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}