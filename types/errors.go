package types

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

// ErrorCode represents TypeScript-like error codes
type ErrorCode string

const (
	ValidationError   ErrorCode = "VALIDATION_ERROR"
	NetworkError      ErrorCode = "NETWORK_ERROR"
	AuthError         ErrorCode = "AUTH_ERROR"
	NotFoundError     ErrorCode = "NOT_FOUND_ERROR"
	InternalError     ErrorCode = "INTERNAL_ERROR"
	TimeoutError      ErrorCode = "TIMEOUT_ERROR"
	CancelledError    ErrorCode = "CANCELLED_ERROR"
	UnknownError      ErrorCode = "UNKNOWN_ERROR"
)

// StackFrame represents a stack frame with file, line, and function info
type StackFrame struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Function string `json:"function"`
	Package  string `json:"package"`
}

// String returns string representation of stack frame
func (sf StackFrame) String() string {
	return fmt.Sprintf("%s:%d in %s", sf.File, sf.Line, sf.Function)
}

// StackTrace represents a collection of stack frames
type StackTrace []StackFrame

// String returns string representation of stack trace
func (st StackTrace) String() string {
	var lines []string
	for _, frame := range st {
		lines = append(lines, "  at "+frame.String())
	}
	return strings.Join(lines, "\n")
}

// EnhancedError represents TypeScript-like error with additional metadata
type EnhancedError struct {
	message   string
	code      ErrorCode
	cause     error
	data      map[string]interface{}
	stack     StackTrace
	timestamp time.Time
}

// NewError creates a new enhanced error (like new Error() in TypeScript)
func NewError(message string, code ...ErrorCode) *EnhancedError {
	errorCode := UnknownError
	if len(code) > 0 {
		errorCode = code[0]
	}
	
	return &EnhancedError{
		message:   message,
		code:      errorCode,
		data:      make(map[string]interface{}),
		stack:     captureStackTrace(2), // Skip NewError and caller
		timestamp: time.Now(),
	}
}

// NewErrorWithCause creates an error with a cause (error chaining)
func NewErrorWithCause(message string, cause error, code ...ErrorCode) *EnhancedError {
	err := NewError(message, code...)
	err.cause = cause
	return err
}

// captureStackTrace captures the current stack trace
func captureStackTrace(skip int) StackTrace {
	var frames StackTrace
	
	// Capture up to 32 stack frames
	for i := skip; i < skip+32; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}
		
		funcName := fn.Name()
		
		// Extract package name
		parts := strings.Split(funcName, ".")
		packageName := ""
		if len(parts) > 1 {
			packageName = strings.Join(parts[:len(parts)-1], ".")
		}
		
		frame := StackFrame{
			File:     file,
			Line:     line,
			Function: funcName,
			Package:  packageName,
		}
		
		frames = append(frames, frame)
		
		// Stop at main.main
		if strings.HasSuffix(funcName, "main.main") {
			break
		}
	}
	
	return frames
}

// Error implements the error interface
func (e *EnhancedError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.message, e.cause)
	}
	return e.message
}

// Message returns the error message
func (e *EnhancedError) Message() string {
	return e.message
}

// Code returns the error code
func (e *EnhancedError) Code() ErrorCode {
	return e.code
}

// Cause returns the cause error (for error chaining)
func (e *EnhancedError) Cause() error {
	return e.cause
}

// Unwrap returns the cause error (for Go 1.13+ error unwrapping)
func (e *EnhancedError) Unwrap() error {
	return e.cause
}

// Stack returns the stack trace
func (e *EnhancedError) Stack() StackTrace {
	return e.stack
}

// Timestamp returns when the error was created
func (e *EnhancedError) Timestamp() time.Time {
	return e.timestamp
}

// Data returns the error data map
func (e *EnhancedError) Data() map[string]interface{} {
	return e.data
}

// WithData adds data to the error (fluent interface)
func (e *EnhancedError) WithData(key string, value interface{}) *EnhancedError {
	e.data[key] = value
	return e
}

// WithCode sets the error code (fluent interface)
func (e *EnhancedError) WithCode(code ErrorCode) *EnhancedError {
	e.code = code
	return e
}

