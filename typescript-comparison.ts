// TypeScript Implementation Comparison
// This file demonstrates the same functionality we've implemented in Go
// to validate how closely our Go implementation matches TypeScript patterns

// ============================================================================
// 1. OPTIONAL TYPES (TypeScript's nullable types and our Optional<T>)
// ============================================================================

type Optional<T> = T | undefined | null;

function someValue<T>(value: T): T {
    return value;
}

function getValue<T>(opt: Optional<T>, defaultValue: T): T {
    return opt ?? defaultValue;
}

// Usage
const userName: Optional<string> = "John Doe";
const userAge: Optional<number> = undefined;

console.log("📦 Optional Types Demo");
console.log("Name:", getValue(userName, "Anonymous"));
console.log("Age:", getValue(userAge, 0));

// ============================================================================
// 2. ARRAY UTILITIES (Native TypeScript array methods)
// ============================================================================

console.log("\n🔧 Array Utilities Demo");

const numbers = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10];

// Map
const doubled = numbers.map(x => x * 2);
console.log("Doubled:", doubled);

// Filter
const evens = numbers.filter(x => x % 2 === 0);
console.log("Evens:", evens);

// Reduce
const sum = numbers.reduce((acc, x) => acc + x, 0);
console.log("Sum:", sum);

// Find
const found = numbers.find(x => x > 5);
console.log("First number > 5:", found);

// Some and Every
const hasEven = numbers.some(x => x % 2 === 0);
const allPositive = numbers.every(x => x > 0);
console.log(`Has even: ${hasEven}, All positive: ${allPositive}`);

// ============================================================================
// 3. PROMISES AND ASYNC/AWAIT
// ============================================================================

console.log("\n⚡ Async/Promise Demo");

// Simple promises
const promise1 = new Promise<string>((resolve) => {
    setTimeout(() => resolve("First result"), 100);
});

const promise2 = new Promise<string>((resolve) => {
    setTimeout(() => resolve("Second result"), 150);
});

// Promise.all
Promise.all([promise1, promise2]).then(results => {
    console.log("All results:", results);
});

// Promise chaining
Promise.resolve(5)
    .then(x => `Value: ${x * 2}`)
    .then(result => console.log("Chained result:", result));

// Async/await
async function asyncDemo() {
    const result = await new Promise<string>(resolve => {
        setTimeout(() => resolve("Async result"), 50);
    });
    console.log("Async/await result:", result);
}

asyncDemo();

// ============================================================================
// 4. CLASSES AND INHERITANCE
// ============================================================================

console.log("\n🏗️  Classes Demo");

class Person {
    constructor(public name: string, public age: number) {}
    
    toString(): string {
        return `Person{Name: ${this.name}, Age: ${this.age}}`;
    }
    
    isAdult(): boolean {
        return this.age >= 18;
    }
}

class Employee extends Person {
    constructor(name: string, age: number, public jobTitle: string, public salary: number) {
        super(name, age);
    }
    
    toString(): string {
        return `Employee{Name: ${this.name}, Age: ${this.age}, JobTitle: ${this.jobTitle}, Salary: ${this.salary}}`;
    }
    
    getAnnualSalary(): number {
        return this.salary * 12;
    }
}

class Manager extends Employee {
    constructor(name: string, age: number, jobTitle: string, salary: number, public teamSize: number) {
        super(name, age, jobTitle, salary);
    }
    
    toString(): string {
        return `Manager{Name: ${this.name}, Age: ${this.age}, JobTitle: ${this.jobTitle}, Salary: ${this.salary}, TeamSize: ${this.teamSize}}`;
    }
    
    isExecutive(): boolean {
        return this.teamSize > 10;
    }
}

// Abstract classes
abstract class Shape {
    constructor(public color: string) {}
    
    abstract area(): number;
    abstract perimeter(): number;
    
    getColor(): string {
        return this.color;
    }
}

class Rectangle extends Shape {
    constructor(color: string, public width: number, public height: number) {
        super(color);
    }
    
    area(): number {
        return this.width * this.height;
    }
    
    perimeter(): number {
        return 2 * (this.width + this.height);
    }
    
    toString(): string {
        return `Rectangle{Color: ${this.color}, Width: ${this.width}, Height: ${this.height}}`;
    }
}

// Usage
const person = new Person("Alice", 25);
console.log("Person:", person.toString());
console.log("Is adult:", person.isAdult());

