package testing

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"
	"typescript-golang/types"
	"typescript-golang/utils"
)

// TestStatus represents the status of a test
type TestStatus int

const (
	Pending TestStatus = iota
	Running
	Passed
	Failed
	Skipped
)

func (s TestStatus) String() string {
	switch s {
	case Pending:
		return "pending"
	case Running:
		return "running"
	case Passed:
		return "passed"
	case Failed:
		return "failed"
	case Skipped:
		return "skipped"
	default:
		return "unknown"
	}
}

// TestResult represents the result of a test
type TestResult struct {
	Name        string        `json:"name"`
	Status      TestStatus    `json:"status"`
	Duration    time.Duration `json:"duration"`
	Error       error         `json:"error,omitempty"`
	StartTime   time.Time     `json:"start_time"`
	EndTime     time.Time     `json:"end_time"`
	Description string        `json:"description"`
}

// TestFunction represents a test function
type TestFunction func(*TestContext)

// TestContext provides context and utilities for tests
type TestContext struct {
	name        string
	description string
	startTime   time.Time
	endTime     time.Time
	status      TestStatus
	error       error
	cleanup     []func()
	timeout     time.Duration
}

// NewTestContext creates a new test context
func NewTestContext(name, description string) *TestContext {
	return &TestContext{
		name:        name,
		description: description,
		status:      Pending,
		cleanup:     make([]func(), 0),
		timeout:     30 * time.Second, // Default timeout
	}
}

// SetTimeout sets the test timeout
func (tc *TestContext) SetTimeout(timeout time.Duration) {
	tc.timeout = timeout
}

// AddCleanup adds a cleanup function
func (tc *TestContext) AddCleanup(fn func()) {
	tc.cleanup = append(tc.cleanup, fn)
}

// runCleanup runs all cleanup functions
func (tc *TestContext) runCleanup() {
	for i := len(tc.cleanup) - 1; i >= 0; i-- {
		tc.cleanup[i]()
	}
}

// Suite represents a test suite (like describe block in Jest)
type Suite struct {
	name         string
	description  string
	tests        []*Test
	suites       []*Suite
	beforeEach   []func(*TestContext)
	afterEach    []func(*TestContext)
	beforeAll    []func()
	afterAll     []func()
	parent       *Suite
	mu           sync.RWMutex
}

// NewSuite creates a new test suite
func NewSuite(name, description string) *Suite {
	return &Suite{
		name:        name,
		description: description,
		tests:       make([]*Test, 0),
		suites:      make([]*Suite, 0),
		beforeEach:  make([]func(*TestContext), 0),
		afterEach:   make([]func(*TestContext), 0),
		beforeAll:   make([]func(), 0),
		afterAll:    make([]func(), 0),
	}
}

// Test represents an individual test (like it block in Jest)
type Test struct {
	name        string
	description string
	fn          TestFunction
	suite       *Suite
	result      *TestResult
	skip        bool
	only        bool
}

// NewTest creates a new test
func NewTest(name, description string, fn TestFunction) *Test {
	return &Test{
		name:        name,
		description: description,
		fn:          fn,
		skip:        false,
		only:        false,
	}
}

// Describe creates a test suite (like describe() in Jest)
func Describe(name string, fn func()) *Suite {
	suite := NewSuite(name, "")
	currentSuite = suite
	fn()
	return suite
}

// It creates a test (like it() in Jest)
func It(description string, fn TestFunction) *Test {
	test := NewTest("", description, fn)
	if currentSuite != nil {
		currentSuite.AddTest(test)
	}
	return test
}

// TestFunc creates a test (alias for It)
func TestFunc(description string, fn TestFunction) *Test {
	return It(description, fn)
}

// BeforeEach adds a before each hook
func BeforeEach(fn func(*TestContext)) {
	if currentSuite != nil {
		currentSuite.BeforeEach(fn)
	}
}

