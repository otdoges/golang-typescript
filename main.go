package main

import (
	"fmt"
	"reflect"
	"time"
	"typescript-golang/async"
	"typescript-golang/classes"
	"typescript-golang/enums"
	"typescript-golang/testing"
	"typescript-golang/types"
	"typescript-golang/utils"
)

func main() {
	fmt.Println("🚀 TypeScript-like Go Implementation Demo")
	fmt.Println("=========================================")
	
	// Run all demos
	demoOptionalTypes()
	demoArrayUtilities()
	demoAsyncPromises()
	demoClasses()
	demoEnums()
	demoUnionTypes()
	demoStringUtilities()
	demoJSONHandling()
	demoDecorators()
	demoCollections()
	demoEvents()
	demoErrorHandling()
	demoTestingFramework()
	
	fmt.Println("\n✅ All demos completed successfully!")
}

func demoOptionalTypes() {
	fmt.Println("\n📦 Optional Types Demo")
	fmt.Println("----------------------")
	
	// Some value
	name := types.Some("John Doe")
	fmt.Printf("Name: %s\n", name.GetOrDefault("Anonymous"))
	
	// None value
	var age types.Optional[int] = types.None[int]()
	fmt.Printf("Age: %d\n", age.GetOrDefault(0))
	
	// Result types
	result := types.Ok[string, error]("Operation successful")
	if result.IsOk() {
		fmt.Printf("Result: %s\n", result.Unwrap())
	}
	
	// Error result
	errorResult := types.Err[string, error](fmt.Errorf("something went wrong"))
	fmt.Printf("Error result: %s\n", errorResult.UnwrapOr("default value"))
}

func demoArrayUtilities() {
	fmt.Println("\n🔧 Array Utilities Demo")
	fmt.Println("-----------------------")
	
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	
	// Map
	doubled := utils.Map(numbers, func(x int) int { return x * 2 })
	fmt.Printf("Doubled: %v\n", doubled)
	
	// Filter
	evens := utils.Filter(numbers, func(x int) bool { return x%2 == 0 })
	fmt.Printf("Evens: %v\n", evens)
	
	// Reduce
	sum := utils.Reduce(numbers, func(acc, x int) int { return acc + x }, 0)
	fmt.Printf("Sum: %d\n", sum)
	
	// Find
	found := utils.Find(numbers, func(x int) bool { return x > 5 })
	if found.IsSome() {
		fmt.Printf("First number > 5: %d\n", found.Get())
	}
	
	// Some and Every
	hasEven := utils.Some(numbers, func(x int) bool { return x%2 == 0 })
	allPositive := utils.Every(numbers, func(x int) bool { return x > 0 })
	fmt.Printf("Has even: %t, All positive: %t\n", hasEven, allPositive)
}

func demoAsyncPromises() {
	fmt.Println("\n⚡ Async/Promise Demo")
	fmt.Println("--------------------")
	
	// Simple promise
	promise1 := async.NewPromise(func() (string, error) {
		time.Sleep(100 * time.Millisecond)
		return "First result", nil
	})
	
	promise2 := async.NewPromise(func() (string, error) {
		time.Sleep(150 * time.Millisecond)
		return "Second result", nil
	})
	
	// Promise.all
	allPromise := async.All(promise1, promise2)
	results, err := allPromise.Await()
	if err == nil {
		fmt.Printf("All results: %v\n", results)
	}
	
	// Promise chaining
	chainedPromise := async.Then(
		async.Resolve(5),
		func(x int) string { return fmt.Sprintf("Value: %d", x*2) },
		func(err error) string { return "Error occurred" },
	)
	
	result, _ := chainedPromise.Await()
	fmt.Printf("Chained result: %s\n", result)
	
	// Sleep promise
	sleepPromise := async.Sleep(50*time.Millisecond, "Woke up!")
	sleepResult, _ := sleepPromise.Await()
	fmt.Printf("Sleep result: %s\n", sleepResult)
}

func demoClasses() {
	fmt.Println("\n🏗️  Classes Demo")
	fmt.Println("----------------")
	
	// Person class
	person := classes.NewPerson("Alice", 25)
	fmt.Printf("Person: %s\n", person.ToString())
	fmt.Printf("Is adult: %t\n", person.IsAdult())
	
	// Employee class (inheritance)
	employee := classes.NewEmployee("Bob", 30, "Software Engineer", 75000)
	fmt.Printf("Employee: %s\n", employee.ToString())
	fmt.Printf("Annual salary: $%.2f\n", employee.GetAnnualSalary())
	
	// Manager class (multi-level inheritance)
	manager := classes.NewManager("Carol", 35, "Engineering Manager", 95000, 8)
	fmt.Printf("Manager: %s\n", manager.ToString())
	fmt.Printf("Is executive: %t\n", manager.IsExecutive())
	
	// Abstract class example
	rectangle := classes.NewRectangle("blue", 10, 20)
	fmt.Printf("Rectangle: %s\n", rectangle.ToString())
	fmt.Printf("Area: %.2f, Perimeter: %.2f\n", rectangle.Area(), rectangle.Perimeter())
}

