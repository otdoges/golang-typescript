package types

import (
	"fmt"
	"reflect"
	"sync"
	"time"
)

// EventListener represents a function that handles events
type EventListener[T any] func(event T)

// EventEmitter represents TypeScript's EventEmitter pattern
type EventEmitter[T any] struct {
	listeners    map[string][]EventListener[T]
	onceListeners map[string][]EventListener[T]
	mu           sync.RWMutex
	maxListeners int
}

// NewEventEmitter creates a new EventEmitter (like new EventEmitter() in TypeScript)
func NewEventEmitter[T any]() *EventEmitter[T] {
	return &EventEmitter[T]{
		listeners:     make(map[string][]EventListener[T]),
		onceListeners: make(map[string][]EventListener[T]),
		maxListeners:  10, // Default max listeners like Node.js
	}
}

// On adds a listener for the specified event (like emitter.on() in TypeScript)
func (ee *EventEmitter[T]) On(event string, listener EventListener[T]) *EventEmitter[T] {
	ee.mu.Lock()
	defer ee.mu.Unlock()
	
	ee.listeners[event] = append(ee.listeners[event], listener)
	
	// Check max listeners
	if len(ee.listeners[event]) > ee.maxListeners {
		fmt.Printf("Warning: EventEmitter has %d listeners for event '%s'. This may indicate a memory leak.\n", 
			len(ee.listeners[event]), event)
	}
	
	return ee
}

// AddListener is an alias for On (like emitter.addListener() in Node.js)
func (ee *EventEmitter[T]) AddListener(event string, listener EventListener[T]) *EventEmitter[T] {
	return ee.On(event, listener)
}

// Once adds a one-time listener (like emitter.once() in TypeScript)
func (ee *EventEmitter[T]) Once(event string, listener EventListener[T]) *EventEmitter[T] {
	ee.mu.Lock()
	defer ee.mu.Unlock()
	
	ee.onceListeners[event] = append(ee.onceListeners[event], listener)
	return ee
}

// Off removes a listener (like emitter.off() in TypeScript)
func (ee *EventEmitter[T]) Off(event string, listener EventListener[T]) *EventEmitter[T] {
	ee.mu.Lock()
	defer ee.mu.Unlock()
	
	// Remove from regular listeners
	if listeners, exists := ee.listeners[event]; exists {
		for i, l := range listeners {
			if reflect.ValueOf(l).Pointer() == reflect.ValueOf(listener).Pointer() {
				ee.listeners[event] = append(listeners[:i], listeners[i+1:]...)
				break
			}
		}
	}
	
	// Remove from once listeners
	if onceListeners, exists := ee.onceListeners[event]; exists {
		for i, l := range onceListeners {
			if reflect.ValueOf(l).Pointer() == reflect.ValueOf(listener).Pointer() {
				ee.onceListeners[event] = append(onceListeners[:i], onceListeners[i+1:]...)
				break
			}
		}
	}
	
	return ee
}

// RemoveListener is an alias for Off
func (ee *EventEmitter[T]) RemoveListener(event string, listener EventListener[T]) *EventEmitter[T] {
	return ee.Off(event, listener)
}

// RemoveAllListeners removes all listeners for an event or all events
func (ee *EventEmitter[T]) RemoveAllListeners(event ...string) *EventEmitter[T] {
	ee.mu.Lock()
	defer ee.mu.Unlock()
	
	if len(event) == 0 {
		// Remove all listeners for all events
		ee.listeners = make(map[string][]EventListener[T])
		ee.onceListeners = make(map[string][]EventListener[T])
	} else {
		// Remove all listeners for specific event
		eventName := event[0]
		delete(ee.listeners, eventName)
		delete(ee.onceListeners, eventName)
	}
	
	return ee
}

// Emit triggers all listeners for the specified event (like emitter.emit() in TypeScript)
func (ee *EventEmitter[T]) Emit(event string, data T) bool {
	ee.mu.Lock()
	
	// Get regular listeners
	listeners := make([]EventListener[T], len(ee.listeners[event]))
	copy(listeners, ee.listeners[event])
	
	// Get once listeners and remove them
	onceListeners := make([]EventListener[T], len(ee.onceListeners[event]))
	copy(onceListeners, ee.onceListeners[event])
	delete(ee.onceListeners, event)
	
	ee.mu.Unlock()
	
	hadListeners := len(listeners) > 0 || len(onceListeners) > 0
	
	// Execute regular listeners
	for _, listener := range listeners {
		go listener(data) // Execute asynchronously like in Node.js
	}
	
	// Execute once listeners
	for _, listener := range onceListeners {
		go listener(data)
	}
	
	return hadListeners
}

