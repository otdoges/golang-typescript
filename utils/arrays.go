package utils

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"typescript-golang/types"
)

// Map transforms each element in a slice (like Array.map() in TypeScript)
func Map[T, U any](slice []T, fn func(T) U) []U {
	result := make([]U, len(slice))
	for i, item := range slice {
		result[i] = fn(item)
	}
	return result
}

// MapWithIndex transforms each element with index (like Array.map() with index)
func MapWithIndex[T, U any](slice []T, fn func(T, int) U) []U {
	result := make([]U, len(slice))
	for i, item := range slice {
		result[i] = fn(item, i)
	}
	return result
}

// Filter creates a new slice with elements that pass the test (like Array.filter())
func Filter[T any](slice []T, predicate func(T) bool) []T {
	var result []T
	for _, item := range slice {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// FilterWithIndex filters with index access
func FilterWithIndex[T any](slice []T, predicate func(T, int) bool) []T {
	var result []T
	for i, item := range slice {
		if predicate(item, i) {
			result = append(result, item)
		}
	}
	return result
}

// Reduce reduces slice to a single value (like Array.reduce())
func Reduce[T, U any](slice []T, fn func(U, T) U, initialValue U) U {
	accumulator := initialValue
	for _, item := range slice {
		accumulator = fn(accumulator, item)
	}
	return accumulator
}

// ReduceWithIndex reduces with index access
func ReduceWithIndex[T, U any](slice []T, fn func(U, T, int) U, initialValue U) U {
	accumulator := initialValue
	for i, item := range slice {
		accumulator = fn(accumulator, item, i)
	}
	return accumulator
}

// ForEach executes a function for each element (like Array.forEach())
func ForEach[T any](slice []T, fn func(T)) {
	for _, item := range slice {
		fn(item)
	}
}

// ForEachWithIndex executes function with index access
func ForEachWithIndex[T any](slice []T, fn func(T, int)) {
	for i, item := range slice {
		fn(item, i)
	}
}

// Find returns the first element that matches predicate (like Array.find())
func Find[T any](slice []T, predicate func(T) bool) types.Optional[T] {
	for _, item := range slice {
		if predicate(item) {
			return types.Some(item)
		}
	}
	return types.None[T]()
}

// FindWithIndex finds with index access
func FindWithIndex[T any](slice []T, predicate func(T, int) bool) types.Optional[T] {
	for i, item := range slice {
		if predicate(item, i) {
			return types.Some(item)
		}
	}
	return types.None[T]()
}

// FindIndex returns the index of first matching element (like Array.findIndex())
func FindIndex[T any](slice []T, predicate func(T) bool) int {
	for i, item := range slice {
		if predicate(item) {
			return i
		}
	}
	return -1
}

// Some checks if at least one element passes the test (like Array.some())
func Some[T any](slice []T, predicate func(T) bool) bool {
	for _, item := range slice {
		if predicate(item) {
			return true
		}
	}
	return false
}

// Every checks if all elements pass the test (like Array.every())
func Every[T any](slice []T, predicate func(T) bool) bool {
	for _, item := range slice {
		if !predicate(item) {
			return false
		}
	}
	return true
}

// Includes checks if slice contains a specific value (like Array.includes())
func Includes[T comparable](slice []T, value T) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

// IndexOf returns the index of first occurrence (like Array.indexOf())
func IndexOf[T comparable](slice []T, value T) int {
	for i, item := range slice {
		if item == value {
			return i
		}
	}
	return -1
}

// LastIndexOf returns the index of last occurrence (like Array.lastIndexOf())
func LastIndexOf[T comparable](slice []T, value T) int {
	for i := len(slice) - 1; i >= 0; i-- {
		if slice[i] == value {
			return i
		}
	}
	return -1
}

// Concat concatenates multiple slices (like Array.concat())
func Concat[T any](slices ...[]T) []T {
	var result []T
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return result
}

// Slice returns a portion of slice (like Array.slice())
func Slice[T any](slice []T, start, end int) []T {
	length := len(slice)
	
	// Handle negative indices
	if start < 0 {
		start = length + start
	}
	if end < 0 {
		end = length + end
	}
	
	// Clamp to bounds
	if start < 0 {
		start = 0
	}
	if end > length {
		end = length
	}
	if start > end {
		start = end
	}
	
	return slice[start:end]
}

// Splice modifies slice by removing/adding elements (like Array.splice())
func Splice[T any](slice []T, start int, deleteCount int, items ...T) []T {
	length := len(slice)
	
	// Handle negative start
	if start < 0 {
		start = length + start
		if start < 0 {
			start = 0
		}
	} else if start > length {
		start = length
	}
	
	// Handle deleteCount
	if deleteCount < 0 {
		deleteCount = 0
	}
	if start+deleteCount > length {
		deleteCount = length - start
	}
	
	// Create result
	result := make([]T, 0, length-deleteCount+len(items))
	result = append(result, slice[:start]...)
	result = append(result, items...)
	result = append(result, slice[start+deleteCount:]...)
	
	return result
}

// Push adds elements to end of slice (like Array.push())
func Push[T any](slice []T, items ...T) []T {
	return append(slice, items...)
}

// Pop removes and returns last element (like Array.pop())
func Pop[T any](slice []T) ([]T, types.Optional[T]) {
	if len(slice) == 0 {
		return slice, types.None[T]()
	}
	last := slice[len(slice)-1]
	return slice[:len(slice)-1], types.Some(last)
}

// Shift removes and returns first element (like Array.shift())
func Shift[T any](slice []T) ([]T, types.Optional[T]) {
	if len(slice) == 0 {
		return slice, types.None[T]()
	}
	first := slice[0]
	return slice[1:], types.Some(first)
}

// Unshift adds elements to beginning of slice (like Array.unshift())
func Unshift[T any](slice []T, items ...T) []T {
	return append(items, slice...)
}

// Reverse reverses the slice in place (like Array.reverse())
func Reverse[T any](slice []T) []T {
	length := len(slice)
	for i := 0; i < length/2; i++ {
		slice[i], slice[length-1-i] = slice[length-1-i], slice[i]
	}
	return slice
}

// Reversed returns a new reversed slice (non-mutating version)
func Reversed[T any](slice []T) []T {
	result := make([]T, len(slice))
	for i, item := range slice {
		result[len(slice)-1-i] = item
	}
	return result
}

// Sort sorts the slice using a comparison function (like Array.sort())
func Sort[T any](slice []T, compare func(T, T) int) []T {
	sort.Slice(slice, func(i, j int) bool {
		return compare(slice[i], slice[j]) < 0
	})
	return slice
}

// Sorted returns a new sorted slice (non-mutating version)
func Sorted[T any](slice []T, compare func(T, T) int) []T {
	result := make([]T, len(slice))
	copy(result, slice)
	return Sort(result, compare)
}

// SortBy sorts by a key function
func SortBy[T any, K types.Ordered](slice []T, keyFn func(T) K) []T {
	return Sort(slice, func(a, b T) int {
		keyA, keyB := keyFn(a), keyFn(b)
		if keyA < keyB {
			return -1
		}
		if keyA > keyB {
			return 1
		}
		return 0
	})
}

// Join converts slice to string with separator (like Array.join())
func Join[T any](slice []T, separator string) string {
	if len(slice) == 0 {
		return ""
	}
	
	strSlice := make([]string, len(slice))
	for i, item := range slice {
		strSlice[i] = fmt.Sprintf("%v", item)
	}
	
	return strings.Join(strSlice, separator)
}

// Flatten flattens one level of nesting (like Array.flat())
func Flatten[T any](slice [][]T) []T {
	var result []T
	for _, subSlice := range slice {
		result = append(result, subSlice...)
	}
	return result
}

// FlatMap maps and flattens (like Array.flatMap())
func FlatMap[T, U any](slice []T, fn func(T) []U) []U {
	var result []U
	for _, item := range slice {
		result = append(result, fn(item)...)
	}
	return result
}

// Unique returns slice with unique elements
func Unique[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	var result []T
	
	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	
	return result
}

// UniqueBy returns slice with unique elements based on key function
func UniqueBy[T any, K comparable](slice []T, keyFn func(T) K) []T {
	seen := make(map[K]bool)
	var result []T
	
	for _, item := range slice {
		key := keyFn(item)
		if !seen[key] {
			seen[key] = true
			result = append(result, item)
		}
	}
	
	return result
}

// GroupBy groups elements by key function
func GroupBy[T any, K comparable](slice []T, keyFn func(T) K) map[K][]T {
	groups := make(map[K][]T)
	
	for _, item := range slice {
		key := keyFn(item)
		groups[key] = append(groups[key], item)
	}
	
	return groups
}

// Partition splits slice into two based on predicate
func Partition[T any](slice []T, predicate func(T) bool) ([]T, []T) {
	var truthy, falsy []T
	
	for _, item := range slice {
		if predicate(item) {
			truthy = append(truthy, item)
		} else {
			falsy = append(falsy, item)
		}
	}
	
	return truthy, falsy
}

// Chunk splits slice into chunks of specified size
func Chunk[T any](slice []T, size int) [][]T {
	if size <= 0 {
		return [][]T{}
	}
	
	var chunks [][]T
	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	
	return chunks
}

// Take returns first n elements
func Take[T any](slice []T, n int) []T {
	if n <= 0 {
		return []T{}
	}
	if n >= len(slice) {
		return slice
	}
	return slice[:n]
}

// TakeWhile takes elements while predicate is true
func TakeWhile[T any](slice []T, predicate func(T) bool) []T {
	var result []T
	for _, item := range slice {
		if !predicate(item) {
			break
		}
		result = append(result, item)
	}
	return result
}

// Drop returns slice without first n elements
func Drop[T any](slice []T, n int) []T {
	if n <= 0 {
		return slice
	}
	if n >= len(slice) {
		return []T{}
	}
	return slice[n:]
}

// DropWhile drops elements while predicate is true
func DropWhile[T any](slice []T, predicate func(T) bool) []T {
	for i, item := range slice {
		if !predicate(item) {
			return slice[i:]
		}
	}
	return []T{}
}

// Zip combines two slices into pairs
func Zip[T, U any](slice1 []T, slice2 []U) []types.Tuple2[T, U] {
	minLen := len(slice1)
	if len(slice2) < minLen {
		minLen = len(slice2)
	}
	
	result := make([]types.Tuple2[T, U], minLen)
	for i := 0; i < minLen; i++ {
		result[i] = types.NewTuple2(slice1[i], slice2[i])
	}
	
	return result
}

// ZipWith combines two slices using a function
func ZipWith[T, U, V any](slice1 []T, slice2 []U, fn func(T, U) V) []V {
	minLen := len(slice1)
	if len(slice2) < minLen {
		minLen = len(slice2)
	}
	
	result := make([]V, minLen)
	for i := 0; i < minLen; i++ {
		result[i] = fn(slice1[i], slice2[i])
	}
	
	return result
}

// IsEmpty checks if slice is empty
func IsEmpty[T any](slice []T) bool {
	return len(slice) == 0
}

// Length returns the length of slice (for consistency with TypeScript)
func Length[T any](slice []T) int {
	return len(slice)
}

// First returns the first element
func First[T any](slice []T) types.Optional[T] {
	if IsEmpty(slice) {
		return types.None[T]()
	}
	return types.Some(slice[0])
}

// Last returns the last element
func Last[T any](slice []T) types.Optional[T] {
	if IsEmpty(slice) {
		return types.None[T]()
	}
	return types.Some(slice[len(slice)-1])
}

// At returns element at index (supports negative indices like TypeScript)
func At[T any](slice []T, index int) types.Optional[T] {
	length := len(slice)
	if index < 0 {
		index = length + index
	}
	
	if index < 0 || index >= length {
		return types.None[T]()
	}
	
	return types.Some(slice[index])
}

// Count counts elements matching predicate
func Count[T any](slice []T, predicate func(T) bool) int {
	count := 0
	for _, item := range slice {
		if predicate(item) {
			count++
		}
	}
	return count
}

// Min returns minimum element (for ordered types)
func Min[T types.Ordered](slice []T) types.Optional[T] {
	if IsEmpty(slice) {
		return types.None[T]()
	}
	
	min := slice[0]
	for _, item := range slice[1:] {
		if item < min {
			min = item
		}
	}
	
	return types.Some(min)
}

// Max returns maximum element (for ordered types)
func Max[T types.Ordered](slice []T) types.Optional[T] {
	if IsEmpty(slice) {
		return types.None[T]()
	}
	
	max := slice[0]
	for _, item := range slice[1:] {
		if item > max {
			max = item
		}
	}
	
	return types.Some(max)
}

// Sum calculates sum of numeric slice
func Sum[T types.Ordered](slice []T) T {
	var sum T
	for _, item := range slice {
		sum += item
	}
	return sum
}

// Average calculates average of numeric slice
func Average[T types.Ordered](slice []T) types.Optional[float64] {
	if IsEmpty(slice) {
		return types.None[float64]()
	}
	
	sum := Sum(slice)
	return types.Some(float64(reflect.ValueOf(sum).Convert(reflect.TypeOf(float64(0))).Float()) / float64(len(slice)))
}