package types

import (
	"fmt"
	"reflect"
)

// Map represents TypeScript's Map<K, V> data structure
type Map[K comparable, V any] struct {
	data map[K]V
	size int
}

// NewMap creates a new Map instance (like new Map() in TypeScript)
func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		data: make(map[K]V),
		size: 0,
	}
}

// NewMapWithEntries creates a Map from initial entries (like new Map(entries))
func NewMapWithEntries[K comparable, V any](entries []Tuple2[K, V]) *Map[K, V] {
	m := NewMap[K, V]()
	for _, entry := range entries {
		m.Set(entry.First, entry.Second)
	}
	return m
}

// Set adds or updates a key-value pair (like map.set() in TypeScript)
func (m *Map[K, V]) Set(key K, value V) *Map[K, V] {
	if _, exists := m.data[key]; !exists {
		m.size++
	}
	m.data[key] = value
	return m
}

// Get retrieves a value by key (like map.get() in TypeScript)
func (m *Map[K, V]) Get(key K) Optional[V] {
	if value, exists := m.data[key]; exists {
		return Some(value)
	}
	return None[V]()
}

// Has checks if a key exists (like map.has() in TypeScript)
func (m *Map[K, V]) Has(key K) bool {
	_, exists := m.data[key]
	return exists
}

// Delete removes a key-value pair (like map.delete() in TypeScript)
func (m *Map[K, V]) Delete(key K) bool {
	if _, exists := m.data[key]; exists {
		delete(m.data, key)
		m.size--
		return true
	}
	return false
}

// Clear removes all entries (like map.clear() in TypeScript)
func (m *Map[K, V]) Clear() {
	m.data = make(map[K]V)
	m.size = 0
}

// Size returns the number of entries (like map.size in TypeScript)
func (m *Map[K, V]) Size() int {
	return m.size
}

// IsEmpty checks if the map is empty
func (m *Map[K, V]) IsEmpty() bool {
	return m.size == 0
}

// Keys returns all keys (like map.keys() in TypeScript)
func (m *Map[K, V]) Keys() []K {
	keys := make([]K, 0, m.size)
	for key := range m.data {
		keys = append(keys, key)
	}
	return keys
}

// Values returns all values (like map.values() in TypeScript)
func (m *Map[K, V]) Values() []V {
	values := make([]V, 0, m.size)
	for _, value := range m.data {
		values = append(values, value)
	}
	return values
}

// Entries returns all key-value pairs (like map.entries() in TypeScript)
func (m *Map[K, V]) Entries() []Tuple2[K, V] {
	entries := make([]Tuple2[K, V], 0, m.size)
	for key, value := range m.data {
		entries = append(entries, NewTuple2(key, value))
	}
	return entries
}

// ForEach iterates over all entries (like map.forEach() in TypeScript)
func (m *Map[K, V]) ForEach(fn func(value V, key K, map_ *Map[K, V])) {
	for key, value := range m.data {
		fn(value, key, m)
	}
}

// Filter creates a new map with entries that pass the test
func (m *Map[K, V]) Filter(predicate func(value V, key K) bool) *Map[K, V] {
	result := NewMap[K, V]()
	for key, value := range m.data {
		if predicate(value, key) {
			result.Set(key, value)
		}
	}
	return result
}

// Map transforms values and returns a new map
func MapTransform[K comparable, V, U any](m *Map[K, V], fn func(value V, key K) U) *Map[K, U] {
	result := NewMap[K, U]()
	for key, value := range m.data {
		result.Set(key, fn(value, key))
	}
	return result
}

// Clone creates a shallow copy of the map
func (m *Map[K, V]) Clone() *Map[K, V] {
	result := NewMap[K, V]()
	for key, value := range m.data {
		result.Set(key, value)
	}
	return result
}

// String returns string representation
func (m *Map[K, V]) String() string {
	return fmt.Sprintf("Map{size: %d}", m.size)
}

// Set represents TypeScript's Set<T> data structure
type Set[T comparable] struct {
	data map[T]struct{}
	size int
}

// NewSet creates a new Set instance (like new Set() in TypeScript)
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		data: make(map[T]struct{}),
		size: 0,
	}
}

// NewSetWithValues creates a Set from initial values (like new Set(values))
func NewSetWithValues[T comparable](values []T) *Set[T] {
	s := NewSet[T]()
	for _, value := range values {
		s.Add(value)
	}
	return s
}

// Add adds a value to the set (like set.add() in TypeScript)
func (s *Set[T]) Add(value T) *Set[T] {
	if _, exists := s.data[value]; !exists {
		s.data[value] = struct{}{}
		s.size++
	}
	return s
}

// Has checks if a value exists in the set (like set.has() in TypeScript)
func (s *Set[T]) Has(value T) bool {
	_, exists := s.data[value]
	return exists
}

// Delete removes a value from the set (like set.delete() in TypeScript)
func (s *Set[T]) Delete(value T) bool {
	if _, exists := s.data[value]; exists {
		delete(s.data, value)
		s.size--
		return true
	}
	return false
}

// Clear removes all values (like set.clear() in TypeScript)
func (s *Set[T]) Clear() {
	s.data = make(map[T]struct{})
	s.size = 0
}

// Size returns the number of values (like set.size in TypeScript)
func (s *Set[T]) Size() int {
	return s.size
}

// IsEmpty checks if the set is empty
func (s *Set[T]) IsEmpty() bool {
	return s.size == 0
}

