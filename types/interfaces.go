package types

import (
	"encoding/json"
	"fmt"
)

// Comparable interface for types that can be compared
// Equivalent to TypeScript's comparison operators
type Comparable[T any] interface {
	CompareTo(other T) int // -1: less, 0: equal, 1: greater
	Equals(other T) bool
}

// Serializable interface for JSON serialization
// Similar to TypeScript's JSON.stringify/parse
type Serializable interface {
	ToJSON() ([]byte, error)
	FromJSON(data []byte) error
}

// Iterable interface for collection types
// Mimics TypeScript's Iterable<T>
type Iterable[T any] interface {
	Iterator() Iterator[T]
	ForEach(fn func(T))
}

// Iterator interface for iteration protocol
// Similar to TypeScript's Iterator<T>
type Iterator[T any] interface {
	Next() (T, bool) // value, hasValue
	HasNext() bool
}

// Disposable interface for resource cleanup
// Similar to TypeScript's disposable pattern
type Disposable interface {
	Dispose() error
}

// Cloneable interface for object cloning
// Mimics TypeScript's object spreading/cloning
type Cloneable[T any] interface {
	Clone() T
}

// Hashable interface for types that can be hashed
// Similar to TypeScript's Map key requirements
type Hashable interface {
	Hash() uint64
}

// Stringable interface for string conversion
// Equivalent to TypeScript's toString()
type Stringable interface {
	ToString() string
}

// DefaultStringable provides default string conversion
type DefaultStringable struct{}

func (d DefaultStringable) ToString() string {
	return fmt.Sprintf("%+v", d)
}

// JSONSerializable provides default JSON serialization
type JSONSerializable struct{}

func (j JSONSerializable) ToJSON() ([]byte, error) {
	return json.Marshal(j)
}

func (j *JSONSerializable) FromJSON(data []byte) error {
	return json.Unmarshal(data, j)
}

// TypeAssert performs safe type assertion
// Similar to TypeScript's type guards
func TypeAssert[T any](value interface{}) (T, bool) {
	if v, ok := value.(T); ok {
		return v, true
	}
	var zero T
	return zero, false
}

// ImplementsInterface checks if a value implements an interface
// Similar to TypeScript's instanceof checks
func ImplementsInterface[T any](value interface{}) bool {
	_, ok := value.(T)
	return ok
}

// IsType checks if a value is of a specific type
func IsType[T any](value interface{}) bool {
	_, ok := value.(T)
	return ok
}

// AsType safely converts interface{} to specific type
func AsType[T any](value interface{}) *T {
	if v, ok := value.(T); ok {
		return &v
	}
	return nil
}