const employee = new Employee("Bob", 30, "Software Engineer", 75000);
console.log("Employee:", employee.toString());
console.log("Annual salary:", employee.getAnnualSalary());

const manager = new Manager("Carol", 35, "Engineering Manager", 95000, 8);
console.log("Manager:", manager.toString());
console.log("Is executive:", manager.isExecutive());

const rectangle = new Rectangle("blue", 10, 20);
console.log("Rectangle:", rectangle.toString());
console.log(`Area: ${rectangle.area()}, Perimeter: ${rectangle.perimeter()}`);

// ============================================================================
// 5. ENUMS
// ============================================================================

console.log("\n📚 Enums Demo");

// Numeric enum
enum Direction {
    Up,
    Down,
    Left,
    Right
}

// String enum
enum Color {
    Red = "red",
    Green = "green",
    Blue = "blue",
    White = "white",
    Black = "black"
}

// Mixed enum
enum Status {
    Pending = 0,
    InProgress = 1,
    Completed = 2,
    Failed = 999
}

enum LogLevel {
    Debug = 0,
    Info = 1,
    Warn = 2,
    Error = 3,
    Fatal = 4
}

// Usage
const direction = Direction.Up;
console.log(`Direction: ${Direction[direction]} (value: ${direction})`);

const color = Color.Blue;
console.log(`Color: ${color} (name: Blue)`);

const currentStatus = Status.InProgress;
console.log(`Status: ${Status[currentStatus]} (value: ${currentStatus})`);

const logLevel = LogLevel.Info;
const isAtLeastWarn = logLevel >= LogLevel.Warn;
console.log(`Log level: ${LogLevel[logLevel]}, is at least warn: ${isAtLeastWarn}`);

// ============================================================================
// 6. UNION TYPES
// ============================================================================

console.log("\n🔀 Union Types Demo");

type StringOrNumber = string | number;

function processValue(value: StringOrNumber): string {
    if (typeof value === "string") {
        return `Got string: ${value}`;
    } else {
        return `Got number: ${value}`;
    }
}

// Usage
const stringValue: StringOrNumber = "hello";
const numberValue: StringOrNumber = 42;

console.log("String value:", typeof stringValue === "string" ? stringValue : "");
console.log("Number value:", typeof numberValue === "number" ? numberValue : "");
console.log("Pattern match result:", processValue("TypeScript-like"));

// Type guards
function isString(value: any): value is string {
    return typeof value === "string";
}

function isNumber(value: any): value is number {
    return typeof value === "number";
}

const testValue = "hello world";
console.log(`IsString: ${isString(testValue)}, IsNumber: ${isNumber(testValue)}`);

// ============================================================================
// 7. STRING UTILITIES
// ============================================================================

console.log("\n🔤 String Utilities Demo");

const text = "Hello TypeScript World";

// Basic operations
console.log("Original:", text);
console.log("Length:", text.length);
console.log("Upper:", text.toUpperCase());
console.log("Substring(0, 5):", text.substring(0, 5));

// Advanced operations
console.log("Starts with 'Hello':", text.startsWith("Hello"));
console.log("Includes 'Script':", text.includes("Script"));
console.log("Index of 'Type':", text.indexOf("Type"));

// Case conversions (would need utility functions in real TypeScript)
function toCamelCase(str: string): string {
    return str.replace(/[-_\s]+(.)?/g, (_, c) => c ? c.toUpperCase() : '');
}

function toPascalCase(str: string): string {
    const camel = toCamelCase(str);
    return camel.charAt(0).toUpperCase() + camel.slice(1);
}

function toKebabCase(str: string): string {
    return str.replace(/([a-z])([A-Z])/g, '$1-$2')
              .replace(/[\s_]+/g, '-')
              .toLowerCase();
}

function toSnakeCase(str: string): string {
    return str.replace(/([a-z])([A-Z])/g, '$1_$2')
              .replace(/[\s-]+/g, '_')
              .toLowerCase();
}

const original = "hello-world_example";
console.log("Original:", original);
console.log("CamelCase:", toCamelCase(original));
console.log("PascalCase:", toPascalCase(original));
console.log("KebabCase:", toKebabCase(original));
console.log("SnakeCase:", toSnakeCase(original));

// Split and join
const words = text.split(" ");
console.log("Words:", words);
const joined = words.join("-");
console.log("Joined:", joined);

// ============================================================================
// 8. JSON HANDLING
// ============================================================================

console.log("\n📄 JSON Handling Demo");

