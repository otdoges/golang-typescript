// TypeScript Implementation Comparison
// This file demonstrates the same functionality we've implemented in Go
// to validate how closely our Go implementation matches TypeScript patterns
function someValue(value) {
    return value;
}
function getValue(opt, defaultValue) {
    return opt ?? defaultValue;
}
// Usage
const userName = "John Doe";
const userAge = undefined;
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
const promise1 = new Promise((resolve) => {
    setTimeout(() => resolve("First result"), 100);
});
const promise2 = new Promise((resolve) => {
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
    const result = await new Promise(resolve => {
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
    constructor(name, age) {
        this.name = name;
        this.age = age;
    }
    toString() {
        return `Person{Name: ${this.name}, Age: ${this.age}}`;
    }
    isAdult() {
        return this.age >= 18;
    }
}
class Employee extends Person {
    constructor(name, age, jobTitle, salary) {
        super(name, age);
        this.jobTitle = jobTitle;
        this.salary = salary;
    }
    toString() {
        return `Employee{Name: ${this.name}, Age: ${this.age}, JobTitle: ${this.jobTitle}, Salary: ${this.salary}}`;
    }
    getAnnualSalary() {
        return this.salary * 12;
    }
}
class Manager extends Employee {
    constructor(name, age, jobTitle, salary, teamSize) {
        super(name, age, jobTitle, salary);
        this.teamSize = teamSize;
    }
    toString() {
        return `Manager{Name: ${this.name}, Age: ${this.age}, JobTitle: ${this.jobTitle}, Salary: ${this.salary}, TeamSize: ${this.teamSize}}`;
    }
    isExecutive() {
        return this.teamSize > 10;
    }
}
// Abstract classes
class Shape {
    constructor(color) {
        this.color = color;
    }
    getColor() {
        return this.color;
    }
}
class Rectangle extends Shape {
    constructor(color, width, height) {
        super(color);
        this.width = width;
        this.height = height;
    }
    area() {
        return this.width * this.height;
    }
    perimeter() {
        return 2 * (this.width + this.height);
    }
    toString() {
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
var Direction;
(function (Direction) {
    Direction[Direction["Up"] = 0] = "Up";
    Direction[Direction["Down"] = 1] = "Down";
    Direction[Direction["Left"] = 2] = "Left";
    Direction[Direction["Right"] = 3] = "Right";
})(Direction || (Direction = {}));
// String enum
var Color;
(function (Color) {
    Color["Red"] = "red";
    Color["Green"] = "green";
    Color["Blue"] = "blue";
    Color["White"] = "white";
    Color["Black"] = "black";
})(Color || (Color = {}));
// Mixed enum
var Status;
(function (Status) {
    Status[Status["Pending"] = 0] = "Pending";
    Status[Status["InProgress"] = 1] = "InProgress";
    Status[Status["Completed"] = 2] = "Completed";
    Status[Status["Failed"] = 999] = "Failed";
})(Status || (Status = {}));
var LogLevel;
(function (LogLevel) {
    LogLevel[LogLevel["Debug"] = 0] = "Debug";
    LogLevel[LogLevel["Info"] = 1] = "Info";
    LogLevel[LogLevel["Warn"] = 2] = "Warn";
    LogLevel[LogLevel["Error"] = 3] = "Error";
    LogLevel[LogLevel["Fatal"] = 4] = "Fatal";
})(LogLevel || (LogLevel = {}));
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
function processValue(value) {
    if (typeof value === "string") {
        return `Got string: ${value}`;
    }
    else {
        return `Got number: ${value}`;
    }
}
// Usage
const stringValue = "hello";
const numberValue = 42;
console.log("String value:", typeof stringValue === "string" ? stringValue : "");
console.log("Number value:", typeof numberValue === "number" ? numberValue : "");
console.log("Pattern match result:", processValue("TypeScript-like"));
// Type guards
function isString(value) {
    return typeof value === "string";
}
function isNumber(value) {
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
function toCamelCase(str) {
    return str.replace(/[-_\s]+(.)?/g, (_, c) => c ? c.toUpperCase() : '');
}
function toPascalCase(str) {
    const camel = toCamelCase(str);
    return camel.charAt(0).toUpperCase() + camel.slice(1);
}
function toKebabCase(str) {
    return str.replace(/([a-z])([A-Z])/g, '$1-$2')
        .replace(/[\s_]+/g, '-')
        .toLowerCase();
}
function toSnakeCase(str) {
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
function pick(obj, ...keys) {
    const result = {};
    keys.forEach(key => {
        if (key in obj) {
            result[key] = obj[key];
        }
    });
    return result;
}
function omit(obj, ...keys) {
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
const userMap = new Map();
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
const numbersSet = new Set();
numbersSet.add(1).add(2).add(3).add(2); // Adding duplicate
console.log("\nSet size:", numbersSet.size);
console.log("Set contains 2:", numbersSet.has(2));
console.log("Set values:", Array.from(numbersSet.values()));
// Set operations (utility functions)
function union(set1, set2) {
    return new Set([...set1, ...set2]);
}
function intersection(set1, set2) {
    return new Set([...set1].filter(x => set2.has(x)));
}
function difference(set1, set2) {
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
class EventEmitter {
    constructor() {
        this.listeners = new Map();
        this.onceListeners = new Map();
    }
    on(event, listener) {
        if (!this.listeners.has(event)) {
            this.listeners.set(event, []);
        }
        this.listeners.get(event).push(listener);
        return this;
    }
    once(event, listener) {
        if (!this.onceListeners.has(event)) {
            this.onceListeners.set(event, []);
        }
        this.onceListeners.get(event).push(listener);
        return this;
    }
    off(event, listener) {
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
    emit(event, data) {
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
    listenerCount(event) {
        const regularCount = this.listeners.get(event)?.length || 0;
        const onceCount = this.onceListeners.get(event)?.length || 0;
        return regularCount + onceCount;
    }
    eventNames() {
        const events = new Set();
        this.listeners.forEach((_, event) => events.add(event));
        this.onceListeners.forEach((_, event) => events.add(event));
        return Array.from(events);
    }
}
// Usage
const emitter = new EventEmitter();
// Add listeners
emitter.on("message", (data) => {
    console.log("Listener 1 received:", data);
});
emitter.on("message", (data) => {
    console.log("Listener 2 received:", data);
});
// One-time listener
emitter.once("startup", (data) => {
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
    constructor(message, code = "UNKNOWN_ERROR", cause) {
        super(message);
        this.name = "EnhancedError";
        this.code = code;
        this.timestamp = new Date();
        this.data = {};
        this.cause = cause;
    }
    withData(key, value) {
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
class Try {
    constructor(fn) {
        this.fn = fn;
    }
    catch(handler) {
        this.errorHandler = handler;
        return this;
    }
    finally(handler) {
        this.finallyHandler = handler;
        return this;
    }
    execute() {
        try {
            const result = this.fn();
            return result;
        }
        catch (error) {
            if (this.errorHandler) {
                this.errorHandler(error);
            }
            return undefined;
        }
        finally {
            if (this.finallyHandler) {
                this.finallyHandler();
            }
        }
    }
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
    .catch((error) => {
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
class TestSuite {
    constructor(name) {
        this.name = name;
        this.tests = [];
        this.beforeEachHooks = [];
        this.afterEachHooks = [];
    }
    it(description, fn) {
        this.tests.push({ description, fn });
    }
    beforeEach(fn) {
        this.beforeEachHooks.push(fn);
    }
    afterEach(fn) {
        this.afterEachHooks.push(fn);
    }
    async run() {
        console.log(`\n📁 ${this.name}`);
        for (const test of this.tests) {
            const ctx = { description: test.description };
            try {
                // Run beforeEach hooks
                this.beforeEachHooks.forEach(hook => hook(ctx));
                // Run the test
                await test.fn(ctx);
                console.log(`  ✅ ${test.description}`);
                // Run afterEach hooks
                this.afterEachHooks.forEach(hook => hook(ctx));
            }
            catch (error) {
                console.log(`  ❌ ${test.description}`);
                console.log(`    Error: ${error}`);
            }
        }
    }
}
// Expectation/assertion framework
class Expectation {
    constructor(actual) {
        this.actual = actual;
    }
    toBe(expected) {
        if (this.actual !== expected) {
            throw new Error(`Expected ${this.actual} to be ${expected}`);
        }
    }
    toEqual(expected) {
        if (JSON.stringify(this.actual) !== JSON.stringify(expected)) {
            throw new Error(`Expected ${JSON.stringify(this.actual)} to equal ${JSON.stringify(expected)}`);
        }
    }
    toContain(expected) {
        if (Array.isArray(this.actual)) {
            if (!this.actual.includes(expected)) {
                throw new Error(`Expected ${JSON.stringify(this.actual)} to contain ${expected}`);
            }
        }
        else if (typeof this.actual === 'string') {
            if (!this.actual.includes(expected)) {
                throw new Error(`Expected "${this.actual}" to contain "${expected}"`);
            }
        }
    }
    toHaveLength(expected) {
        const length = Array.isArray(this.actual) ? this.actual.length :
            typeof this.actual === 'string' ? this.actual.length : 0;
        if (length !== expected) {
            throw new Error(`Expected length ${length} to be ${expected}`);
        }
    }
}
function expect(actual) {
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
