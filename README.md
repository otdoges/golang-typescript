# TypeScript-like Go Implementation

A comprehensive Go implementation that provides TypeScript-like language features and patterns, making it easier for TypeScript developers to work with Go while leveraging Go's performance and concurrency features.

## 🚀 Features

### Core Language Features
- **Generic Types**: TypeScript-like generic type system using Go generics
- **Optional Types**: `Optional<T>` similar to TypeScript's `T | undefined`
- **Union Types**: Type-safe union types with pattern matching
- **Result Types**: Error handling similar to Rust/TypeScript's Result pattern
- **Structural Typing**: Interface-based structural typing system

### Collections & Data Structures
- **Map<K,V>**: TypeScript-like Map with full API compatibility
- **Set<T>**: TypeScript-like Set with union, intersection, difference operations
- **WeakMap & WeakSet**: Simplified weak reference collections
- **Tuple Types**: Strongly-typed tuple implementations (Tuple2, Tuple3)

### Class-like Structures
- **Classes**: Class-like structs with constructors, methods, and inheritance
- **Abstract Classes**: Abstract base classes with virtual methods
- **Method Overriding**: Support for method overriding in derived classes
- **Access Patterns**: Public/private-like access patterns

### Async Programming
- **Promises**: TypeScript-like Promise implementation using goroutines
- **Async/Await**: Promise chaining and async patterns
- **Promise Utilities**: `Promise.all()`, `Promise.race()`, `Promise.any()`, etc.
- **Timeout Support**: Promise timeout and cancellation

### Event System
- **EventEmitter**: TypeScript/Node.js-like EventEmitter with full API
- **Observable Pattern**: RxJS-like Observable and Subject implementations
- **Event Bus**: Global event communication system
- **Promise-based Events**: Event waiting with timeout support

### Error Handling
- **Enhanced Errors**: Rich error objects with stack traces and error chaining
- **Try-Catch Pattern**: TypeScript-like try-catch-finally blocks
- **Error Boundaries**: React-like error boundary pattern
- **Assertions**: Comprehensive assertion framework
- **Error Formatting**: Multiple error output formats (short, detailed, JSON)

### Utility Functions
- **Array Methods**: Complete set of TypeScript array methods (`map`, `filter`, `reduce`, etc.)
- **String Methods**: TypeScript-like string manipulation methods
- **Object Utilities**: Object manipulation similar to TypeScript's Object methods
- **JSON Handling**: TypeScript-like JSON stringify/parse functionality

### Testing Framework
- **Jest/Mocha Style**: Familiar `describe`, `it`, `beforeEach`, `afterEach` syntax
- **Rich Assertions**: Comprehensive `expect` API with matchers
- **Test Suites**: Nested test organization and reporting
- **Lifecycle Hooks**: Full setup/teardown lifecycle support
- **Test Reporting**: Detailed console reporting with timing and status

### Advanced Features
- **Enums**: Numeric and string enums with TypeScript-like syntax
- **Decorators**: Function decorators for logging, caching, timing, etc.
- **Type Guards**: Runtime type checking and narrowing
- **Performance Optimized**: Zero-cost abstractions where possible

## 📦 Installation

```bash
# Clone the repository
git clone <repository-url>
cd typescript-golang

# Install dependencies
make deps

# Build the project
make build
```

## 🎯 Quick Start

```go
package main

import (
    "fmt"
    "typescript-golang/types"
    "typescript-golang/utils"
    "typescript-golang/async"
    "typescript-golang/classes"
)

func main() {
    // Optional types
    name := types.Some("John")
    fmt.Println(name.GetOrDefault("Anonymous")) // "John"

    // Array utilities
    numbers := []int{1, 2, 3, 4, 5}
    doubled := utils.Map(numbers, func(x int) int { return x * 2 })
    fmt.Println(doubled) // [2, 4, 6, 8, 10]

    // Promises
    promise := async.NewPromise(func() (string, error) {
        return "Hello from async!", nil
    })
    result, _ := promise.Await()
    fmt.Println(result) // "Hello from async!"

    // Classes
    person := classes.NewPerson("Alice", 25)
    fmt.Println(person.ToString()) // "Person{Name: Alice, Age: 25}"
}
```

## 📚 Documentation

### Types Package

The `types` package provides TypeScript-like type system features:

```go
// Optional types
var name types.Optional[string] = types.Some("John")
if name.IsSome() {
    fmt.Println("Name:", name.Get())
}

// Result types
result := types.Ok[string, error]("success")
if result.IsOk() {
    fmt.Println("Result:", result.Unwrap())
}

// Union types
union := types.NewStringOrNumber("hello")
if str, ok := union.AsString(); ok {
    fmt.Println("String value:", str)
}
```

### Utils Package

Comprehensive utility functions for arrays, strings, and objects:

```go
// Array utilities
numbers := []int{1, 2, 3, 4, 5}
evens := utils.Filter(numbers, func(x int) bool { return x%2 == 0 })
sum := utils.Reduce(numbers, func(acc, x int) int { return acc + x }, 0)

// String utilities
text := "Hello World"
upper := utils.Strings.ToUpperCase(text)
substr := utils.Strings.Substring(text, 0, 5)

// Object utilities
obj := map[string]interface{}{"name": "John", "age": 30}
keys := utils.Object.Keys(obj)
values := utils.Object.Values(obj)
```

### Async Package

Promise-based asynchronous programming:

```go
// Create a promise
promise := async.NewPromise(func() (int, error) {
    time.Sleep(1 * time.Second)
    return 42, nil
})

// Chain promises
result := async.Then(promise, 
    func(x int) string { return fmt.Sprintf("Result: %d", x) },
    func(err error) string { return "Error occurred" },
)

// Promise utilities
promises := []*async.Promise[int]{promise1, promise2, promise3}
allResults := async.All(promises...)
```