// AfterEach adds an after each hook
func AfterEach(fn func(*TestContext)) {
	if currentSuite != nil {
		currentSuite.AfterEach(fn)
	}
}

// BeforeAll adds a before all hook
func BeforeAll(fn func()) {
	if currentSuite != nil {
		currentSuite.BeforeAll(fn)
	}
}

// AfterAll adds an after all hook
func AfterAll(fn func()) {
	if currentSuite != nil {
		currentSuite.AfterAll(fn)
	}
}

// Global current suite for building test structure
var currentSuite *Suite

// AddTest adds a test to the suite
func (s *Suite) AddTest(test *Test) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	test.suite = s
	s.tests = append(s.tests, test)
}

// AddSuite adds a nested suite
func (s *Suite) AddSuite(nested *Suite) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	nested.parent = s
	s.suites = append(s.suites, nested)
}

// BeforeEach adds a before each hook
func (s *Suite) BeforeEach(fn func(*TestContext)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.beforeEach = append(s.beforeEach, fn)
}

// AfterEach adds an after each hook
func (s *Suite) AfterEach(fn func(*TestContext)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.afterEach = append(s.afterEach, fn)
}

// BeforeAll adds a before all hook
func (s *Suite) BeforeAll(fn func()) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.beforeAll = append(s.beforeAll, fn)
}

// AfterAll adds an after all hook
func (s *Suite) AfterAll(fn func()) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.afterAll = append(s.afterAll, fn)
}

// Run executes the test suite
func (s *Suite) Run() *SuiteResult {
	result := &SuiteResult{
		Name:        s.name,
		Description: s.description,
		StartTime:   time.Now(),
		Tests:       make([]*TestResult, 0),
		Suites:      make([]*SuiteResult, 0),
	}
	
	// Run before all hooks
	for _, hook := range s.beforeAll {
		hook()
	}
	
	// Run tests
	for _, test := range s.tests {
		if test.skip {
			continue
		}
		
		testResult := s.runTest(test)
		result.Tests = append(result.Tests, testResult)
		
		switch testResult.Status {
		case Passed:
			result.PassedCount++
		case Failed:
			result.FailedCount++
		case Skipped:
			result.SkippedCount++
		}
		result.TotalCount++
	}
	
	// Run nested suites
	for _, nested := range s.suites {
		nestedResult := nested.Run()
		result.Suites = append(result.Suites, nestedResult)
		
		result.TotalCount += nestedResult.TotalCount
		result.PassedCount += nestedResult.PassedCount
		result.FailedCount += nestedResult.FailedCount
		result.SkippedCount += nestedResult.SkippedCount
	}
	
	// Run after all hooks
	for _, hook := range s.afterAll {
		hook()
	}
	
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)
	
	return result
}

// runTest executes a single test
func (s *Suite) runTest(test *Test) *TestResult {
	ctx := NewTestContext(test.name, test.description)
	ctx.startTime = time.Now()
	ctx.status = Running
	
	// Run before each hooks
	for _, hook := range s.beforeEach {
		hook(ctx)
	}
	
	// Run the test with timeout
	done := make(chan bool, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				ctx.status = Failed
				ctx.error = fmt.Errorf("panic: %v", r)
			}
			done <- true
		}()
		
		test.fn(ctx)
		if ctx.status != Failed {
			ctx.status = Passed
		}
	}()
	
	// Wait for completion or timeout
	select {
	case <-done:
		// Test completed
	case <-time.After(ctx.timeout):
		ctx.status = Failed
		ctx.error = fmt.Errorf("test timeout after %v", ctx.timeout)
	}
	
	ctx.endTime = time.Now()
	
	// Run after each hooks
	for _, hook := range s.afterEach {
		hook(ctx)
	}
	
	// Run cleanup
	ctx.runCleanup()
	
	return &TestResult{
		Name:        test.name,
		Description: test.description,
		Status:      ctx.status,
		Duration:    ctx.endTime.Sub(ctx.startTime),
		Error:       ctx.error,
		StartTime:   ctx.startTime,
		EndTime:     ctx.endTime,
	}
}