// EmitSync emits event synchronously (all listeners execute before returning)
func (ee *EventEmitter[T]) EmitSync(event string, data T) bool {
	ee.mu.Lock()
	
	// Get regular listeners
	listeners := make([]EventListener[T], len(ee.listeners[event]))
	copy(listeners, ee.listeners[event])
	
	// Get once listeners and remove them
	onceListeners := make([]EventListener[T], len(ee.onceListeners[event]))
	copy(onceListeners, ee.onceListeners[event])
	delete(ee.onceListeners, event)
	
	ee.mu.Unlock()
	
	hadListeners := len(listeners) > 0 || len(onceListeners) > 0
	
	// Execute regular listeners synchronously
	for _, listener := range listeners {
		listener(data)
	}
	
	// Execute once listeners synchronously
	for _, listener := range onceListeners {
		listener(data)
	}
	
	return hadListeners
}

// ListenerCount returns the number of listeners for an event
func (ee *EventEmitter[T]) ListenerCount(event string) int {
	ee.mu.RLock()
	defer ee.mu.RUnlock()
	
	return len(ee.listeners[event]) + len(ee.onceListeners[event])
}

// EventNames returns all event names that have listeners
func (ee *EventEmitter[T]) EventNames() []string {
	ee.mu.RLock()
	defer ee.mu.RUnlock()
	
	eventSet := make(map[string]bool)
	
	for event := range ee.listeners {
		eventSet[event] = true
	}
	
	for event := range ee.onceListeners {
		eventSet[event] = true
	}
	
	var events []string
	for event := range eventSet {
		events = append(events, event)
	}
	
	return events
}

// SetMaxListeners sets the maximum number of listeners per event
func (ee *EventEmitter[T]) SetMaxListeners(max int) *EventEmitter[T] {
	ee.mu.Lock()
	defer ee.mu.Unlock()
	
	ee.maxListeners = max
	return ee
}

// GetMaxListeners returns the maximum number of listeners per event
func (ee *EventEmitter[T]) GetMaxListeners() int {
	ee.mu.RLock()
	defer ee.mu.RUnlock()
	
	return ee.maxListeners
}

// Observable represents TypeScript's Observable pattern (RxJS-like)
type Observable[T any] struct {
	emitter *EventEmitter[T]
}

// NewObservable creates a new Observable
func NewObservable[T any]() *Observable[T] {
	return &Observable[T]{
		emitter: NewEventEmitter[T](),
	}
}

// Subscribe subscribes to the observable (like observable.subscribe() in RxJS)
func (obs *Observable[T]) Subscribe(observer EventListener[T]) *Subscription {
	obs.emitter.On("data", observer)
	
	return &Subscription{
		unsubscribe: func() {
			obs.emitter.Off("data", observer)
		},
	}
}

// Next emits the next value
func (obs *Observable[T]) Next(value T) {
	obs.emitter.Emit("data", value)
}

// Error emits an error
func (obs *Observable[T]) Error(err error) {
	obs.emitter.Emit("error", *new(T)) // Placeholder for error event
}

// Complete signals completion
func (obs *Observable[T]) Complete() {
	obs.emitter.Emit("complete", *new(T))
}

// Subscription represents a subscription to an observable
type Subscription struct {
	unsubscribe func()
}

// Unsubscribe cancels the subscription
func (s *Subscription) Unsubscribe() {
	if s.unsubscribe != nil {
		s.unsubscribe()
	}
}

// Subject represents a subject (both observable and observer)
type Subject[T any] struct {
	*Observable[T]
	subscribers []EventListener[T]
	mu          sync.RWMutex
}

// NewSubject creates a new Subject
func NewSubject[T any]() *Subject[T] {
	return &Subject[T]{
		Observable:  NewObservable[T](),
		subscribers: make([]EventListener[T], 0),
	}
}

// Next emits a value to all subscribers
func (s *Subject[T]) Next(value T) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	for _, subscriber := range s.subscribers {
		go subscriber(value)
	}
}

// Subscribe adds a subscriber
func (s *Subject[T]) Subscribe(observer EventListener[T]) *Subscription {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.subscribers = append(s.subscribers, observer)
	
	return &Subscription{
		unsubscribe: func() {
			s.mu.Lock()
			defer s.mu.Unlock()
			
			for i, sub := range s.subscribers {
				if reflect.ValueOf(sub).Pointer() == reflect.ValueOf(observer).Pointer() {
					s.subscribers = append(s.subscribers[:i], s.subscribers[i+1:]...)
					break
				}
			}
		},
	}
}

// Event represents a generic event with data
type Event[T any] struct {
	Type      string
	Data      T
	Timestamp time.Time
	Source    interface{}
}

// NewEvent creates a new event
func NewEvent[T any](eventType string, data T, source interface{}) *Event[T] {
	return &Event[T]{
		Type:      eventType,
		Data:      data,
		Timestamp: time.Now(),
		Source:    source,
	}
}

