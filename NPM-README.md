# TypeScript-Golang NPM Package

A comprehensive Node.js package that provides TypeScript developers with the tools to work with Go using familiar TypeScript-like patterns.

## 🚀 Quick Start

```bash
# Install globally for CLI usage
npm install -g typescript-golang

# Or install locally in your project
npm install typescript-golang
```

## 📦 What's Included

- **🛠️ CLI Tools**: Project scaffolding and code generation
- **📚 Go Library**: Complete TypeScript-like Go implementation
- **🎨 Project Templates**: Ready-to-use project templates
- **🔄 Code Generator**: TypeScript to Go conversion utilities
- **📖 Documentation**: Comprehensive guides and examples

## 🎯 CLI Usage

### Create a New Project

```bash
# Create a basic Go project with TypeScript-like patterns
ts-go init my-awesome-project

# Use a specific template
ts-go init my-web-api --template web-api
ts-go init my-cli-tool --template cli-app
```

### Generate Go Code from TypeScript

```bash
# Convert TypeScript interfaces to Go structs
ts-go generate types.ts --output types.go

# Convert entire TypeScript files
ts-go generate src/models.ts --output models/models.go
```

### Available Templates

```bash
# List all available templates
ts-go templates

# Show usage examples
ts-go examples

# Run the demo
ts-go demo
```

## 📋 Available Templates

### `basic-project`
A basic Go project showcasing all TypeScript-like features:
- Optional types
- Array utilities
- Async/Promise patterns
- Classes and inheritance
- Collections (Map/Set)
- Event system
- Error handling
- Testing framework

### `web-api`
Web API server template with:
- HTTP server with middleware
- TypeScript-like error handling
- Promise-based async operations
- JSON utilities
- Event-driven architecture
- API testing framework

### `cli-app`
Command-line application template with:
- CLI argument parsing
- Interactive prompts
- File system operations
- Progress indicators
- Configuration management

## 🔧 Programmatic Usage

```typescript
import { createProject, generateGoCode, GoUtils } from 'typescript-golang';

// Create a new project programmatically
await createProject({
  projectName: 'my-project',
  template: 'web-api',
  outputDir: './my-project'
});

// Generate Go code from TypeScript
await generateGoCode('types.ts', 'types.go');

// Utility functions
const goType = GoUtils.convertType('string[]'); // Returns: []string
const isValid = await GoUtils.validateSyntax(goCode);
```

## 🎯 TypeScript to Go Conversion Examples

### Interfaces → Structs

**TypeScript:**
```typescript
interface User {
  id: number;
  name: string;
  email?: string;
}
```

**Generated Go:**
```go
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email,omitempty"`
}
```

### Classes → Structs with Methods

**TypeScript:**
```typescript
class Person {
  constructor(public name: string, public age: number) {}
  
  isAdult(): boolean {
    return this.age >= 18;
  }
}
```

**Generated Go:**
```go
type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func NewPerson(name string, age int) *Person {
    return &Person{
        Name: name,
        Age:  age,
    }
}

func (p *Person) IsAdult() bool {
    return p.Age >= 18
}
```

### Enums → Constants

**TypeScript:**
```typescript
enum Color {
  Red = "red",
  Green = "green",
  Blue = "blue"
}
```

**Generated Go:**
```go
type Color string

const (
    Red   Color = "red"
    Green Color = "green"
    Blue  Color = "blue"
)

func (c Color) String() string {
    switch c {
    case Red:
        return "Red"
    case Green:
        return "Green"
    case Blue:
        return "Blue"
    default:
        return "Unknown"
    }
}
```

## 🏗️ Generated Project Structure

When you create a new project, you get:

```
my-project/
├── main.go                 # Entry point with examples
├── go.mod                  # Go module definition
├── types/                  # TypeScript-like type system
│   ├── generics.go        # Optional, Result, Tuple types
│   ├── unions.go          # Union types and pattern matching
│   ├── collections.go     # Map, Set implementations
│   ├── events.go          # EventEmitter, Observable
│   └── errors.go          # Enhanced error handling
├── utils/                  # Utility functions
│   ├── arrays.go          # Array methods (map, filter, reduce)
│   ├── strings.go         # String utilities
│   ├── json.go            # JSON and object utilities
│   └── decorators.go      # Function decorators
├── async/                  # Async programming
│   └── promise.go         # Promise implementation
├── classes/                # Class-like structures
│   └── base.go            # Base classes and inheritance
├── enums/                  # Enum implementations
│   └── enums.go           # Numeric and string enums
└── testing/                # Testing framework
    └── framework.go       # Jest/Mocha-style testing