// SuiteResult represents the result of running a test suite
type SuiteResult struct {
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	StartTime    time.Time       `json:"start_time"`
	EndTime      time.Time       `json:"end_time"`
	Duration     time.Duration   `json:"duration"`
	Tests        []*TestResult   `json:"tests"`
	Suites       []*SuiteResult  `json:"suites"`
	TotalCount   int             `json:"total_count"`
	PassedCount  int             `json:"passed_count"`
	FailedCount  int             `json:"failed_count"`
	SkippedCount int             `json:"skipped_count"`
}

// Expect creates an expectation for testing (like expect() in Jest)
func Expect(actual interface{}) *Expectation {
	return &Expectation{
		actual: actual,
		not:    false,
	}
}

// Expectation represents a test expectation with matchers
type Expectation struct {
	actual interface{}
	not    bool
}

// Not negates the expectation
func (e *Expectation) Not() *Expectation {
	return &Expectation{
		actual: e.actual,
		not:    !e.not,
	}
}

// ToBe checks for strict equality (like toBe() in Jest)
func (e *Expectation) ToBe(expected interface{}) {
	equal := reflect.DeepEqual(e.actual, expected)
	if e.not {
		equal = !equal
	}
	
	if !equal {
		var message string
		if e.not {
			message = fmt.Sprintf("Expected %v not to be %v", e.actual, expected)
		} else {
			message = fmt.Sprintf("Expected %v to be %v", e.actual, expected)
		}
		panic(types.NewValidationError(message))
	}
}

// ToEqual checks for deep equality (like toEqual() in Jest)
func (e *Expectation) ToEqual(expected interface{}) {
	e.ToBe(expected) // Same as ToBe for now
}

// ToBeTrue checks if value is true
func (e *Expectation) ToBeTrue() {
	e.ToBe(true)
}

// ToBeFalse checks if value is false
func (e *Expectation) ToBeFalse() {
	e.ToBe(false)
}

// ToBeNil checks if value is nil
func (e *Expectation) ToBeNil() {
	if e.not {
		if e.actual == nil {
			panic(types.NewValidationError("Expected value not to be nil"))
		}
	} else {
		if e.actual != nil {
			panic(types.NewValidationError(fmt.Sprintf("Expected %v to be nil", e.actual)))
		}
	}
}

// ToContain checks if slice/string contains value
func (e *Expectation) ToContain(expected interface{}) {
	contains := false
	
	actualValue := reflect.ValueOf(e.actual)
	switch actualValue.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < actualValue.Len(); i++ {
			if reflect.DeepEqual(actualValue.Index(i).Interface(), expected) {
				contains = true
				break
			}
		}
	case reflect.String:
		if expectedStr, ok := expected.(string); ok {
			contains = strings.Contains(actualValue.String(), expectedStr)
		}
	}
	
	if e.not {
		contains = !contains
	}
	
	if !contains {
		var message string
		if e.not {
			message = fmt.Sprintf("Expected %v not to contain %v", e.actual, expected)
		} else {
			message = fmt.Sprintf("Expected %v to contain %v", e.actual, expected)
		}
		panic(types.NewValidationError(message))
	}
}

// ToHaveLength checks the length of slice/string/map
func (e *Expectation) ToHaveLength(expected int) {
	var length int
	
	actualValue := reflect.ValueOf(e.actual)
	switch actualValue.Kind() {
	case reflect.Slice, reflect.Array, reflect.String, reflect.Map:
		length = actualValue.Len()
	default:
		panic(types.NewValidationError(fmt.Sprintf("Expected %v to have a length property", e.actual)))
	}
	
	equal := length == expected
	if e.not {
		equal = !equal
	}
	
	if !equal {
		var message string
		if e.not {
			message = fmt.Sprintf("Expected %v not to have length %d, but got %d", e.actual, expected, length)
		} else {
			message = fmt.Sprintf("Expected %v to have length %d, but got %d", e.actual, expected, length)
		}
		panic(types.NewValidationError(message))
	}
}

