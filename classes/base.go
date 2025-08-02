package classes

import (
	"fmt"
	"reflect"
	"typescript-golang/types"
)

// Class represents the base class interface
// Similar to TypeScript's class declaration
type Class interface {
	GetClassName() string
	ToString() string
	GetType() reflect.Type
	Clone() Class
}

// BaseClass provides default class functionality
// Similar to TypeScript's base class with constructor
type BaseClass struct {
	className string
	types.DefaultStringable
	types.JSONSerializable
}

// NewBaseClass creates a new base class instance (constructor)
func NewBaseClass(className string) *BaseClass {
	return &BaseClass{
		className: className,
	}
}

// GetClassName returns the class name
func (b *BaseClass) GetClassName() string {
	return b.className
}

// ToString provides string representation
func (b *BaseClass) ToString() string {
	return fmt.Sprintf("%s{}", b.className)
}

// GetType returns the reflect.Type
func (b *BaseClass) GetType() reflect.Type {
	return reflect.TypeOf(b)
}

// Clone creates a copy of the class
func (b *BaseClass) Clone() Class {
	return &BaseClass{
		className: b.className,
	}
}

// Constructor function type for class instantiation
type Constructor[T any] func() T

// ClassDefinition represents a class definition with methods
type ClassDefinition[T Class] struct {
	name        string
	constructor Constructor[T]
	methods     map[string]interface{}
	properties  map[string]interface{}
}

// NewClassDefinition creates a new class definition
func NewClassDefinition[T Class](name string, constructor Constructor[T]) *ClassDefinition[T] {
	return &ClassDefinition[T]{
		name:        name,
		constructor: constructor,
		methods:     make(map[string]interface{}),
		properties:  make(map[string]interface{}),
	}
}

// AddMethod adds a method to the class
func (cd *ClassDefinition[T]) AddMethod(name string, method interface{}) {
	cd.methods[name] = method
}

// AddProperty adds a property to the class
func (cd *ClassDefinition[T]) AddProperty(name string, value interface{}) {
	cd.properties[name] = value
}

// Instantiate creates a new instance of the class
func (cd *ClassDefinition[T]) Instantiate() T {
	return cd.constructor()
}

// GetMethod returns a method by name
func (cd *ClassDefinition[T]) GetMethod(name string) (interface{}, bool) {
	method, exists := cd.methods[name]
	return method, exists
}

