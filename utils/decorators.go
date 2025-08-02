package utils

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
)

// cacheEntry represents a cached value with expiry time
type cacheEntry struct {
	value  []reflect.Value
	expiry time.Time
}

// Decorator represents a function decorator
type Decorator[T any] func(T) T

// MethodDecorator represents a method decorator
type MethodDecorator func(reflect.Value, []reflect.Value) []reflect.Value

// FunctionWrapper wraps a function with additional behavior
type FunctionWrapper[T any] struct {
	original T
	wrapper  func(T) T
}

// NewFunctionWrapper creates a new function wrapper
func NewFunctionWrapper[T any](original T, wrapper func(T) T) *FunctionWrapper[T] {
	return &FunctionWrapper[T]{
		original: original,
		wrapper:  wrapper,
	}
}

// Call executes the wrapped function
func (fw *FunctionWrapper[T]) Call() T {
	return fw.wrapper(fw.original)
}

// GetOriginal returns the original function
func (fw *FunctionWrapper[T]) GetOriginal() T {
	return fw.original
}

// Log decorator that logs function calls (like @log in TypeScript)
func Log[T any](name string) Decorator[T] {
	return func(fn T) T {
		// For functions, we need to use reflection
		fnValue := reflect.ValueOf(fn)
		if fnValue.Kind() != reflect.Func {
			return fn
		}

		// Create wrapper function
		wrapper := reflect.MakeFunc(fnValue.Type(), func(args []reflect.Value) []reflect.Value {
			fmt.Printf("[LOG] Calling %s with %d arguments\n", name, len(args))
			start := time.Now()
			
			results := fnValue.Call(args)
			
			duration := time.Since(start)
			fmt.Printf("[LOG] %s completed in %v\n", name, duration)
			
			return results
		})

		return wrapper.Interface().(T)
	}
}

// Timer decorator that measures execution time (like @timer in TypeScript)
func Timer[T any](name string) Decorator[T] {
	return func(fn T) T {
		fnValue := reflect.ValueOf(fn)
		if fnValue.Kind() != reflect.Func {
			return fn
		}

		wrapper := reflect.MakeFunc(fnValue.Type(), func(args []reflect.Value) []reflect.Value {
			start := time.Now()
			results := fnValue.Call(args)
			duration := time.Since(start)
			
			fmt.Printf("[TIMER] %s executed in %v\n", name, duration)
			return results
		})

		return wrapper.Interface().(T)
	}
}

// Memoize decorator for caching function results (like @memoize in TypeScript)
func Memoize[T any](fn T) T {
	fnValue := reflect.ValueOf(fn)
	if fnValue.Kind() != reflect.Func {
		return fn
	}

	cache := make(map[string][]reflect.Value)
	
	wrapper := reflect.MakeFunc(fnValue.Type(), func(args []reflect.Value) []reflect.Value {
		// Create cache key from arguments
		key := fmt.Sprintf("%v", args)
		
		// Check cache
		if cached, exists := cache[key]; exists {
			fmt.Printf("[MEMOIZE] Cache hit for key: %s\n", key)
			return cached
		}
		
		// Call original function
		results := fnValue.Call(args)
		
		// Store in cache
		cache[key] = results
		fmt.Printf("[MEMOIZE] Cached result for key: %s\n", key)
		
		return results
	})

	return wrapper.Interface().(T)
}