// ToThrow checks if function throws an error
func (e *Expectation) ToThrow(expectedError ...string) {
	fn, ok := e.actual.(func())
	if !ok {
		panic(types.NewValidationError("ToThrow can only be used with functions"))
	}
	
	var thrownError error
	func() {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					thrownError = err
				} else {
					thrownError = fmt.Errorf("%v", r)
				}
			}
		}()
		
		fn()
	}()
	
	threw := thrownError != nil
	if e.not {
		threw = !threw
	}
	
	if !threw {
		var message string
		if e.not {
			message = "Expected function not to throw an error"
		} else {
			message = "Expected function to throw an error"
		}
		panic(types.NewValidationError(message))
	}
	
	// Check error message if provided
	if len(expectedError) > 0 && thrownError != nil {
		expectedMsg := expectedError[0]
		actualMsg := thrownError.Error()
		if !strings.Contains(actualMsg, expectedMsg) {
			message := fmt.Sprintf("Expected error to contain '%s', but got '%s'", expectedMsg, actualMsg)
			panic(types.NewValidationError(message))
		}
	}
}

// TestRunner manages and runs tests
type TestRunner struct {
	suites   []*Suite
	reporter Reporter
}

// NewTestRunner creates a new test runner
func NewTestRunner() *TestRunner {
	return &TestRunner{
		suites:   make([]*Suite, 0),
		reporter: &ConsoleReporter{},
	}
}

// AddSuite adds a suite to the runner
func (tr *TestRunner) AddSuite(suite *Suite) {
	tr.suites = append(tr.suites, suite)
}

// SetReporter sets the test reporter
func (tr *TestRunner) SetReporter(reporter Reporter) {
	tr.reporter = reporter
}

// Run executes all test suites
func (tr *TestRunner) Run() *RunResult {
	result := &RunResult{
		StartTime: time.Now(),
		Suites:    make([]*SuiteResult, 0),
	}
	
	for _, suite := range tr.suites {
		suiteResult := suite.Run()
		result.Suites = append(result.Suites, suiteResult)
		
		result.TotalCount += suiteResult.TotalCount
		result.PassedCount += suiteResult.PassedCount
		result.FailedCount += suiteResult.FailedCount
		result.SkippedCount += suiteResult.SkippedCount
	}
	
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)
	
	tr.reporter.Report(result)
	
	return result
}

// RunResult represents the overall test run result
type RunResult struct {
	StartTime    time.Time       `json:"start_time"`
	EndTime      time.Time       `json:"end_time"`
	Duration     time.Duration   `json:"duration"`
	Suites       []*SuiteResult  `json:"suites"`
	TotalCount   int             `json:"total_count"`
	PassedCount  int             `json:"passed_count"`
	FailedCount  int             `json:"failed_count"`
	SkippedCount int             `json:"skipped_count"`
}

// Reporter interface for test reporting
type Reporter interface {
	Report(result *RunResult)
}

// ConsoleReporter reports test results to console
type ConsoleReporter struct{}

// Report prints test results to console
func (cr *ConsoleReporter) Report(result *RunResult) {
	fmt.Println("\n🧪 Test Results")
	fmt.Println("===============")
	
	for _, suite := range result.Suites {
		cr.reportSuite(suite, 0)
	}
	
	fmt.Printf("\nSummary:\n")
	fmt.Printf("  Total: %d\n", result.TotalCount)
	fmt.Printf("  Passed: %d\n", result.PassedCount)
	fmt.Printf("  Failed: %d\n", result.FailedCount)
	fmt.Printf("  Skipped: %d\n", result.SkippedCount)
	fmt.Printf("  Duration: %v\n", result.Duration)
	
	if result.FailedCount > 0 {
		fmt.Printf("\n❌ %d test(s) failed\n", result.FailedCount)
	} else {
		fmt.Printf("\n✅ All tests passed!\n")
	}
}