### Classes Package

Object-oriented programming patterns:

```go
// Base class
person := classes.NewPerson("John", 25)
person.SetAge(26)
fmt.Println(person.IsAdult()) // true

// Inheritance
employee := classes.NewEmployee("Jane", 30, "Developer", 75000)
fmt.Println(employee.GetAnnualSalary()) // 900000

// Abstract classes and method overriding
rectangle := classes.NewRectangle("red", 10, 20)
fmt.Println(rectangle.Area()) // 200
fmt.Println(rectangle.Perimeter()) // 60
```

### Enums Package

TypeScript-like enum implementations:

```go
// Numeric enums
direction := enums.Up
fmt.Println(direction.String()) // "Up"
fmt.Println(direction.Value())  // 0

// String enums
color := enums.Red
fmt.Println(color.String()) // "red"
fmt.Println(color.Name())   // "Red"

// Parse enums
if parsed, err := enums.ParseColor("blue"); err == nil {
    fmt.Println("Parsed color:", parsed)
}
```

## 🛠️ Development

### Available Make Targets

```bash
make help          # Show all available targets
make build         # Build the project
make run           # Run the project
make test          # Run tests
make fmt           # Format code
make lint          # Lint code
make dev           # Development workflow (fmt, vet, test, build)
make clean         # Clean build artifacts
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run benchmarks
make bench
```

### Code Quality

```bash
# Run all quality checks
make quality

# Individual checks
make fmt           # Format code
make vet           # Run go vet
make lint          # Run golangci-lint
```

## 📋 Project Structure

```
typescript-golang/
├── main.go              # Main entry point
├── go.mod              # Go module definition
├── Makefile            # Build and development scripts
├── README.md           # This file
├── types/              # Type system implementations
│   ├── interfaces.go   # Structural typing and interfaces
│   ├── generics.go     # Generic types and utilities
│   └── unions.go       # Union types and type guards
├── utils/              # Utility functions
│   ├── arrays.go       # Array/slice utilities
│   ├── strings.go      # String manipulation utilities
│   ├── json.go         # JSON and object utilities
│   └── decorators.go   # Function decorators
├── async/              # Asynchronous programming
│   └── promise.go      # Promise implementation
├── classes/            # Class-like structures
│   └── base.go         # Base classes and inheritance
└── enums/              # Enum implementations
    └── enums.go        # Numeric and string enums
```

## 🎨 Examples

### Advanced Promise Usage

```go
// Promise chaining
promise := async.NewPromise(func() (int, error) {
    return 10, nil
})

chained := async.ThenPromise(promise, func(x int) *async.Promise[string] {
    return async.NewPromise(func() (string, error) {
        return fmt.Sprintf("Value: %d", x*2), nil
    })
})

result, _ := chained.Await()
fmt.Println(result) // "Value: 20"
```

### Decorator Usage

```go
// Function with multiple decorators
decoratedFunc := utils.NewDecoratorChain[func(int) (string, error)]().
    WithLog("MyFunction").
    WithTimer("MyFunction").
    WithRetry(3, 100*time.Millisecond).
    WithCache(5*time.Second).
    Apply(myFunction)

result, err := decoratedFunc(42)
```

### Complex Type Operations

```go
// Working with union types and pattern matching
union := types.NewStringOrNumber(42)
result := types.NewMatcher[string](union).Match(
    map[reflect.Type]func(interface{}) string{
        reflect.TypeOf(""): func(v interface{}) string {
            return "Got string: " + v.(string)
        },
        reflect.TypeOf(0): func(v interface{}) string {
            return fmt.Sprintf("Got number: %d", v.(int))
        },
    },
    func(v interface{}) string {
        return "Unknown type"
    },
)
```

## 🔧 Configuration

### Environment Variables

- `GO_ENV`: Set to `development` or `production`
- `LOG_LEVEL`: Set logging level (`debug`, `info`, `warn`, `error`)

### Build Tags

```bash
# Build with debug information
go build -tags debug

# Build for production
go build -tags production -ldflags="-w -s"
```

## 📊 Performance

This implementation is designed to provide TypeScript-like APIs while maintaining Go's performance characteristics:

- **Zero-cost abstractions**: Most utilities compile to efficient Go code
- **Minimal allocations**: Careful memory management in hot paths
- **Concurrent-safe**: Thread-safe implementations where appropriate
- **Benchmarked**: All utilities include performance benchmarks

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/new-feature`
3. Make your changes and add tests
4. Run quality checks: `make quality`
5. Commit your changes: `git commit -am 'Add new feature'`
6. Push to the branch: `git push origin feature/new-feature`
7. Create a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by TypeScript's elegant type system and developer experience
- Built on Go's powerful concurrency primitives and performance
- Designed for developers transitioning from TypeScript to Go

## 📈 Roadmap

- [x] **Collections**: Map, Set, WeakMap, WeakSet implementations ✅
- [x] **Event System**: EventEmitter, Observable, Subject patterns ✅
- [x] **Enhanced Error Handling**: Stack traces, error chaining, try-catch ✅
- [x] **Testing Framework**: Jest/Mocha-style testing with assertions ✅
- [x] **TypeScript Comparison**: Side-by-side implementation comparison ✅
- [ ] **Namespace System**: TypeScript-like module organization
- [ ] **Advanced Async**: Async generators, streams, and iterators
- [ ] **Performance Benchmarks**: Comprehensive performance analysis
- [ ] **Example Projects**: Real-world usage demonstrations
- [ ] **VS Code Extension**: Enhanced developer experience
- [ ] **Documentation Site**: Interactive documentation and examples

---

**Made with ❤️ for TypeScript developers learning Go**