// Person example class similar to TypeScript class
type Person struct {
	*BaseClass
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// NewPerson creates a new Person (constructor)
func NewPerson(name string, age int) *Person {
	return &Person{
		BaseClass: NewBaseClass("Person"),
		Name:      name,
		Age:       age,
	}
}

// GetName method (public)
func (p *Person) GetName() string {
	return p.Name
}

// SetName method (public)
func (p *Person) SetName(name string) {
	p.Name = name
}

// GetAge method (public)
func (p *Person) GetAge() int {
	return p.Age
}

// SetAge method (public)
func (p *Person) SetAge(age int) {
	p.Age = age
}

// ToString override
func (p *Person) ToString() string {
	return fmt.Sprintf("Person{Name: %s, Age: %d}", p.Name, p.Age)
}

// Clone creates a copy of Person
func (p *Person) Clone() Class {
	return &Person{
		BaseClass: NewBaseClass("Person"),
		Name:      p.Name,
		Age:       p.Age,
	}
}

// IsAdult method
func (p *Person) IsAdult() bool {
	return p.Age >= 18
}

// Employee extends Person (inheritance)
type Employee struct {
	*Person // Embedding for inheritance
	JobTitle string  `json:"jobTitle"`
	Salary   float64 `json:"salary"`
}

// NewEmployee creates a new Employee (constructor)
func NewEmployee(name string, age int, jobTitle string, salary float64) *Employee {
	return &Employee{
		Person:   NewPerson(name, age),
		JobTitle: jobTitle,
		Salary:   salary,
	}
}

// GetJobTitle method
func (e *Employee) GetJobTitle() string {
	return e.JobTitle
}

// SetJobTitle method
func (e *Employee) SetJobTitle(title string) {
	e.JobTitle = title
}

// GetSalary method
func (e *Employee) GetSalary() float64 {
	return e.Salary
}

// SetSalary method
func (e *Employee) SetSalary(salary float64) {
	e.Salary = salary
}

// ToString override (method overriding)
func (e *Employee) ToString() string {
	return fmt.Sprintf("Employee{Name: %s, Age: %d, JobTitle: %s, Salary: %.2f}", 
		e.Name, e.Age, e.JobTitle, e.Salary)
}

// Clone creates a copy of Employee
func (e *Employee) Clone() Class {
	return &Employee{
		Person:   NewPerson(e.Name, e.Age),
		JobTitle: e.JobTitle,
		Salary:   e.Salary,
	}
}

// GetAnnualSalary method (new method in derived class)
func (e *Employee) GetAnnualSalary() float64 {
	return e.Salary * 12
}

// Manager extends Employee (multi-level inheritance)
type Manager struct {
	*Employee
	TeamSize int `json:"teamSize"`
}

// NewManager creates a new Manager (constructor)
func NewManager(name string, age int, jobTitle string, salary float64, teamSize int) *Manager {
	return &Manager{
		Employee: NewEmployee(name, age, jobTitle, salary),
		TeamSize: teamSize,
	}
}

// GetTeamSize method
func (m *Manager) GetTeamSize() int {
	return m.TeamSize
}

// SetTeamSize method
func (m *Manager) SetTeamSize(size int) {
	m.TeamSize = size
}

// ToString override
func (m *Manager) ToString() string {
	return fmt.Sprintf("Manager{Name: %s, Age: %d, JobTitle: %s, Salary: %.2f, TeamSize: %d}", 
		m.Name, m.Age, m.JobTitle, m.Salary, m.TeamSize)
}

// Clone creates a copy of Manager
func (m *Manager) Clone() Class {
	return &Manager{
		Employee: NewEmployee(m.Name, m.Age, m.JobTitle, m.Salary),
		TeamSize: m.TeamSize,
	}
}

// IsExecutive method
func (m *Manager) IsExecutive() bool {
	return m.TeamSize > 10
}

// AbstractClass represents abstract class functionality
type AbstractClass interface {
	Class
	IsAbstract() bool
}

// Shape abstract class example
type Shape struct {
	*BaseClass
	Color string `json:"color"`
}

// NewShape creates a new Shape
func NewShape(color string) *Shape {
	return &Shape{
		BaseClass: NewBaseClass("Shape"),
		Color:     color,
	}
}

// IsAbstract implementation
func (s *Shape) IsAbstract() bool {
	return true
}

// Area is an abstract method (should be overridden)
func (s *Shape) Area() float64 {
	panic("Abstract method Area() must be implemented by derived class")
}

// Perimeter is an abstract method (should be overridden)
func (s *Shape) Perimeter() float64 {
	panic("Abstract method Perimeter() must be implemented by derived class")
}

// GetColor method
func (s *Shape) GetColor() string {
	return s.Color
}

// Rectangle extends Shape
type Rectangle struct {
	*Shape
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// NewRectangle creates a new Rectangle
func NewRectangle(color string, width, height float64) *Rectangle {
	return &Rectangle{
		Shape:  NewShape(color),
		Width:  width,
		Height: height,
	}
}

// Area implementation (overriding abstract method)
func (r *Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Perimeter implementation (overriding abstract method)
func (r *Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// IsAbstract override
func (r *Rectangle) IsAbstract() bool {
	return false
}

// ToString override
func (r *Rectangle) ToString() string {
	return fmt.Sprintf("Rectangle{Color: %s, Width: %.2f, Height: %.2f}", 
		r.Color, r.Width, r.Height)
}

// Clone creates a copy of Rectangle
func (r *Rectangle) Clone() Class {
	return &Rectangle{
		Shape:  NewShape(r.Color),
		Width:  r.Width,
		Height: r.Height,
	}
}