// Retry decorator that retries function on failure (like @retry in TypeScript)
func Retry[T any](maxAttempts int, delay time.Duration) Decorator[T] {
	return func(fn T) T {
		fnValue := reflect.ValueOf(fn)
		if fnValue.Kind() != reflect.Func {
			return fn
		}

		wrapper := reflect.MakeFunc(fnValue.Type(), func(args []reflect.Value) []reflect.Value {
			var lastResults []reflect.Value
			
			for attempt := 1; attempt <= maxAttempts; attempt++ {
				fmt.Printf("[RETRY] Attempt %d/%d\n", attempt, maxAttempts)
				
				results := fnValue.Call(args)
				lastResults = results
				
				// Check if last result is an error (assumes error is last return value)
				if len(results) > 0 {
					lastResult := results[len(results)-1]
					if lastResult.Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) && !lastResult.IsNil() {
						if attempt < maxAttempts {
							fmt.Printf("[RETRY] Attempt %d failed, retrying in %v\n", attempt, delay)
							time.Sleep(delay)
							continue
						} else {
							fmt.Printf("[RETRY] All %d attempts failed\n", maxAttempts)
							break
						}
					} else {
						fmt.Printf("[RETRY] Attempt %d succeeded\n", attempt)
						break
					}
				} else {
					break
				}
			}
			
			return lastResults
		})

		return wrapper.Interface().(T)
	}
}

// RateLimit decorator that limits function call frequency
func RateLimit[T any](callsPerSecond int) Decorator[T] {
	return func(fn T) T {
		fnValue := reflect.ValueOf(fn)
		if fnValue.Kind() != reflect.Func {
			return fn
		}

		interval := time.Second / time.Duration(callsPerSecond)
		lastCall := time.Time{}
		
		wrapper := reflect.MakeFunc(fnValue.Type(), func(args []reflect.Value) []reflect.Value {
			now := time.Now()
			if elapsed := now.Sub(lastCall); elapsed < interval {
				sleep := interval - elapsed
				fmt.Printf("[RATE_LIMIT] Rate limited, sleeping for %v\n", sleep)
				time.Sleep(sleep)
			}
			
			lastCall = time.Now()
			return fnValue.Call(args)
		})

		return wrapper.Interface().(T)
	}
}

// Validate decorator that validates function arguments
func Validate[T any](validator func([]reflect.Value) error) Decorator[T] {
	return func(fn T) T {
		fnValue := reflect.ValueOf(fn)
		if fnValue.Kind() != reflect.Func {
			return fn
		}

		wrapper := reflect.MakeFunc(fnValue.Type(), func(args []reflect.Value) []reflect.Value {
			// Validate arguments
			if err := validator(args); err != nil {
				// Return error (assumes function returns error as last value)
				fnType := fnValue.Type()
				results := make([]reflect.Value, fnType.NumOut())
				
				// Set all return values to zero values except last (error)
				for i := 0; i < fnType.NumOut()-1; i++ {
					results[i] = reflect.Zero(fnType.Out(i))
				}
				// Set error as last return value
				if fnType.NumOut() > 0 && fnType.Out(fnType.NumOut()-1).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
					results[fnType.NumOut()-1] = reflect.ValueOf(err)
				}
				
				return results
			}
			
			return fnValue.Call(args)
		})

		return wrapper.Interface().(T)
	}
}

// Deprecated decorator that logs deprecation warning
func Deprecated[T any](message string) Decorator[T] {
	return func(fn T) T {
		fnValue := reflect.ValueOf(fn)
		if fnValue.Kind() != reflect.Func {
			return fn
		}

		wrapper := reflect.MakeFunc(fnValue.Type(), func(args []reflect.Value) []reflect.Value {
			// Get function name
			pc := reflect.ValueOf(fn).Pointer()
			funcName := runtime.FuncForPC(pc).Name()
			
			fmt.Printf("[DEPRECATED] %s is deprecated: %s\n", funcName, message)
			return fnValue.Call(args)
		})

		return wrapper.Interface().(T)
	}
}

// Cache decorator with TTL (Time To Live)
func CacheWithTTL[T any](ttl time.Duration) Decorator[T] {
	return func(fn T) T {
		fnValue := reflect.ValueOf(fn)
		if fnValue.Kind() != reflect.Func {
			return fn
		}

		cache := make(map[string]cacheEntry)
		
		wrapper := reflect.MakeFunc(fnValue.Type(), func(args []reflect.Value) []reflect.Value {
			key := fmt.Sprintf("%v", args)
			now := time.Now()
			
			// Check cache and TTL
			if entry, exists := cache[key]; exists && now.Before(entry.expiry) {
				fmt.Printf("[CACHE] Cache hit for key: %s\n", key)
				return entry.value
			}
			
			// Call original function
			results := fnValue.Call(args)
			
			// Store in cache with TTL
			cache[key] = cacheEntry{
				value:  results,
				expiry: now.Add(ttl),
			}
			fmt.Printf("[CACHE] Cached result for key: %s (expires in %v)\n", key, ttl)
			
			return results
		})

		return wrapper.Interface().(T)
	}
}