// Values returns all values (like set.values() in TypeScript)
func (s *Set[T]) Values() []T {
	values := make([]T, 0, s.size)
	for value := range s.data {
		values = append(values, value)
	}
	return values
}

// ForEach iterates over all values (like set.forEach() in TypeScript)
func (s *Set[T]) ForEach(fn func(value T, index int, set *Set[T])) {
	index := 0
	for value := range s.data {
		fn(value, index, s)
		index++
	}
}

// Filter creates a new set with values that pass the test
func (s *Set[T]) Filter(predicate func(value T) bool) *Set[T] {
	result := NewSet[T]()
	for value := range s.data {
		if predicate(value) {
			result.Add(value)
		}
	}
	return result
}

// Map transforms values and returns a new set
func SetMap[T, U comparable](s *Set[T], fn func(value T) U) *Set[U] {
	result := NewSet[U]()
	for value := range s.data {
		result.Add(fn(value))
	}
	return result
}

// Union returns a new set with all values from both sets (like set union)
func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	result := s.Clone()
	for value := range other.data {
		result.Add(value)
	}
	return result
}

// Intersection returns a new set with values present in both sets
func (s *Set[T]) Intersection(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	for value := range s.data {
		if other.Has(value) {
			result.Add(value)
		}
	}
	return result
}

// Difference returns a new set with values in this set but not in other
func (s *Set[T]) Difference(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	for value := range s.data {
		if !other.Has(value) {
			result.Add(value)
		}
	}
	return result
}

// SymmetricDifference returns values in either set but not in both
func (s *Set[T]) SymmetricDifference(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	
	// Add values from this set that are not in other
	for value := range s.data {
		if !other.Has(value) {
			result.Add(value)
		}
	}
	
	// Add values from other set that are not in this
	for value := range other.data {
		if !s.Has(value) {
			result.Add(value)
		}
	}
	
	return result
}

// IsSubsetOf checks if this set is a subset of other
func (s *Set[T]) IsSubsetOf(other *Set[T]) bool {
	if s.size > other.size {
		return false
	}
	for value := range s.data {
		if !other.Has(value) {
			return false
		}
	}
	return true
}

// IsSupersetOf checks if this set is a superset of other
func (s *Set[T]) IsSupersetOf(other *Set[T]) bool {
	return other.IsSubsetOf(s)
}

// IsDisjoint checks if this set has no common elements with other
func (s *Set[T]) IsDisjoint(other *Set[T]) bool {
	for value := range s.data {
		if other.Has(value) {
			return false
		}
	}
	return true
}

// Clone creates a shallow copy of the set
func (s *Set[T]) Clone() *Set[T] {
	result := NewSet[T]()
	for value := range s.data {
		result.Add(value)
	}
	return result
}

// ToSlice converts the set to a slice
func (s *Set[T]) ToSlice() []T {
	return s.Values()
}

// String returns string representation
func (s *Set[T]) String() string {
	return fmt.Sprintf("Set{size: %d}", s.size)
}

// WeakMap represents a weak reference map (simplified version)
// Note: Go doesn't have true weak references, so this is a simplified implementation
type WeakMap[K, V any] struct {
	data map[uintptr]V
}

// NewWeakMap creates a new WeakMap
func NewWeakMap[K, V any]() *WeakMap[K, V] {
	return &WeakMap[K, V]{
		data: make(map[uintptr]V),
	}
}

// Set sets a value for an object key
func (wm *WeakMap[K, V]) Set(key K, value V) {
	ptr := reflect.ValueOf(key).Pointer()
	wm.data[ptr] = value
}

// Get retrieves a value by object key
func (wm *WeakMap[K, V]) Get(key K) Optional[V] {
	ptr := reflect.ValueOf(key).Pointer()
	if value, exists := wm.data[ptr]; exists {
		return Some(value)
	}
	return None[V]()
}

// Has checks if an object key exists
func (wm *WeakMap[K, V]) Has(key K) bool {
	ptr := reflect.ValueOf(key).Pointer()
	_, exists := wm.data[ptr]
	return exists
}

// Delete removes an entry by object key
func (wm *WeakMap[K, V]) Delete(key K) bool {
	ptr := reflect.ValueOf(key).Pointer()
	if _, exists := wm.data[ptr]; exists {
		delete(wm.data, ptr)
		return true
	}
	return false
}

// WeakSet represents a weak reference set (simplified version)
type WeakSet[T any] struct {
	data map[uintptr]struct{}
}

// NewWeakSet creates a new WeakSet
func NewWeakSet[T any]() *WeakSet[T] {
	return &WeakSet[T]{
		data: make(map[uintptr]struct{}),
	}
}

// Add adds an object to the set
func (ws *WeakSet[T]) Add(value T) *WeakSet[T] {
	ptr := reflect.ValueOf(value).Pointer()
	ws.data[ptr] = struct{}{}
	return ws
}

// Has checks if an object is in the set
func (ws *WeakSet[T]) Has(value T) bool {
	ptr := reflect.ValueOf(value).Pointer()
	_, exists := ws.data[ptr]
	return exists
}

// Delete removes an object from the set
func (ws *WeakSet[T]) Delete(value T) bool {
	ptr := reflect.ValueOf(value).Pointer()
	if _, exists := ws.data[ptr]; exists {
		delete(ws.data, ptr)
		return true
	}
	return false
}