// Create object
const personData = {
    name: "John Doe",
    age: 30,
    city: "New York",
    active: true
};

// JSON stringify
const jsonStr = JSON.stringify(personData, null, 2);
console.log("JSON:", jsonStr);

// JSON parse
const parsed = JSON.parse(jsonStr);
console.log("Parsed:", parsed);

// Object utilities
const keys = Object.keys(personData);
const values = Object.values(personData);
console.log("Keys:", keys);
console.log("Values:", values);

// Object manipulation
function pick<T extends object, K extends keyof T>(obj: T, ...keys: K[]): Pick<T, K> {
    const result = {} as Pick<T, K>;
    keys.forEach(key => {
        if (key in obj) {
            result[key] = obj[key];
        }
    });
    return result;
}

function omit<T extends object, K extends keyof T>(obj: T, ...keys: K[]): Omit<T, K> {
    const result = { ...obj };
    keys.forEach(key => delete result[key]);
    return result;
}

const picked = pick(personData, "name", "age");
console.log("Picked:", picked);

const omitted = omit(personData, "active");
console.log("Omitted:", omitted);

// Merge objects
const extra = { country: "USA", zipcode: "10001" };
const merged = { ...personData, ...extra };
console.log("Merged:", merged);

// ============================================================================
// 9. COLLECTIONS (Map and Set)
// ============================================================================

console.log("\n🗂️  Collections Demo");

// Map demo
const userMap = new Map<string, string>();
userMap.set("alice", "Alice Johnson");
userMap.set("bob", "Bob Smith");
userMap.set("charlie", "Charlie Brown");

console.log("Map size:", userMap.size);
console.log("Alice's full name:", userMap.get("alice"));
console.log("Map keys:", Array.from(userMap.keys()));
console.log("Map values:", Array.from(userMap.values()));

// Map iteration
console.log("Map entries:");
userMap.forEach((value, key) => {
    console.log(`  ${key} -> ${value}`);
});

// Set demo
const numbersSet = new Set<number>();
numbersSet.add(1).add(2).add(3).add(2); // Adding duplicate

console.log("\nSet size:", numbersSet.size);
console.log("Set contains 2:", numbersSet.has(2));
console.log("Set values:", Array.from(numbersSet.values()));

// Set operations (utility functions)
function union<T>(set1: Set<T>, set2: Set<T>): Set<T> {
    return new Set([...set1, ...set2]);
}

function intersection<T>(set1: Set<T>, set2: Set<T>): Set<T> {
    return new Set([...set1].filter(x => set2.has(x)));
}

function difference<T>(set1: Set<T>, set2: Set<T>): Set<T> {
    return new Set([...set1].filter(x => !set2.has(x)));
}

const otherSet = new Set([3, 4, 5]);
console.log("Union:", Array.from(union(numbersSet, otherSet)));
console.log("Intersection:", Array.from(intersection(numbersSet, otherSet)));
console.log("Difference:", Array.from(difference(numbersSet, otherSet)));

// ============================================================================
// 10. EVENT HANDLING (EventTarget/EventEmitter pattern)
// ============================================================================

console.log("\n📡 Events Demo");

// EventEmitter-like class
class EventEmitter<T = any> {
    private listeners = new Map<string, Array<(data: T) => void>>();
    private onceListeners = new Map<string, Array<(data: T) => void>>();

    on(event: string, listener: (data: T) => void): this {
        if (!this.listeners.has(event)) {
            this.listeners.set(event, []);
        }
        this.listeners.get(event)!.push(listener);
        return this;
    }

    once(event: string, listener: (data: T) => void): this {
        if (!this.onceListeners.has(event)) {
            this.onceListeners.set(event, []);
        }
        this.onceListeners.get(event)!.push(listener);
        return this;
    }

    off(event: string, listener: (data: T) => void): this {
        // Remove from regular listeners
        const listeners = this.listeners.get(event);
        if (listeners) {
            const index = listeners.indexOf(listener);
            if (index > -1) {
                listeners.splice(index, 1);
            }
        }

        // Remove from once listeners
        const onceListeners = this.onceListeners.get(event);
        if (onceListeners) {
            const index = onceListeners.indexOf(listener);
            if (index > -1) {
                onceListeners.splice(index, 1);
            }
        }

        return this;
    }

