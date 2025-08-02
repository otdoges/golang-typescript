package enums

import (
	"fmt"
	"strconv"
	"strings"
)

// Enum interface represents TypeScript enum functionality
type Enum interface {
	String() string
	Value() interface{}
	Name() string
	Ordinal() int
}

// NumericEnum represents TypeScript's numeric enums
type NumericEnum struct {
	name    string
	value   int
	ordinal int
}

// NewNumericEnum creates a new numeric enum value
func NewNumericEnum(name string, value int, ordinal int) *NumericEnum {
	return &NumericEnum{
		name:    name,
		value:   value,
		ordinal: ordinal,
	}
}

func (e *NumericEnum) String() string {
	return e.name
}

func (e *NumericEnum) Value() interface{} {
	return e.value
}

func (e *NumericEnum) Name() string {
	return e.name
}

func (e *NumericEnum) Ordinal() int {
	return e.ordinal
}

// StringEnum represents TypeScript's string enums
type StringEnum struct {
	name    string
	value   string
	ordinal int
}

// NewStringEnum creates a new string enum value
func NewStringEnum(name string, value string, ordinal int) *StringEnum {
	return &StringEnum{
		name:    name,
		value:   value,
		ordinal: ordinal,
	}
}

func (e *StringEnum) String() string {
	return e.name
}

func (e *StringEnum) Value() interface{} {
	return e.value
}

func (e *StringEnum) Name() string {
	return e.name
}

func (e *StringEnum) Ordinal() int {
	return e.ordinal
}

// Direction numeric enum (like TypeScript: enum Direction { Up, Down, Left, Right })
type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

var directionNames = map[Direction]string{
	Up:    "Up",
	Down:  "Down",
	Left:  "Left",
	Right: "Right",
}

var directionValues = map[string]Direction{
	"Up":    Up,
	"Down":  Down,
	"Left":  Left,
	"Right": Right,
}

func (d Direction) String() string {
	if name, ok := directionNames[d]; ok {
		return name
	}
	return strconv.Itoa(int(d))
}

func (d Direction) Value() interface{} {
	return int(d)
}

func (d Direction) Name() string {
	return d.String()
}

func (d Direction) Ordinal() int {
	return int(d)
}

// ParseDirection parses string to Direction
func ParseDirection(s string) (Direction, error) {
	if dir, ok := directionValues[s]; ok {
		return dir, nil
	}
	return 0, fmt.Errorf("invalid direction: %s", s)
}

// GetAllDirections returns all direction values
func GetAllDirections() []Direction {
	return []Direction{Up, Down, Left, Right}
}

// Color string enum (like TypeScript: enum Color { Red = "red", Green = "green", Blue = "blue" })
type Color string

const (
	Red   Color = "red"
	Green Color = "green"
	Blue  Color = "blue"
	White Color = "white"
	Black Color = "black"
)

var colorOrdinals = map[Color]int{
	Red:   0,
	Green: 1,
	Blue:  2,
	White: 3,
	Black: 4,
}

var colorValues = map[string]Color{
	"red":   Red,
	"green": Green,
	"blue":  Blue,
	"white": White,
	"black": Black,
}

func (c Color) String() string {
	return string(c)
}

func (c Color) Value() interface{} {
	return string(c)
}

func (c Color) Name() string {
	return strings.Title(string(c))
}

func (c Color) Ordinal() int {
	if ord, ok := colorOrdinals[c]; ok {
		return ord
	}
	return -1
}

// ParseColor parses string to Color
func ParseColor(s string) (Color, error) {
	if color, ok := colorValues[strings.ToLower(s)]; ok {
		return color, nil
	}
	return "", fmt.Errorf("invalid color: %s", s)
}

// GetAllColors returns all color values
func GetAllColors() []Color {
	return []Color{Red, Green, Blue, White, Black}
}

// Status enum with mixed values (like TypeScript heterogeneous enums)
type Status int

const (
	Pending Status = iota
	InProgress
	Completed
	Failed = 999
)

var statusNames = map[Status]string{
	Pending:    "Pending",
	InProgress: "InProgress",
	Completed:  "Completed",
	Failed:     "Failed",
}

var statusValues = map[string]Status{
	"Pending":    Pending,
	"InProgress": InProgress,
	"Completed":  Completed,
	"Failed":     Failed,
}