```

## 🎨 Features Showcase

### Optional Types
```go
name := types.Some("John Doe")
age := types.None[int]()

fmt.Println(name.GetOrDefault("Anonymous"))  // John Doe
fmt.Println(age.GetOrDefault(0))             // 0
```

### Array Utilities
```go
numbers := []int{1, 2, 3, 4, 5}
doubled := utils.Map(numbers, func(x int) int { return x * 2 })
evens := utils.Filter(numbers, func(x int) bool { return x%2 == 0 })
sum := utils.Reduce(numbers, func(acc, x int) int { return acc + x }, 0)
```

### Promises
```go
promise := async.NewPromise(func() (string, error) {
    time.Sleep(1 * time.Second)
    return "Hello World", nil
})

result, err := promise.Await()
```

### Collections
```go
userMap := types.NewMap[string, string]()
userMap.Set("alice", "Alice Johnson")

numbersSet := types.NewSet[int]()
numbersSet.Add(1).Add(2).Add(3)
```

### Events
```go
emitter := types.NewEventEmitter[string]()
emitter.On("message", func(data string) {
    fmt.Printf("Received: %s\n", data)
})
emitter.Emit("message", "Hello!")
```

## 🔧 Advanced Usage

### Custom Code Generation

```typescript
import { TypeScriptToGoGenerator } from 'typescript-golang';

const generator = new TypeScriptToGoGenerator({
  preserveComments: true,
  generateTests: true,
  outputPackage: 'models'
});

const result = await generator.convertCode(typescriptCode);
console.log(result.goCode);
```

### Project Analysis

```typescript
import { CodeAnalysis } from 'typescript-golang';

const interfaces = CodeAnalysis.extractInterfaces(tsCode);
const classes = CodeAnalysis.extractClasses(tsCode);
const enums = CodeAnalysis.extractEnums(tsCode);
```

### Build Integration

```typescript
import { BuildUtils } from 'typescript-golang';

// Build the Go project
const buildResult = await BuildUtils.goBuild('./my-project');

// Run tests
const testResult = await BuildUtils.goTest('./my-project');

// Format code
const formatResult = await BuildUtils.goFormat('./my-project');
```

## 🛠️ Development Workflow

1. **Create Project**: `ts-go init my-project`
2. **Develop**: Write Go code using TypeScript-like patterns
3. **Generate**: `ts-go generate` for TypeScript→Go conversion
4. **Build**: `go build` or use `BuildUtils.goBuild()`
5. **Test**: Use the built-in testing framework
6. **Deploy**: Standard Go deployment

## 📚 Documentation

- [Complete API Reference](https://github.com/typescript-golang/typescript-golang)
- [TypeScript vs Go Comparison](https://github.com/typescript-golang/typescript-golang/blob/main/COMPARISON.md)
- [Migration Guide](https://github.com/typescript-golang/typescript-golang/blob/main/MIGRATION.md)
- [Examples Repository](https://github.com/typescript-golang/examples)

## 🤝 Contributing

We welcome contributions! See our [Contributing Guide](https://github.com/typescript-golang/typescript-golang/blob/main/CONTRIBUTING.md).

## 📄 License

MIT License - see [LICENSE](https://github.com/typescript-golang/typescript-golang/blob/main/LICENSE) for details.

## 🙏 Acknowledgments

- Inspired by TypeScript's elegant type system
- Built on Go's powerful performance and concurrency
- Designed for smooth language transitions

---

**Made with ❤️ for TypeScript developers learning Go**