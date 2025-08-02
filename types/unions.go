package types

import (
	"fmt"
	"reflect"
)

// Union represents TypeScript's union types (T | U)
type Union interface {
	Value() interface{}
	Type() reflect.Type
	Is(targetType reflect.Type) bool
	As(targetType reflect.Type) (interface{}, bool)
}

// union is the concrete implementation of Union
type union struct {
	value interface{}
	typ   reflect.Type
}

// NewUnion creates a new union value
func NewUnion(value interface{}) Union {
	return &union{
		value: value,
		typ:   reflect.TypeOf(value),
	}
}

func (u *union) Value() interface{} {
	return u.value
}

func (u *union) Type() reflect.Type {
	return u.typ
}

func (u *union) Is(targetType reflect.Type) bool {
	return u.typ == targetType
}

func (u *union) As(targetType reflect.Type) (interface{}, bool) {
	if u.typ == targetType {
		return u.value, true
	}
	return nil, false
}

// String representation
func (u *union) String() string {
	return fmt.Sprintf("Union{%v: %v}", u.typ, u.value)
}

// StringOrNumber represents TypeScript's string | number
type StringOrNumber struct {
	Union
}

// NewStringOrNumber creates string | number union
func NewStringOrNumber(value interface{}) (*StringOrNumber, error) {
	switch value.(type) {
	case string, int, int32, int64, float32, float64:
		return &StringOrNumber{NewUnion(value)}, nil
	default:
		return nil, fmt.Errorf("value must be string or number, got %T", value)
	}
}

// AsString safely extracts string value
func (s *StringOrNumber) AsString() (string, bool) {
	if val, ok := s.As(reflect.TypeOf("")); ok {
		return val.(string), true
	}
	return "", false
}

// AsNumber safely extracts numeric value as float64
func (s *StringOrNumber) AsNumber() (float64, bool) {
	switch v := s.Value().(type) {
	case int:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case float32:
		return float64(v), true
	case float64:
		return v, true
	default:
		return 0, false
	}
}

// StringOrBool represents TypeScript's string | boolean
type StringOrBool struct {
	Union
}

// NewStringOrBool creates string | boolean union
func NewStringOrBool(value interface{}) (*StringOrBool, error) {
	switch value.(type) {
	case string, bool:
		return &StringOrBool{NewUnion(value)}, nil
	default:
		return nil, fmt.Errorf("value must be string or boolean, got %T", value)
	}
}

// AsString safely extracts string value
func (s *StringOrBool) AsString() (string, bool) {
	if val, ok := s.As(reflect.TypeOf("")); ok {
		return val.(string), true
	}
	return "", false
}

// AsBool safely extracts boolean value
func (s *StringOrBool) AsBool() (bool, bool) {
	if val, ok := s.As(reflect.TypeOf(true)); ok {
		return val.(bool), true
	}
	return false, false
}

// UnionMatcher provides pattern matching for unions
// Similar to TypeScript's type guards and switch statements
type UnionMatcher[T any] struct {
	union Union
}

// NewMatcher creates a new union matcher
func NewMatcher[T any](u Union) *UnionMatcher[T] {
	return &UnionMatcher[T]{union: u}
}

// Match performs pattern matching on union types
func (m *UnionMatcher[T]) Match(cases map[reflect.Type]func(interface{}) T, defaultCase func(interface{}) T) T {
	if handler, exists := cases[m.union.Type()]; exists {
		return handler(m.union.Value())
	}
	if defaultCase != nil {
		return defaultCase(m.union.Value())
	}
	var zero T
	return zero
}

// TypeGuard creates a type guard function
// Similar to TypeScript's user-defined type guards
type TypeGuard[T any] func(value interface{}) bool

// IsString type guard
func IsString(value interface{}) bool {
	_, ok := value.(string)
	return ok
}

// IsNumber type guard
func IsNumber(value interface{}) bool {
	switch value.(type) {
	case int, int32, int64, float32, float64:
		return true
	default:
		return false
	}
}

// IsBool type guard
func IsBool(value interface{}) bool {
	_, ok := value.(bool)
	return ok
}

// IsArray type guard
func IsArray(value interface{}) bool {
	return reflect.TypeOf(value).Kind() == reflect.Slice
}

// IsObject type guard (checks for struct)
func IsObject(value interface{}) bool {
	return reflect.TypeOf(value).Kind() == reflect.Struct
}

// CreateTypeGuard creates a generic type guard
func CreateTypeGuard[T any]() TypeGuard[T] {
	var zero T
	targetType := reflect.TypeOf(zero)
	return func(value interface{}) bool {
		return reflect.TypeOf(value) == targetType
	}
}

// Discriminated union support
// DiscriminatedUnion represents TypeScript's discriminated unions
type DiscriminatedUnion interface {
	Union
	Discriminator() string
}

type discriminatedUnion struct {
	*union
	discriminator string
}

// NewDiscriminatedUnion creates a discriminated union
func NewDiscriminatedUnion(value interface{}, discriminator string) DiscriminatedUnion {
	return &discriminatedUnion{
		union:         &union{value: value, typ: reflect.TypeOf(value)},
		discriminator: discriminator,
	}
}

func (d *discriminatedUnion) Discriminator() string {
	return d.discriminator
}

// Narrow narrows a union type based on a type guard
func Narrow[T any](u Union, guard TypeGuard[T]) *T {
	if guard(u.Value()) {
		if val, ok := u.Value().(T); ok {
			return &val
		}
	}
	return nil
}