    emit(event: string, data: T): boolean {
        let hadListeners = false;

        // Emit to regular listeners
        const listeners = this.listeners.get(event);
        if (listeners) {
            listeners.forEach(listener => listener(data));
            hadListeners = true;
        }

        // Emit to once listeners and remove them
        const onceListeners = this.onceListeners.get(event);
        if (onceListeners && onceListeners.length > 0) {
            onceListeners.forEach(listener => listener(data));
            this.onceListeners.delete(event);
            hadListeners = true;
        }

        return hadListeners;
    }

    listenerCount(event: string): number {
        const regularCount = this.listeners.get(event)?.length || 0;
        const onceCount = this.onceListeners.get(event)?.length || 0;
        return regularCount + onceCount;
    }

    eventNames(): string[] {
        const events = new Set<string>();
        this.listeners.forEach((_, event) => events.add(event));
        this.onceListeners.forEach((_, event) => events.add(event));
        return Array.from(events);
    }
}

// Usage
const emitter = new EventEmitter<string>();

// Add listeners
emitter.on("message", (data: string) => {
    console.log("Listener 1 received:", data);
});

emitter.on("message", (data: string) => {
    console.log("Listener 2 received:", data);
});

// One-time listener
emitter.once("startup", (data: string) => {
    console.log("Startup event:", data);
});

// Emit events
console.log("Emitting 'message' event:");
emitter.emit("message", "Hello from EventEmitter!");

console.log("\nEmitting 'startup' event (once):");
emitter.emit("startup", "System initialized");
emitter.emit("startup", "This won't be heard"); // Won't trigger once listener

console.log("Listener count for 'message':", emitter.listenerCount("message"));
console.log("Event names:", emitter.eventNames());

// ============================================================================
// 11. ERROR HANDLING
// ============================================================================

console.log("\n🚨 Error Handling Demo");

// Enhanced Error class
class EnhancedError extends Error {
    public readonly code: string;
    public readonly timestamp: Date;
    public readonly data: Record<string, any>;
    public readonly cause?: Error;

    constructor(message: string, code: string = "UNKNOWN_ERROR", cause?: Error) {
        super(message);
        this.name = "EnhancedError";
        this.code = code;
        this.timestamp = new Date();
        this.data = {};
        this.cause = cause;
    }

    withData(key: string, value: any): this {
        this.data[key] = value;
        return this;
    }

    toJSON() {
        return {
            name: this.name,
            message: this.message,
            code: this.code,
            timestamp: this.timestamp.toISOString(),
            data: this.data,
            stack: this.stack,
            cause: this.cause?.message
        };
    }
}

// Try-catch pattern
class Try<T> {
    constructor(private fn: () => T) {}

    catch(handler: (error: Error) => void): this {
        this.errorHandler = handler;
        return this;
    }

    finally(handler: () => void): this {
        this.finallyHandler = handler;
        return this;
    }

    execute(): T | undefined {
        try {
            const result = this.fn();
            return result;
        } catch (error) {
            if (this.errorHandler) {
                this.errorHandler(error as Error);
            }
            return undefined;
        } finally {
            if (this.finallyHandler) {
                this.finallyHandler();
            }
        }
    }

    private errorHandler?: (error: Error) => void;
    private finallyHandler?: () => void;
}

// Usage
const err1 = new EnhancedError("Something went wrong", "VALIDATION_ERROR");
err1.withData("field", "username").withData("value", "");

console.log("Basic error:", err1.message);
console.log("Error code:", err1.code);
console.log("Error data:", err1.data);

// Error with cause
const originalErr = new Error("network connection failed");
const err2 = new EnhancedError("Failed to fetch user data", "NETWORK_ERROR", originalErr);

console.log("\nChained error:", err2.message);
console.log("Cause:", err2.cause?.message);

// Try-catch pattern
console.log("\nTry-catch demo:");
const result = new Try(() => {
    throw new EnhancedError("Invalid input", "VALIDATION_ERROR");
    return "success";
})
.catch((error: Error) => {
    if (error instanceof EnhancedError) {
        console.log(`Caught error: [${error.code}] ${error.message}`);
    }
})
.finally(() => {
    console.log("Finally block executed");
})
.execute();

console.log("Try-catch result:", result);

// ============================================================================
// 12. TESTING FRAMEWORK CONCEPTS (Jest/Mocha style)
// ============================================================================

console.log("\n🧪 Testing Framework Demo");

// This would typically be in a separate test file with Jest or Mocha
interface TestContext {
    description: string;
}

type TestFunction = (ctx: TestContext) => void | Promise<void>;