func demoEnums() {
	fmt.Println("\n📚 Enums Demo")
	fmt.Println("-------------")
	
	// Direction enum (numeric)
	direction := enums.Up
	fmt.Printf("Direction: %s (value: %v, ordinal: %d)\n", 
		direction.String(), direction.Value(), direction.Ordinal())
	
	// Color enum (string)
	color := enums.Blue
	fmt.Printf("Color: %s (value: %v, ordinal: %d)\n", 
		color.String(), color.Value(), color.Ordinal())
	
	// Status enum with custom values
	status := enums.InProgress
	fmt.Printf("Status: %s (value: %v, ordinal: %d)\n", 
		status.String(), status.Value(), status.Ordinal())
	
	// Parse enum
	if parsedColor, err := enums.ParseColor("red"); err == nil {
		fmt.Printf("Parsed color: %s\n", parsedColor.String())
	}
	
	// Log level enum
	logLevel := enums.Info
	fmt.Printf("Log level: %s, is at least warn: %t\n", 
		logLevel.String(), logLevel.IsAtLeast(enums.Warn))
}

func demoUnionTypes() {
	fmt.Println("\n🔀 Union Types Demo")
	fmt.Println("------------------")
	
	// String or Number union
	stringOrNum, _ := types.NewStringOrNumber("hello")
	if str, ok := stringOrNum.AsString(); ok {
		fmt.Printf("String value: %s\n", str)
	}
	
	numberOrStr, _ := types.NewStringOrNumber(42)
	if num, ok := numberOrStr.AsNumber(); ok {
		fmt.Printf("Number value: %.2f\n", num)
	}
	
	// Pattern matching with unions
	union := types.NewUnion("TypeScript-like")
	result := types.NewMatcher[string](union).Match(
		map[reflect.Type]func(interface{}) string{
			reflect.TypeOf(""): func(v interface{}) string {
				return fmt.Sprintf("Got string: %s", v.(string))
			},
			reflect.TypeOf(0): func(v interface{}) string {
				return fmt.Sprintf("Got number: %d", v.(int))
			},
		},
		func(v interface{}) string {
			return "Unknown type"
		},
	)
	fmt.Printf("Pattern match result: %s\n", result)
	
	// Type guards
	value := "hello world"
	fmt.Printf("IsString: %t, IsNumber: %t\n", 
		types.IsString(value), types.IsNumber(value))
}

func demoStringUtilities() {
	fmt.Println("\n🔤 String Utilities Demo")
	fmt.Println("------------------------")
	
	text := "Hello TypeScript World"
	
	// Basic string operations
	fmt.Printf("Original: %s\n", text)
	fmt.Printf("Length: %d\n", utils.Strings.Length(text))
	fmt.Printf("Upper: %s\n", utils.Strings.ToUpperCase(text))
	fmt.Printf("Substring(0, 5): %s\n", utils.Strings.Substring(text, 0, 5))
	
	// Advanced operations
	fmt.Printf("Starts with 'Hello': %t\n", utils.Strings.StartsWith(text, "Hello"))
	fmt.Printf("Includes 'Script': %t\n", utils.Strings.Includes(text, "Script"))
	fmt.Printf("Index of 'Type': %d\n", utils.Strings.IndexOf(text, "Type"))
	
	// Case conversions
	original := "hello-world_example"
	fmt.Printf("Original: %s\n", original)
	fmt.Printf("CamelCase: %s\n", utils.Strings.ToCamelCase(original))
	fmt.Printf("PascalCase: %s\n", utils.Strings.ToPascalCase(original))
	fmt.Printf("KebabCase: %s\n", utils.Strings.ToKebabCase(original))
	fmt.Printf("SnakeCase: %s\n", utils.Strings.ToSnakeCase(original))
	
	// Split and join
	words := utils.Strings.Split(text, " ")
	fmt.Printf("Words: %v\n", words)
	joined := utils.Join(words, "-")
	fmt.Printf("Joined: %s\n", joined)
}