// Compose combines multiple decorators
func Compose[T any](decorators ...Decorator[T]) Decorator[T] {
	return func(fn T) T {
		result := fn
		// Apply decorators in reverse order (last decorator is applied first)
		for i := len(decorators) - 1; i >= 0; i-- {
			result = decorators[i](result)
		}
		return result
	}
}

// DecoratorChain allows chaining decorators fluently
type DecoratorChain[T any] struct {
	decorators []Decorator[T]
}

// NewDecoratorChain creates a new decorator chain
func NewDecoratorChain[T any]() *DecoratorChain[T] {
	return &DecoratorChain[T]{
		decorators: make([]Decorator[T], 0),
	}
}

// Add adds a decorator to the chain
func (dc *DecoratorChain[T]) Add(decorator Decorator[T]) *DecoratorChain[T] {
	dc.decorators = append(dc.decorators, decorator)
	return dc
}

// WithLog adds log decorator
func (dc *DecoratorChain[T]) WithLog(name string) *DecoratorChain[T] {
	return dc.Add(Log[T](name))
}

// WithTimer adds timer decorator
func (dc *DecoratorChain[T]) WithTimer(name string) *DecoratorChain[T] {
	return dc.Add(Timer[T](name))
}

// WithMemoize adds memoize decorator
func (dc *DecoratorChain[T]) WithMemoize() *DecoratorChain[T] {
	return dc.Add(Memoize[T])
}

// WithRetry adds retry decorator
func (dc *DecoratorChain[T]) WithRetry(maxAttempts int, delay time.Duration) *DecoratorChain[T] {
	return dc.Add(Retry[T](maxAttempts, delay))
}

// WithRateLimit adds rate limit decorator
func (dc *DecoratorChain[T]) WithRateLimit(callsPerSecond int) *DecoratorChain[T] {
	return dc.Add(RateLimit[T](callsPerSecond))
}

// WithCache adds cache decorator with TTL
func (dc *DecoratorChain[T]) WithCache(ttl time.Duration) *DecoratorChain[T] {
	return dc.Add(CacheWithTTL[T](ttl))
}

// WithDeprecated adds deprecated decorator
func (dc *DecoratorChain[T]) WithDeprecated(message string) *DecoratorChain[T] {
	return dc.Add(Deprecated[T](message))
}

// Apply applies all decorators to the function
func (dc *DecoratorChain[T]) Apply(fn T) T {
	return Compose(dc.decorators...)(fn)
}

// Example usage struct with decorated methods
type ExampleService struct {
	name string
}

// NewExampleService creates a new example service
func NewExampleService(name string) *ExampleService {
	return &ExampleService{name: name}
}

// GetData is an example method that can be decorated
func (es *ExampleService) GetData(id int) (string, error) {
	if id <= 0 {
		return "", fmt.Errorf("invalid id: %d", id)
	}
	time.Sleep(100 * time.Millisecond) // Simulate work
	return fmt.Sprintf("data-%d from %s", id, es.name), nil
}

// CreateDecoratedGetData creates a decorated version of GetData
func (es *ExampleService) CreateDecoratedGetData() func(int) (string, error) {
	originalMethod := es.GetData
	
	decorated := NewDecoratorChain[func(int) (string, error)]().
		WithLog("GetData").
		WithTimer("GetData").
		WithRetry(3, 100*time.Millisecond).
		WithCache(5*time.Second).
		Apply(originalMethod)
	
	return decorated
}