// EventBus represents a global event bus for decoupled communication
type EventBus struct {
	emitters map[string]*EventEmitter[interface{}]
	mu       sync.RWMutex
}

// NewEventBus creates a new EventBus
func NewEventBus() *EventBus {
	return &EventBus{
		emitters: make(map[string]*EventEmitter[interface{}]),
	}
}

// Publish publishes an event to the bus
func (eb *EventBus) Publish(eventType string, data interface{}) {
	eb.mu.RLock()
	emitter, exists := eb.emitters[eventType]
	eb.mu.RUnlock()
	
	if exists {
		emitter.Emit("data", data)
	}
}

// Subscribe subscribes to events of a specific type
func (eb *EventBus) Subscribe(eventType string, listener EventListener[interface{}]) *Subscription {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	
	if _, exists := eb.emitters[eventType]; !exists {
		eb.emitters[eventType] = NewEventEmitter[interface{}]()
	}
	
	eb.emitters[eventType].On("data", listener)
	
	return &Subscription{
		unsubscribe: func() {
			eb.mu.RLock()
			defer eb.mu.RUnlock()
			
			if emitter, exists := eb.emitters[eventType]; exists {
				emitter.Off("data", listener)
			}
		},
	}
}

// WaitForChan waits for an event and returns a channel
func (ee *EventEmitter[T]) WaitForChan(event string, timeout ...time.Duration) <-chan T {
	resultChan := make(chan T, 1)
	timeoutDuration := 30 * time.Second // Default timeout
	
	if len(timeout) > 0 {
		timeoutDuration = timeout[0]
	}
	
	// Set up one-time listener
	listener := func(data T) {
		select {
		case resultChan <- data:
		default:
		}
	}
	
	ee.Once(event, listener)
	
	// Set up timeout
	go func() {
		time.Sleep(timeoutDuration)
		close(resultChan)
	}()
	
	return resultChan
}

// WaitFor waits for an event with optional timeout
func (ee *EventEmitter[T]) WaitFor(event string, timeout ...time.Duration) (T, error) {
	resultChan := ee.WaitForChan(event, timeout...)
	
	if result, ok := <-resultChan; ok {
		return result, nil
	}
	
	var zero T
	timeoutDuration := 30 * time.Second
	if len(timeout) > 0 {
		timeoutDuration = timeout[0]
	}
	return zero, fmt.Errorf("timeout waiting for event '%s' after %v", event, timeoutDuration)
}

// Pipe creates a pipeline of event transformations (like RxJS operators)
type EventPipe[T, U any] struct {
	source      *Observable[T]
	transformer func(T) U
}

// NewEventPipe creates a new event pipe
func NewEventPipe[T, U any](source *Observable[T], transformer func(T) U) *EventPipe[T, U] {
	return &EventPipe[T, U]{
		source:      source,
		transformer: transformer,
	}
}

// Subscribe subscribes to the transformed stream
func (ep *EventPipe[T, U]) Subscribe(observer EventListener[U]) *Subscription {
	return ep.source.Subscribe(func(value T) {
		transformed := ep.transformer(value)
		observer(transformed)
	})
}

// Filter operator for event streams
func Filter[T any](source *Observable[T], predicate func(T) bool) *Observable[T] {
	result := NewObservable[T]()
	
	source.Subscribe(func(value T) {
		if predicate(value) {
			result.Next(value)
		}
	})
	
	return result
}

// Map operator for event streams  
func ObservableMap[T, U any](source *Observable[T], transformer func(T) U) *Observable[U] {
	result := NewObservable[U]()
	
	source.Subscribe(func(value T) {
		transformed := transformer(value)
		result.Next(transformed)
	})
	
	return result
}

// Debounce operator - emits only after a pause in events
func Debounce[T any](source *Observable[T], delay time.Duration) *Observable[T] {
	result := NewObservable[T]()
	var timer *time.Timer
	var mu sync.Mutex
	
	source.Subscribe(func(value T) {
		mu.Lock()
		defer mu.Unlock()
		
		if timer != nil {
			timer.Stop()
		}
		
		timer = time.AfterFunc(delay, func() {
			result.Next(value)
		})
	})
	
	return result
}

// Throttle operator - emits at most once per time period
func Throttle[T any](source *Observable[T], interval time.Duration) *Observable[T] {
	result := NewObservable[T]()
	var lastEmit time.Time
	var mu sync.Mutex
	
	source.Subscribe(func(value T) {
		mu.Lock()
		defer mu.Unlock()
		
		now := time.Now()
		if now.Sub(lastEmit) >= interval {
			lastEmit = now
			result.Next(value)
		}
	})
	
	return result
}