func demoJSONHandling() {
	fmt.Println("\n📄 JSON Handling Demo")
	fmt.Println("---------------------")
	
	// Create object
	person := map[string]interface{}{
		"name": "John Doe",
		"age":  30,
		"city": "New York",
		"active": true,
	}
	
	// JSON stringify
	jsonStr, _ := utils.JSON.Stringify(person)
	fmt.Printf("JSON: %s\n", jsonStr)
	
	// JSON parse
	parsed, _ := utils.JSON.Parse(jsonStr)
	fmt.Printf("Parsed: %v\n", parsed)
	
	// Object utilities
	keys := utils.Object.Keys(person)
	values := utils.Object.Values(person)
	fmt.Printf("Keys: %v\n", keys)
	fmt.Printf("Values: %v\n", values)
	
	// Object manipulation
	picked := utils.Object.Pick(person, "name", "age")
	fmt.Printf("Picked: %v\n", picked)
	
	omitted := utils.Object.Omit(person, "active")
	fmt.Printf("Omitted: %v\n", omitted)
	
	// Merge objects
	extra := map[string]interface{}{"country": "USA", "zipcode": "10001"}
	merged := utils.Object.Merge(person, extra)
	fmt.Printf("Merged: %v\n", merged)
}

func demoDecorators() {
	fmt.Println("\n🎭 Decorators Demo")
	fmt.Println("-----------------")
	
	// Original function
	slowFunction := func(n int) (string, error) {
		time.Sleep(50 * time.Millisecond) // Simulate work
		if n < 0 {
			return "", fmt.Errorf("negative number: %d", n)
		}
		return fmt.Sprintf("Result: %d", n*n), nil
	}
	
	// Apply decorators
	decoratedFunc := utils.NewDecoratorChain[func(int) (string, error)]().
		WithLog("SlowFunction").
		WithTimer("SlowFunction").
		WithMemoize().
		Apply(slowFunction)
	
	// Test the decorated function
	fmt.Println("Calling decorated function twice with same input:")
	result1, _ := decoratedFunc(5)
	fmt.Printf("First call result: %s\n", result1)
	
	result2, _ := decoratedFunc(5) // Should hit cache
	fmt.Printf("Second call result: %s\n", result2)
	
	// Test with different input
	result3, _ := decoratedFunc(3)
	fmt.Printf("Different input result: %s\n", result3)
}

func demoCollections() {
	fmt.Println("\n🗂️  Collections Demo")
	fmt.Println("-------------------")
	
	// Map demo
	userMap := types.NewMap[string, string]()
	userMap.Set("alice", "Alice Johnson")
	userMap.Set("bob", "Bob Smith")
	userMap.Set("charlie", "Charlie Brown")
	
	fmt.Printf("Map size: %d\n", userMap.Size())
	if name := userMap.Get("alice"); name.IsSome() {
		fmt.Printf("Alice's full name: %s\n", name.Get())
	}
	
	fmt.Printf("Map keys: %v\n", userMap.Keys())
	fmt.Printf("Map values: %v\n", userMap.Values())
	
	// Map iteration
	fmt.Println("Map entries:")
	userMap.ForEach(func(value, key string, m *types.Map[string, string]) {
		fmt.Printf("  %s -> %s\n", key, value)
	})
	
	// Set demo
	numbersSet := types.NewSet[int]()
	numbersSet.Add(1).Add(2).Add(3).Add(2) // Adding duplicate
	
	fmt.Printf("\nSet size: %d\n", numbersSet.Size())
	fmt.Printf("Set contains 2: %t\n", numbersSet.Has(2))
	fmt.Printf("Set values: %v\n", numbersSet.Values())
	
	// Set operations
	otherSet := types.NewSetWithValues([]int{3, 4, 5})
	union := numbersSet.Union(otherSet)
	intersection := numbersSet.Intersection(otherSet)
	difference := numbersSet.Difference(otherSet)
	
	fmt.Printf("Union: %v\n", union.Values())
	fmt.Printf("Intersection: %v\n", intersection.Values())
	fmt.Printf("Difference: %v\n", difference.Values())
}