func (s Status) String() string {
	if name, ok := statusNames[s]; ok {
		return name
	}
	return strconv.Itoa(int(s))
}

func (s Status) Value() interface{} {
	return int(s)
}

func (s Status) Name() string {
	return s.String()
}

func (s Status) Ordinal() int {
	switch s {
	case Pending:
		return 0
	case InProgress:
		return 1
	case Completed:
		return 2
	case Failed:
		return 3
	default:
		return -1
	}
}

// ParseStatus parses string to Status
func ParseStatus(s string) (Status, error) {
	if status, ok := statusValues[s]; ok {
		return status, nil
	}
	return 0, fmt.Errorf("invalid status: %s", s)
}

// GetAllStatuses returns all status values
func GetAllStatuses() []Status {
	return []Status{Pending, InProgress, Completed, Failed}
}

// EnumSet represents a set of enum values (like TypeScript's const enum usage)
type EnumSet[T Enum] struct {
	values map[string]T
	names  []string
}

// NewEnumSet creates a new enum set
func NewEnumSet[T Enum](values ...T) *EnumSet[T] {
	set := &EnumSet[T]{
		values: make(map[string]T),
		names:  make([]string, 0, len(values)),
	}
	
	for _, value := range values {
		name := value.Name()
		set.values[name] = value
		set.names = append(set.names, name)
	}
	
	return set
}

// Get returns enum value by name
func (es *EnumSet[T]) Get(name string) (T, bool) {
	value, exists := es.values[name]
	return value, exists
}

// Has checks if enum value exists
func (es *EnumSet[T]) Has(name string) bool {
	_, exists := es.values[name]
	return exists
}

// GetNames returns all enum names
func (es *EnumSet[T]) GetNames() []string {
	return append([]string(nil), es.names...)
}

// GetValues returns all enum values
func (es *EnumSet[T]) GetValues() []T {
	values := make([]T, 0, len(es.names))
	for _, name := range es.names {
		values = append(values, es.values[name])
	}
	return values
}

// Size returns the number of enum values
func (es *EnumSet[T]) Size() int {
	return len(es.names)
}

// LogLevel enum example with explicit values
type LogLevel int

const (
	Debug LogLevel = 0
	Info  LogLevel = 1
	Warn  LogLevel = 2
	Error LogLevel = 3
	Fatal LogLevel = 4
)

var logLevelNames = map[LogLevel]string{
	Debug: "Debug",
	Info:  "Info",
	Warn:  "Warn",
	Error: "Error",
	Fatal: "Fatal",
}

var logLevelValues = map[string]LogLevel{
	"Debug": Debug,
	"Info":  Info,
	"Warn":  Warn,
	"Error": Error,
	"Fatal": Fatal,
}

func (l LogLevel) String() string {
	if name, ok := logLevelNames[l]; ok {
		return name
	}
	return strconv.Itoa(int(l))
}

func (l LogLevel) Value() interface{} {
	return int(l)
}

func (l LogLevel) Name() string {
	return l.String()
}

func (l LogLevel) Ordinal() int {
	return int(l)
}

// IsAtLeast checks if log level is at least the specified level
func (l LogLevel) IsAtLeast(level LogLevel) bool {
	return l >= level
}

// ParseLogLevel parses string to LogLevel
func ParseLogLevel(s string) (LogLevel, error) {
	if level, ok := logLevelValues[s]; ok {
		return level, nil
	}
	return 0, fmt.Errorf("invalid log level: %s", s)
}

// GetAllLogLevels returns all log level values
func GetAllLogLevels() []LogLevel {
	return []LogLevel{Debug, Info, Warn, Error, Fatal}
}

// EnumUtils provides utility functions for enums
type EnumUtils struct{}

// ToString converts any enum to string
func (EnumUtils) ToString(e Enum) string {
	return e.String()
}

// GetValue gets the underlying value of any enum
func (EnumUtils) GetValue(e Enum) interface{} {
	return e.Value()
}

// Compare compares two enums by their ordinal values
func (EnumUtils) Compare(e1, e2 Enum) int {
	ord1, ord2 := e1.Ordinal(), e2.Ordinal()
	if ord1 < ord2 {
		return -1
	}
	if ord1 > ord2 {
		return 1
	}
	return 0
}

// Global enum utilities instance
var Utils = EnumUtils{}