// Is checks if the error is of a specific code (like instanceof in TypeScript)
func (e *EnhancedError) Is(code ErrorCode) bool {
	return e.code == code
}

// String returns a detailed string representation
func (e *EnhancedError) String() string {
	var parts []string
	
	parts = append(parts, fmt.Sprintf("Error: %s", e.message))
	parts = append(parts, fmt.Sprintf("Code: %s", e.code))
	parts = append(parts, fmt.Sprintf("Time: %s", e.timestamp.Format(time.RFC3339)))
	
	if e.cause != nil {
		parts = append(parts, fmt.Sprintf("Caused by: %v", e.cause))
	}
	
	if len(e.data) > 0 {
		parts = append(parts, fmt.Sprintf("Data: %+v", e.data))
	}
	
	if len(e.stack) > 0 {
		parts = append(parts, "Stack trace:")
		parts = append(parts, e.stack.String())
	}
	
	return strings.Join(parts, "\n")
}

// ToJSON returns JSON representation of the error
func (e *EnhancedError) ToJSON() map[string]interface{} {
	result := map[string]interface{}{
		"message":   e.message,
		"code":      string(e.code),
		"timestamp": e.timestamp.Format(time.RFC3339),
		"data":      e.data,
	}
	
	if e.cause != nil {
		result["cause"] = e.cause.Error()
	}
	
	var stackStrings []string
	for _, frame := range e.stack {
		stackStrings = append(stackStrings, frame.String())
	}
	result["stack"] = stackStrings
	
	return result
}

// ErrorHandler represents a function that handles errors
type ErrorHandler func(*EnhancedError)

// Try represents TypeScript-like try-catch functionality
type Try[T any] struct {
	fn          func() (T, error)
	errorHandler ErrorHandler
	finallyFn   func()
}

// NewTry creates a new Try instance (like try-catch in TypeScript)
func NewTry[T any](fn func() (T, error)) *Try[T] {
	return &Try[T]{fn: fn}
}

// Catch adds an error handler (like catch in TypeScript)
func (t *Try[T]) Catch(handler ErrorHandler) *Try[T] {
	t.errorHandler = handler
	return t
}

// Finally adds a finally handler (like finally in TypeScript)
func (t *Try[T]) Finally(fn func()) *Try[T] {
	t.finallyFn = fn
	return t
}

// Execute runs the try-catch-finally block
func (t *Try[T]) Execute() (T, error) {
	var result T
	var err error
	
	// Defer finally block
	if t.finallyFn != nil {
		defer t.finallyFn()
	}
	
	// Execute main function
	result, err = t.fn()
	
	// Handle error if present
	if err != nil {
		var enhancedErr *EnhancedError
		
		// Convert to enhanced error if not already
		if e, ok := err.(*EnhancedError); ok {
			enhancedErr = e
		} else {
			enhancedErr = NewErrorWithCause("Wrapped error", err)
		}
		
		// Call error handler if present
		if t.errorHandler != nil {
			t.errorHandler(enhancedErr)
		}
		
		return result, enhancedErr
	}
	
	return result, nil
}

// Assert represents TypeScript-like assertions
type Assert struct{}

// True asserts that a condition is true
func (Assert) True(condition bool, message string) error {
	if !condition {
		return NewError(message, ValidationError)
	}
	return nil
}

// False asserts that a condition is false
func (Assert) False(condition bool, message string) error {
	if condition {
		return NewError(message, ValidationError)
	}
	return nil
}

// Equal asserts that two values are equal
func (Assert) Equal(expected, actual interface{}, message string) error {
	if expected != actual {
		return NewError(fmt.Sprintf("%s: expected %v, got %v", message, expected, actual), ValidationError)
	}
	return nil
}

// NotEqual asserts that two values are not equal
func (Assert) NotEqual(expected, actual interface{}, message string) error {
	if expected == actual {
		return NewError(fmt.Sprintf("%s: expected not %v, got %v", message, expected, actual), ValidationError)
	}
	return nil
}

// NotNil asserts that a value is not nil
func (Assert) NotNil(value interface{}, message string) error {
	if value == nil {
		return NewError(message, ValidationError)
	}
	return nil
}