// reportSuite reports a single suite
func (cr *ConsoleReporter) reportSuite(suite *SuiteResult, indent int) {
	prefix := strings.Repeat("  ", indent)
	fmt.Printf("%s📁 %s\n", prefix, suite.Name)
	
	for _, test := range suite.Tests {
		status := ""
		switch test.Status {
		case Passed:
			status = "✅"
		case Failed:
			status = "❌"
		case Skipped:
			status = "⏭️"
		default:
			status = "❓"
		}
		
		fmt.Printf("%s  %s %s (%v)\n", prefix, status, test.Description, test.Duration)
		
		if test.Error != nil {
			fmt.Printf("%s    Error: %v\n", prefix, test.Error)
		}
	}
	
	for _, nested := range suite.Suites {
		cr.reportSuite(nested, indent+1)
	}
}

// RunExampleTests demonstrates the testing framework
func RunExampleTests() {
	runner := NewTestRunner()
	
	// Array utilities test suite
	arrayTestSuite := Describe("Array Utilities", func() {
		BeforeAll(func() {
			fmt.Println("Setting up array tests...")
		})
		
		AfterAll(func() {
			fmt.Println("Cleaning up array tests...")
		})
		
		BeforeEach(func(ctx *TestContext) {
			fmt.Printf("Starting test: %s\n", ctx.description)
		})
		
		AfterEach(func(ctx *TestContext) {
			fmt.Printf("Finished test: %s\n", ctx.description)
		})
		
		It("should map array elements correctly", func(ctx *TestContext) {
			numbers := []int{1, 2, 3, 4, 5}
			doubled := utils.Map(numbers, func(x int) int { return x * 2 })
			
			Expect(doubled).ToEqual([]int{2, 4, 6, 8, 10})
			Expect(doubled).ToHaveLength(5)
			Expect(doubled).ToContain(6)
		})
		
		It("should filter array elements correctly", func(ctx *TestContext) {
			numbers := []int{1, 2, 3, 4, 5, 6}
			evens := utils.Filter(numbers, func(x int) bool { return x%2 == 0 })
			
			Expect(evens).ToEqual([]int{2, 4, 6})
			Expect(evens).ToHaveLength(3)
		})
		
		It("should reduce array correctly", func(ctx *TestContext) {
			numbers := []int{1, 2, 3, 4, 5}
			sum := utils.Reduce(numbers, func(acc, x int) int { return acc + x }, 0)
			
			Expect(sum).ToBe(15)
		})
		
		It("should find elements correctly", func(ctx *TestContext) {
			numbers := []int{1, 2, 3, 4, 5}
			found := utils.Find(numbers, func(x int) bool { return x > 3 })
			
			Expect(found.IsSome()).ToBeTrue()
			Expect(found.Get()).ToBe(4)
		})
	})
	
	// Optional types test suite
	optionalTestSuite := Describe("Optional Types", func() {
		It("should handle Some values correctly", func(ctx *TestContext) {
			opt := types.Some("test value")
			
			Expect(opt.IsSome()).ToBeTrue()
			Expect(opt.IsNone()).ToBeFalse()
			Expect(opt.Get()).ToBe("test value")
			Expect(opt.GetOrDefault("default")).ToBe("test value")
		})
		
		It("should handle None values correctly", func(ctx *TestContext) {
			opt := types.None[string]()
			
			Expect(opt.IsSome()).ToBeFalse()
			Expect(opt.IsNone()).ToBeTrue()
			Expect(opt.GetOrDefault("default")).ToBe("default")
		})
		
		It("should panic when getting None value", func(ctx *TestContext) {
			opt := types.None[string]()
			
			Expect(func() {
				opt.Get()
			}).ToThrow("panic")
		})
	})
	
	// Add all suites to runner
	runner.AddSuite(arrayTestSuite)
	runner.AddSuite(optionalTestSuite)
	
	// Run all tests
	runner.Run()
}