func demoEvents() {
	fmt.Println("\n📡 Events Demo")
	fmt.Println("--------------")
	
	// EventEmitter demo
	emitter := types.NewEventEmitter[string]()
	
	// Add listeners
	emitter.On("message", func(data string) {
		fmt.Printf("Listener 1 received: %s\n", data)
	})
	
	emitter.On("message", func(data string) {
		fmt.Printf("Listener 2 received: %s\n", data)
	})
	
	// One-time listener
	emitter.Once("startup", func(data string) {
		fmt.Printf("Startup event: %s\n", data)
	})
	
	// Emit events
	fmt.Println("Emitting 'message' event:")
	emitter.EmitSync("message", "Hello from EventEmitter!")
	
	fmt.Println("\nEmitting 'startup' event (once):")
	emitter.EmitSync("startup", "System initialized")
	emitter.EmitSync("startup", "This won't be heard") // Won't trigger once listener
	
	fmt.Printf("Listener count for 'message': %d\n", emitter.ListenerCount("message"))
	fmt.Printf("Event names: %v\n", emitter.EventNames())
	
	// Observable demo
	fmt.Println("\nObservable demo:")
	observable := types.NewObservable[int]()
	
	subscription1 := observable.Subscribe(func(value int) {
		fmt.Printf("Observer 1: %d\n", value)
	})
	
	subscription2 := observable.Subscribe(func(value int) {
		fmt.Printf("Observer 2: %d squared = %d\n", value, value*value)
	})
	
	// Emit some values
	observable.Next(5)
	observable.Next(10)
	
	// Unsubscribe one observer
	subscription1.Unsubscribe()
	fmt.Println("Observer 1 unsubscribed")
	
	observable.Next(15) // Only observer 2 will receive this
	
	subscription2.Unsubscribe()
	
	// Subject demo (acts as both observable and observer)
	fmt.Println("\nSubject demo:")
	subject := types.NewSubject[string]()
	
	subject.Subscribe(func(message string) {
		fmt.Printf("Subject subscriber: %s\n", message)
	})
	
	subject.Next("Hello from Subject!")
	subject.Next("Subjects are powerful!")
}

func demoErrorHandling() {
	fmt.Println("\n🚨 Error Handling Demo")
	fmt.Println("----------------------")
	
	// Basic enhanced error
	err1 := types.NewError("Something went wrong", types.ValidationError)
	err1.WithData("field", "username").WithData("value", "")
	
	fmt.Printf("Basic error: %s\n", err1.Error())
	fmt.Printf("Error code: %s\n", err1.Code())
	fmt.Printf("Error data: %v\n", err1.Data())
	
	// Error with cause (error chaining)
	originalErr := fmt.Errorf("network connection failed")
	err2 := types.NewErrorWithCause("Failed to fetch user data", originalErr, types.NetworkError)
	
	fmt.Printf("\nChained error: %s\n", err2.Error())
	fmt.Printf("Cause: %v\n", err2.Cause())
	
	// Try-catch pattern
	fmt.Println("\nTry-catch demo:")
	result, err := types.NewTry(func() (string, error) {
		// Simulate a function that might fail
		return "", types.NewValidationError("Invalid input")
	}).Catch(func(err *types.EnhancedError) {
		fmt.Printf("Caught error: [%s] %s\n", err.Code(), err.Message())
	}).Finally(func() {
		fmt.Println("Finally block executed")
	}).Execute()
	
	if err != nil {
		fmt.Printf("Try-catch returned error: %s\n", types.Formatter.FormatShort(err))
	} else {
		fmt.Printf("Try-catch result: %s\n", result)
	}
	
	// Assertions
	fmt.Println("\nAssertions demo:")
	if err := types.Assertions.True(2+2 == 4, "Math should work"); err != nil {
		fmt.Printf("Assertion failed: %s\n", err)
	} else {
		fmt.Println("Assertion passed: 2+2 = 4")
	}
	
	if err := types.Assertions.Equal(10, 5*2, "Multiplication check"); err != nil {
		fmt.Printf("Assertion failed: %s\n", err)
	} else {
		fmt.Println("Assertion passed: 5*2 = 10")
	}
	
	// Error boundary demo
	fmt.Println("\nError boundary demo:")
	boundary := types.NewErrorBoundary().
		OnError(func(err *types.EnhancedError) error {
			fmt.Printf("Error boundary caught: %s\n", err.Message())
			return nil // Error handled
		}).
		WithFallback(func(err *types.EnhancedError) interface{} {
			return "Fallback value"
		})
	
	result2, err := types.Wrap(boundary, func() (string, error) {
		return "", types.NewError("Simulated error", types.InternalError)
	})
	
	fmt.Printf("Error boundary result: %s, error: %v\n", result2, err)
	
	// Error formatting
	fmt.Println("\nError formatting demo:")
	complexErr := types.NewError("Complex error", types.AuthError).
		WithData("userId", 123).
		WithData("action", "login")
	
	fmt.Printf("Short format: %s\n", types.Formatter.FormatShort(complexErr))
	fmt.Printf("JSON format: %+v\n", types.Formatter.FormatJSON(complexErr))
}

func demoTestingFramework() {
	fmt.Println("\n🧪 Testing Framework Demo")
	fmt.Println("-------------------------")
	
	fmt.Println("Running comprehensive test suite...")
	testing.RunExampleTests()
}