class TestSuite {
    private tests: Array<{ description: string; fn: TestFunction }> = [];
    private beforeEachHooks: Array<(ctx: TestContext) => void> = [];
    private afterEachHooks: Array<(ctx: TestContext) => void> = [];

    constructor(public name: string) {}

    it(description: string, fn: TestFunction): void {
        this.tests.push({ description, fn });
    }

    beforeEach(fn: (ctx: TestContext) => void): void {
        this.beforeEachHooks.push(fn);
    }

    afterEach(fn: (ctx: TestContext) => void): void {
        this.afterEachHooks.push(fn);
    }

    async run(): Promise<void> {
        console.log(`\n📁 ${this.name}`);
        
        for (const test of this.tests) {
            const ctx: TestContext = { description: test.description };
            
            try {
                // Run beforeEach hooks
                this.beforeEachHooks.forEach(hook => hook(ctx));
                
                // Run the test
                await test.fn(ctx);
                
                console.log(`  ✅ ${test.description}`);
                
                // Run afterEach hooks
                this.afterEachHooks.forEach(hook => hook(ctx));
            } catch (error) {
                console.log(`  ❌ ${test.description}`);
                console.log(`    Error: ${error}`);
            }
        }
    }
}

// Expectation/assertion framework
class Expectation<T> {
    constructor(private actual: T) {}

    toBe(expected: T): void {
        if (this.actual !== expected) {
            throw new Error(`Expected ${this.actual} to be ${expected}`);
        }
    }

    toEqual(expected: T): void {
        if (JSON.stringify(this.actual) !== JSON.stringify(expected)) {
            throw new Error(`Expected ${JSON.stringify(this.actual)} to equal ${JSON.stringify(expected)}`);
        }
    }

    toContain(expected: any): void {
        if (Array.isArray(this.actual)) {
            if (!this.actual.includes(expected)) {
                throw new Error(`Expected ${JSON.stringify(this.actual)} to contain ${expected}`);
            }
        } else if (typeof this.actual === 'string') {
            if (!this.actual.includes(expected)) {
                throw new Error(`Expected "${this.actual}" to contain "${expected}"`);
            }
        }
    }

    toHaveLength(expected: number): void {
        const length = Array.isArray(this.actual) ? this.actual.length : 
                      typeof this.actual === 'string' ? this.actual.length : 0;
        if (length !== expected) {
            throw new Error(`Expected length ${length} to be ${expected}`);
        }
    }
}

function expect<T>(actual: T): Expectation<T> {
    return new Expectation(actual);
}

// Example test suite
async function runTests() {
    const arrayTests = new TestSuite("Array Utilities");
    
    arrayTests.beforeEach((ctx) => {
        console.log(`    Starting: ${ctx.description}`);
    });
    
    arrayTests.it("should map array elements correctly", () => {
        const numbers = [1, 2, 3, 4, 5];
        const doubled = numbers.map(x => x * 2);
        
        expect(doubled).toEqual([2, 4, 6, 8, 10]);
        expect(doubled).toHaveLength(5);
        expect(doubled).toContain(6);
    });
    
    arrayTests.it("should filter array elements correctly", () => {
        const numbers = [1, 2, 3, 4, 5, 6];
        const evens = numbers.filter(x => x % 2 === 0);
        
        expect(evens).toEqual([2, 4, 6]);
        expect(evens).toHaveLength(3);
    });
    
    await arrayTests.run();
}

console.log("Running TypeScript test suite...");
runTests();

// ============================================================================
// SUMMARY
// ============================================================================

console.log("\n🎯 TypeScript Comparison Summary");
console.log("================================");
console.log("This TypeScript file demonstrates all the patterns we've");
console.log("successfully implemented in our Go TypeScript-like library:");
console.log("- Optional/nullable types ✅");
console.log("- Array utilities (map, filter, reduce, etc.) ✅");
console.log("- Promise-based async programming ✅");
console.log("- Class-based OOP with inheritance ✅");
console.log("- Enums (numeric and string) ✅");
console.log("- Union types and type guards ✅");
console.log("- String manipulation utilities ✅");
console.log("- JSON handling ✅");
console.log("- Collections (Map, Set) ✅");
console.log("- Event handling (EventEmitter pattern) ✅");
console.log("- Enhanced error handling ✅");
console.log("- Testing framework (Jest/Mocha style) ✅");
console.log("\nOur Go implementation successfully provides TypeScript-like");
console.log("developer experience while leveraging Go's performance!");