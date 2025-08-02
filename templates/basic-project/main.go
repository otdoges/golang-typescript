package main

import (
	"fmt"
	"time"
	"PROJECT_NAME/async"
	"PROJECT_NAME/classes"
	"PROJECT_NAME/enums"
	"PROJECT_NAME/types"
	"PROJECT_NAME/utils"
)

func main() {
	fmt.Println("🚀 Welcome to PROJECT_NAME!")
	fmt.Println("TypeScript-like Go implementation ready to use!")
	fmt.Println()

	// Optional Types Demo
	fmt.Println("📦 Optional Types:")
	name := types.Some("TypeScript Developer")
	age := types.None[int]()
	fmt.Printf("  Name: %s\n", name.GetOrDefault("Anonymous"))
	fmt.Printf("  Age: %d\n", age.GetOrDefault(25))
	fmt.Println()

	// Array Utilities Demo
	fmt.Println("🔧 Array Utilities:")
	numbers := []int{1, 2, 3, 4, 5}
	doubled := utils.Map(numbers, func(x int) int { return x * 2 })
	evens := utils.Filter(numbers, func(x int) bool { return x%2 == 0 })
	sum := utils.Reduce(numbers, func(acc, x int) int { return acc + x }, 0)
	fmt.Printf("  Original: %v\n", numbers)
	fmt.Printf("  Doubled: %v\n", doubled)
	fmt.Printf("  Evens: %v\n", evens)
	fmt.Printf("  Sum: %d\n", sum)
	fmt.Println()

	// Promises Demo
	fmt.Println("⚡ Async/Promise:")
	promise := async.NewPromise(func() (string, error) {
		time.Sleep(100 * time.Millisecond)
		return "Hello from Promise!", nil
	})
	
	result, err := promise.Await()
	if err == nil {
		fmt.Printf("  Promise result: %s\n", result)
	}
	fmt.Println()

	// Classes Demo
	fmt.Println("🏗️  Classes:")
	person := NewPerson("Alice", 30)
	fmt.Printf("  Person: %s\n", person.ToString())
	fmt.Printf("  Is Adult: %t\n", person.IsAdult())
	fmt.Println()

	// Collections Demo
	fmt.Println("🗂️  Collections:")
	userMap := types.NewMap[string, string]()
	userMap.Set("alice", "Alice Johnson")
	userMap.Set("bob", "Bob Smith")
	fmt.Printf("  Map size: %d\n", userMap.Size())
	if alice := userMap.Get("alice"); alice.IsSome() {
		fmt.Printf("  Alice: %s\n", alice.Get())
	}

	numbersSet := types.NewSet[int]()
	numbersSet.Add(1).Add(2).Add(3).Add(2) // Duplicate ignored
	fmt.Printf("  Set size: %d\n", numbersSet.Size())
	fmt.Printf("  Set contains 2: %t\n", numbersSet.Has(2))
	fmt.Println()

	// Events Demo
	fmt.Println("📡 Events:")
	emitter := types.NewEventEmitter[string]()
	emitter.On("message", func(data string) {
		fmt.Printf("  Received: %s\n", data)
	})
	emitter.Emit("message", "Hello Events!")
	fmt.Println()

	// Enums Demo
	fmt.Println("📚 Enums:")
	direction := enums.Up
	color := enums.Blue
	fmt.Printf("  Direction: %s (value: %d)\n", direction.String(), int(direction))
	fmt.Printf("  Color: %s\n", string(color))
	fmt.Println()

	fmt.Println("✅ All features working! Ready to build something amazing! 🎉")
}

// Example class implementation
type Person struct {
	*classes.BaseClass
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func NewPerson(name string, age int) *Person {
	return &Person{
		BaseClass: classes.NewBaseClass("Person"),
		Name:      name,
		Age:       age,
	}
}

func (p *Person) ToString() string {
	return fmt.Sprintf("Person{Name: %s, Age: %d}", p.Name, p.Age)
}

func (p *Person) IsAdult() bool {
	return p.Age >= 18
}