// Nil asserts that a value is nil
func (Assert) Nil(value interface{}, message string) error {
	if value != nil {
		return NewError(fmt.Sprintf("%s: expected nil, got %v", message, value), ValidationError)
	}
	return nil
}

// Global assert instance
var Assertions = Assert{}

// ErrorBoundary represents TypeScript React-like error boundary
type ErrorBoundary struct {
	errorHandler func(*EnhancedError) error
	fallback     func(*EnhancedError) interface{}
}

// NewErrorBoundary creates a new error boundary
func NewErrorBoundary() *ErrorBoundary {
	return &ErrorBoundary{}
}

// OnError sets the error handler
func (eb *ErrorBoundary) OnError(handler func(*EnhancedError) error) *ErrorBoundary {
	eb.errorHandler = handler
	return eb
}

// WithFallback sets the fallback function
func (eb *ErrorBoundary) WithFallback(fallback func(*EnhancedError) interface{}) *ErrorBoundary {
	eb.fallback = fallback
	return eb
}

// Wrap wraps a function with error boundary
func Wrap[T any](eb *ErrorBoundary, fn func() (T, error)) (T, error) {
	result, err := fn()
	
	if err != nil {
		var enhancedErr *EnhancedError
		
		// Convert to enhanced error if not already
		if e, ok := err.(*EnhancedError); ok {
			enhancedErr = e
		} else {
			enhancedErr = NewErrorWithCause("Wrapped error", err)
		}
		
		// Call error handler if present
		if eb.errorHandler != nil {
			if handlerErr := eb.errorHandler(enhancedErr); handlerErr != nil {
				// Error handler failed, return original error
				return result, enhancedErr
			}
		}
		
		// Use fallback if available
		if eb.fallback != nil {
			if fallbackResult, ok := eb.fallback(enhancedErr).(T); ok {
				return fallbackResult, nil
			}
		}
		
		return result, enhancedErr
	}
	
	return result, nil
}

// ErrorFormatter provides different error formatting options
type ErrorFormatter struct{}

// FormatShort returns a short error format
func (ErrorFormatter) FormatShort(err error) string {
	if e, ok := err.(*EnhancedError); ok {
		return fmt.Sprintf("[%s] %s", e.code, e.message)
	}
	return err.Error()
}

// FormatDetailed returns a detailed error format
func (ErrorFormatter) FormatDetailed(err error) string {
	if e, ok := err.(*EnhancedError); ok {
		return e.String()
	}
	return err.Error()
}

// FormatJSON returns JSON formatted error
func (ErrorFormatter) FormatJSON(err error) map[string]interface{} {
	if e, ok := err.(*EnhancedError); ok {
		return e.ToJSON()
	}
	
	return map[string]interface{}{
		"message": err.Error(),
		"code":    string(UnknownError),
		"timestamp": time.Now().Format(time.RFC3339),
	}
}

// Global error formatter
var Formatter = ErrorFormatter{}

// IsErrorCode checks if an error has a specific error code
func IsErrorCode(err error, code ErrorCode) bool {
	if e, ok := err.(*EnhancedError); ok {
		return e.Is(code)
	}
	return false
}

// GetErrorCode extracts error code from any error
func GetErrorCode(err error) ErrorCode {
	if e, ok := err.(*EnhancedError); ok {
		return e.Code()
	}
	return UnknownError
}

// WrapError wraps a regular error as an enhanced error
func WrapError(err error, message string, code ...ErrorCode) *EnhancedError {
	return NewErrorWithCause(message, err, code...)
}

// ValidationError helper functions
func NewValidationError(message string) *EnhancedError {
	return NewError(message, ValidationError)
}

func NewNetworkError(message string) *EnhancedError {
	return NewError(message, NetworkError)
}

func NewAuthError(message string) *EnhancedError {
	return NewError(message, AuthError)
}

func NewNotFoundError(message string) *EnhancedError {
	return NewError(message, NotFoundError)
}

func NewInternalError(message string) *EnhancedError {
	return NewError(message, InternalError)
}

func NewTimeoutError(message string) *EnhancedError {
	return NewError(message, TimeoutError)
}